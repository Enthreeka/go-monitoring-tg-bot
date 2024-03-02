package balancer

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/spam"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
	"sync/atomic"
	"time"
)

type Balancer interface {
	Prepare(ctx context.Context, notification *entity.Notification, users []entity.User)
}

type BotPool struct {
	spammer  *spam.SpammerBots
	log      *logger.Logger
	current  uint64
	botsName []string

	successCounter int64
}

func NewBalancer(spammer *spam.SpammerBots, log *logger.Logger) Balancer {
	return &BotPool{
		spammer: spammer,
		log:     log,
	}
}

func (b *BotPool) getBots() {
	b.botsName = make([]string, 0, 20)

	b.spammer.Range(func(key, value any) error {
		botName := key.(string)
		b.botsName = append(b.botsName, botName)
		return nil
	})

	b.log.Info("Gets bots: %v", b.botsName)
}

func (b *BotPool) nextIndex() int {
	return int(atomic.AddUint64(&b.current, uint64(1)) % uint64(len(b.botsName)))
}

func (b *BotPool) getNextBot() *tgbotapi.BotAPI {
	next := b.nextIndex()

	l := len(b.botsName) + next

	for i := next; i < l; i++ {
		idx := i % len(b.botsName)

		if i != next {
			atomic.StoreUint64(&b.current, uint64(idx))
		}

		bot, _ := b.spammer.Read(b.botsName[idx])
		return bot
	}
	return nil
}

func (b *BotPool) recursiveSender(notification *entity.Notification, userID int64) func() {
	bot := b.getNextBot()
	if bot == nil {
		b.log.Error("bot has nil pointer")
		return b.recursiveSender(notification, userID)
	}

	err := b.sendMsgToNewUser(notification, userID, bot)
	if err != nil {
		isPossibleToResend := handlerSenderError(err)
		if isPossibleToResend {
			return b.recursiveSender(notification, userID)
		}

		b.log.Error("balancer.sendMsgToNewUser: bot_name:%s %v:", bot.Self.UserName, err)
		return nil
	}

	return nil
}

func (b *BotPool) Prepare(ctx context.Context, notification *entity.Notification, users []entity.User) {
	b.log.Info("Starting balancer sending")
	start := time.Now()

	b.getBots()
	var wg sync.WaitGroup

	for i, user := range users {
		select {
		case <-ctx.Done():
			b.log.Error("context in balancer: %v", ctx.Err())
			return
		default:
			if i%2 == 0 {
				wg.Add(1)
				go func(id int64) {
					defer wg.Done()

					b.recursiveSender(notification, id)

					return
				}(user.ID)
				continue
			}

			b.recursiveSender(notification, user.ID)
		}
	}
	wg.Wait()

	b.log.Info("Balancer is executed fot time: %v", time.Since(start))
	b.log.Info("Count successfully sent messages: %d/%d", b.successCounter, len(users))
}
