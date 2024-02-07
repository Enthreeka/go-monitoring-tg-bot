package bot

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/config"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/callback"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/view"
	pgRepo "github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/excel"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/postgres"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"os/signal"
	"syscall"
)

type Bot struct {
	psql  *postgres.Postgres
	excel *excel.Excel
	store *stateful.Store

	userService         service.UserService
	requestService      service.RequestService
	notificationService service.NotificationService
	channelService      service.ChannelService
	senderService       service.SenderService

	generalViewHandler view.ViewGeneral

	channelCallbackHandler      callback.CallbackChannel
	generalCallbackHandler      callback.CallbackGeneral
	userCallbackHandler         callback.CallbackUser
	requestCallbackHandler      callback.CallbackRequest
	notificationCallbackHandler callback.CallbackNotification
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) initServices(psql *postgres.Postgres, log *logger.Logger) {
	userRepo := pgRepo.NewUserRepo(psql)
	requestRepo := pgRepo.NewRequestRepo(psql)
	notificationRepo := pgRepo.NewNotificationRepo(psql)
	channelRepo := pgRepo.NewChannelRepo(psql)
	senderRepo := pgRepo.NewSenderRepo(psql)

	b.userService = service.NewUserService(userRepo, requestRepo, channelRepo, log)
	b.requestService = service.NewRequestService(requestRepo, log)
	b.notificationService = service.NewNotificationService(notificationRepo, channelRepo, log)
	b.channelService = service.NewChannelService(channelRepo, log)
	b.senderService = service.NewSenderService(senderRepo, channelRepo)
}

func (b *Bot) initHandlers(log *logger.Logger) {
	b.generalViewHandler = view.ViewGeneral{
		Log: log,
	}
	b.channelCallbackHandler = callback.CallbackChannel{
		ChannelService: b.channelService,
		RequestService: b.requestService,
		Log:            log,
	}
	b.generalCallbackHandler = callback.CallbackGeneral{
		Log: log,
	}
	b.userCallbackHandler = callback.CallbackUser{
		UserService:   b.userService,
		SenderService: b.senderService,
		Log:           log,
		Excel:         b.excel,
		Store:         b.store,
	}
	b.requestCallbackHandler = callback.CallbackRequest{
		RequestService:      b.requestService,
		NotificationService: b.notificationService,
		Log:                 log,
	}
	b.notificationCallbackHandler = callback.CallbackNotification{
		NotificationService: b.notificationService,
		Log:                 log,
		Store:               b.store,
	}
}

func (b *Bot) initExcel(log *logger.Logger) {
	b.excel = excel.NewExcel(log)
}

func (b *Bot) initialize(log *logger.Logger) {
	b.initStore()
	b.initExcel(log)
	b.initServices(b.psql, log)
	b.initHandlers(log)
}

func (b *Bot) initStore() {
	b.store = stateful.NewStore()
}

func (b *Bot) Run(log *logger.Logger, cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	bot.Debug = false
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

	newBot := tgbot.NewBot(bot, log, b.store, b.requestService, b.userService, b.channelService, b.notificationService, b.senderService)

	newBot.RegisterCommandView("start", b.generalViewHandler.ViewStart())

	newBot.RegisterCommandCallback("main_menu", b.generalCallbackHandler.CallbackStart())
	newBot.RegisterCommandCallback("channel_setting", b.channelCallbackHandler.CallbackShowAllChannel())
	newBot.RegisterCommandCallback("channel_get", b.channelCallbackHandler.CallbackShowChannelInfo())
	newBot.RegisterCommandCallback("user_setting", b.generalCallbackHandler.CallbackGetUserSettingMenu())
	newBot.RegisterCommandCallback("download_excel", b.userCallbackHandler.CallbackGetExcelFile())
	newBot.RegisterCommandCallback("approved_all", b.requestCallbackHandler.CallbackApproveAllRequest())
	newBot.RegisterCommandCallback("rejected_all", b.requestCallbackHandler.CallbackRejectAllRequest())
	newBot.RegisterCommandCallback("approved_time", b.requestCallbackHandler.CallbackApproveAllThroughTime())
	newBot.RegisterCommandCallback("hello_setting", b.notificationCallbackHandler.CallbackGetSettingNotification())
	newBot.RegisterCommandCallback("add_text_notification", b.notificationCallbackHandler.CallbackUpdateTextNotification())
	newBot.RegisterCommandCallback("add_photo_notification", b.notificationCallbackHandler.CallbackUpdateFileNotification())
	newBot.RegisterCommandCallback("add_button_notification", b.notificationCallbackHandler.CallbackUpdateButtonNotification())
	newBot.RegisterCommandCallback("example_notification", b.notificationCallbackHandler.CallbackGetExampleNotification())
	newBot.RegisterCommandCallback("cancel_setting", b.notificationCallbackHandler.CallbackCancelNotificationSetting())
	newBot.RegisterCommandCallback("delete_text_notification", b.notificationCallbackHandler.CallbackDeleteTextNotification())
	newBot.RegisterCommandCallback("delete_photo_notification", b.notificationCallbackHandler.CallbackDeleteFileNotification())
	newBot.RegisterCommandCallback("delete_button_notification", b.notificationCallbackHandler.CallbackDeleteButtonNotification())
	newBot.RegisterCommandCallback("sender_setting", b.userCallbackHandler.CallbackGetUserSenderSetting())
	newBot.RegisterCommandCallback("send_message", b.userCallbackHandler.CallbackPostMessageToUser())
	newBot.RegisterCommandCallback("update_sender_message", b.userCallbackHandler.CallbackUpdateUserSenderMessage())
	newBot.RegisterCommandCallback("delete_sender_message", b.userCallbackHandler.CallbackDeleteUserSenderMessage())
	newBot.RegisterCommandCallback("example_sender_message", b.userCallbackHandler.CallbackGetExampleUserSenderMessage())
	newBot.RegisterCommandCallback("comeback", b.channelCallbackHandler.CallbackShowChannelInfoByName())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}
	return nil
}
