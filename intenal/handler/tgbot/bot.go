package tgbot

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"runtime/debug"
	"sync"
	"time"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error

const (
	requestInProgress = "in progress"
	requestApproved   = "approved"
	requestRejected   = "rejected"
)

const (
	roleUser  = "user"
	roleAdmin = "admin"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	log *logger.Logger

	cmdView      map[string]ViewFunc
	callbackView map[string]ViewFunc

	requestService service.RequestService
	userService    service.UserService

	mu sync.RWMutex
}

func NewBot(bot *tgbotapi.BotAPI,
	log *logger.Logger,
	requestService service.RequestService,
	userService service.UserService) *Bot {
	return &Bot{
		bot:            bot,
		log:            log,
		requestService: requestService,
		userService:    userService,
	}
}

func (b *Bot) RegisterCommandView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}

	b.cmdView[cmd] = view
}

func (b *Bot) RegisterCommandCallback(callback string, view ViewFunc) {
	if b.callbackView == nil {
		b.callbackView = make(map[string]ViewFunc)
	}

	b.callbackView[callback] = view
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

			b.handlerUpdate(updateCtx, &update)

			cancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handlerUpdate(ctx context.Context, update *tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			b.log.Error("panic recovered: %v, %s", p, string(debug.Stack()))
		}
	}()

	// Если пришло сообщение
	if update.Message != nil {
		b.log.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

		var view ViewFunc

		cmd := update.Message.Command()

		cmdView, ok := b.cmdView[cmd]
		if !ok {
			return
		}

		view = cmdView

		if err := view(ctx, b.bot, update); err != nil {
			b.log.Error("failed to handle update: %v", err)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "internal error")
			if _, err := b.bot.Send(msg); err != nil {
				b.log.Error("failed to send message: %v", err)
			}
			return
		}
		// Если нажали кнопку
	} else if update.CallbackQuery != nil {
		b.log.Info("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackData())

		var callback ViewFunc
		err, callbackView := CallbackStrings(update, b)
		if err != nil {
			b.log.Error("%v", err)
			return
		}

		callback = callbackView

		if err := callback(ctx, b.bot, update); err != nil {
			b.log.Error("failed to handle update: %v", err)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "internal server error")
			if _, err := b.bot.Send(msg); err != nil {
				b.log.Error("failed to send message: %v", err)
			}
			return
		}
	} else if update.ChatJoinRequest != nil {
		b.log.Info("[%s] %s", update.ChatJoinRequest.From.UserName, update.ChatJoinRequest.InviteLink.InviteLink)

		if update.ChatJoinRequest.InviteLink == nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "internal server error")
			if _, err := b.bot.Send(msg); err != nil {
				b.log.Error("failed to send message: %v", err)
			}
			return
		}

		user := &entity.User{
			ID:          update.ChatJoinRequest.From.ID,
			UsernameTg:  update.ChatJoinRequest.From.UserName,
			Phone:       nil,
			ChannelFrom: &update.ChatJoinRequest.InviteLink.InviteLink,
			CreatedAt:   time.Now().Local(),
			Role:        roleUser,
		}
		request := &entity.Request{
			UserID:        update.ChatJoinRequest.From.ID,
			StatusRequest: requestInProgress,
		}

		if err := b.userService.JoinChannel(ctx, user, request); err != nil {
			b.log.Error("userService.JoinChannel: can`t ")
		}
	}
}

func (b *Bot) JoinRequest(update *tgbotapi.Update) error {
	req := update.ChatJoinRequest

	if req == nil {
		return boterror.ErrNil
	}

	return nil
}
