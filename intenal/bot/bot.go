package bot

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/config"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/callback"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/view"
	pgRepo "github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"os/signal"
	"syscall"
)

type Bot struct {
	psql *postgres.Postgres

	userService         service.UserService
	requestService      service.RequestService
	notificationService service.NotificationService
	channelService      service.ChannelService

	generalViewHandler view.ViewGeneral

	channelCallbackHandler callback.CallbackChannel
	generalCallbackHandler callback.CallbackGeneral
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) initServices(psql *postgres.Postgres, log *logger.Logger) {
	userRepo := pgRepo.NewUserRepo(psql)
	requestRepo := pgRepo.NewRequestRepo(psql)
	notificationRepo := pgRepo.NewNotificationRepo(psql)
	channelRepo := pgRepo.NewChannelRepo(psql)

	b.userService = service.NewUserService(userRepo, requestRepo, log)
	b.requestService = service.NewRequestService(requestRepo)
	b.notificationService = service.NewNotificationService(notificationRepo)
	b.channelService = service.NewChannelService(channelRepo, log)
}

func (b *Bot) initHandlers(log *logger.Logger) {
	b.generalViewHandler = view.ViewGeneral{
		Log: log,
	}
	b.channelCallbackHandler = callback.CallbackChannel{
		ChannelService: b.channelService,
		Log:            log,
	}
	b.generalCallbackHandler = callback.CallbackGeneral{
		Log: log,
	}
}

func (b *Bot) initialize(log *logger.Logger) {
	b.initServices(b.psql, log)
	b.initHandlers(log)
}

func (b *Bot) Run(log *logger.Logger, cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("failed to load token %v", err)
	}

	log.Info("Authorized on account %s", bot.Self.UserName)

	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}
	defer psql.Close()
	b.psql = psql

	b.initialize(log)

	newBot := tgbot.NewBot(bot, log, b.requestService, b.userService, b.channelService)

	newBot.RegisterCommandView("start", b.generalViewHandler.ViewStart())

	newBot.RegisterCommandCallback("main_menu", b.generalCallbackHandler.CallbackStart())
	newBot.RegisterCommandCallback("channel_setting", b.channelCallbackHandler.CallbackShowAllChannel())
	//newBot.RegisterCommandCallback("channel_get", b.channelCallbackHandler.CallbackShowAllChannel())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}
	return nil
}
