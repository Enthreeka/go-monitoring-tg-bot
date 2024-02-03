package tgbot

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"runtime/debug"
	"sync"
	"time"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error

const (
	RequestInProgress = "in progress"
	RequestApproved   = "approved"
	RequestRejected   = "rejected"
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
	channelService service.ChannelService

	mu sync.RWMutex
}

func NewBot(bot *tgbotapi.BotAPI,
	log *logger.Logger,
	requestService service.RequestService,
	userService service.UserService,
	channelService service.ChannelService) *Bot {
	return &Bot{
		bot:            bot,
		log:            log,
		requestService: requestService,
		userService:    userService,
		channelService: channelService,
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

		if err := b.userService.CreateUser(ctx, userUpdateToModel(update)); err != nil {
			b.log.Error("userService.CreateUser: failed to create user: %v", err)
			return
		}

		var view ViewFunc

		cmd := update.Message.Command()

		cmdView, ok := b.cmdView[cmd]
		if !ok {
			return
		}

		view = cmdView

		if err := view(ctx, b.bot, update); err != nil {
			b.log.Error("failed to handle update: %v", err)
			handler.HandleError(b.bot, update, handler.InternalServerError)
			return
		}
		// Если нажали кнопку
	} else if update.CallbackQuery != nil {
		b.log.Info("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackData())

		var callback ViewFunc

		err, callbackView := b.CallbackStrings(update.CallbackData())
		if err != nil {
			b.log.Error("%v", err)
			return
		}

		callback = callbackView

		if err := callback(ctx, b.bot, update); err != nil {
			b.log.Error("failed to handle update: %v", err)
			handler.HandleError(b.bot, update, handler.InternalServerError)
			return
		}
		// Если пришла заявка на вступление
	} else if update.ChatJoinRequest != nil {
		b.log.Info("[%s] %s", update.ChatJoinRequest.From.UserName, update.ChatJoinRequest.InviteLink.InviteLink)

		if err := b.userService.CreateUser(ctx, userUpdateToModel(update)); err != nil {
			b.log.Error("userService.CreateUser: %v", err)
			handler.HandleError(b.bot, update, handler.InternalServerError)
			return
		}

		if err := b.requestService.CreateRequest(ctx, requestUpdateToModel(update)); err != nil {
			b.log.Error("requestService.CreateRequest: %v", err)
			handler.HandleError(b.bot, update, handler.InternalServerError)
			return
		}

		// Если добавляют/удаляют канал
	} else if update.MyChatMember != nil {
		b.log.Info("[%s] %s", update.MyChatMember.From.UserName, update.MyChatMember.NewChatMember.Status)

		if err := b.channelService.ChatMember(ctx, channelUpdateToModel(update)); err != nil {
			b.log.Error("channelService.ChatMember: %v", err)
			handler.HandleError(b.bot, update, handler.InternalServerError)
			return
		}
	}
}
