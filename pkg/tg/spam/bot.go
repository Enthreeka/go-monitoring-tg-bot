package spam

import (
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type spammerBots struct {
	storageBot map[string]*tgbotapi.BotAPI

	log *logger.Logger
	mu  sync.RWMutex
}

type NewBot interface {
	InitializeBot(token string)
}

func NewSpammerBot(log *logger.Logger) NewBot {
	return &spammerBots{
		storageBot: make(map[string]*tgbotapi.BotAPI, 20),
		log:        log,
	}
}

func (s *spammerBots) InitializeBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		s.log.Error("failed to load token %v", err)
	}

	s.log.Info("Authorized on account %s", bot.Self.UserName)

	s.mu.Lock()
	s.storageBot[bot.Self.UserName] = bot
	s.mu.Unlock()
}
