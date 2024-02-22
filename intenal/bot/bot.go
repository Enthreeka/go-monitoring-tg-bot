package bot

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/config"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/callback"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/middleware"
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
	spammerCallbackHandler      callback.CallbackSpammer
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
	b.spammerCallbackHandler = callback.CallbackSpammer{
		Log: log,
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

	newBot.RegisterCommandView("start", middleware.AdminMiddleware(b.userService, b.generalViewHandler.ViewStart()))

	newBot.RegisterCommandCallback("main_menu", middleware.AdminMiddleware(b.userService, b.generalCallbackHandler.CallbackStart()))
	newBot.RegisterCommandCallback("channel_setting", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackShowAllChannel()))
	newBot.RegisterCommandCallback("channel_get", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackShowChannelInfo()))
	newBot.RegisterCommandCallback("user_setting", middleware.AdminMiddleware(b.userService, b.generalCallbackHandler.CallbackGetUserSettingMenu()))
	newBot.RegisterCommandCallback("download_excel", middleware.AdminMiddleware(b.userService, b.userCallbackHandler.CallbackGetExcelFile()))
	newBot.RegisterCommandCallback("approved_all", middleware.AdminMiddleware(b.userService, b.requestCallbackHandler.CallbackApproveAllRequest()))
	newBot.RegisterCommandCallback("rejected_all", middleware.AdminMiddleware(b.userService, b.requestCallbackHandler.CallbackRejectAllRequest()))
	newBot.RegisterCommandCallback("approved_time", middleware.AdminMiddleware(b.userService, b.requestCallbackHandler.CallbackApproveAllThroughTime()))
	newBot.RegisterCommandCallback("hello_setting", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGetSettingNotification()))
	newBot.RegisterCommandCallback("add_text_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackUpdateTextNotification()))
	newBot.RegisterCommandCallback("add_photo_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackUpdateFileNotification()))
	newBot.RegisterCommandCallback("add_button_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackUpdateButtonNotification()))
	newBot.RegisterCommandCallback("example_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGetExampleNotification()))
	newBot.RegisterCommandCallback("cancel_setting", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackCancelNotificationSetting()))
	newBot.RegisterCommandCallback("delete_text_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackDeleteTextNotification()))
	newBot.RegisterCommandCallback("delete_photo_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackDeleteFileNotification()))
	newBot.RegisterCommandCallback("delete_button_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackDeleteButtonNotification()))
	newBot.RegisterCommandCallback("sender_setting", middleware.AdminMiddleware(b.userService, b.userCallbackHandler.CallbackGetUserSenderSetting()))
	newBot.RegisterCommandCallback("send_message", middleware.AdminMiddleware(b.userService, b.userCallbackHandler.CallbackPostMessageToUser()))
	newBot.RegisterCommandCallback("update_sender_message", middleware.AdminMiddleware(b.userService, b.userCallbackHandler.CallbackUpdateUserSenderMessage()))
	newBot.RegisterCommandCallback("delete_sender_message", middleware.AdminMiddleware(b.userService, b.userCallbackHandler.CallbackDeleteUserSenderMessage()))
	newBot.RegisterCommandCallback("example_sender_message", middleware.AdminMiddleware(b.userService, b.userCallbackHandler.CallbackGetExampleUserSenderMessage()))
	newBot.RegisterCommandCallback("comeback", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackShowChannelInfoByName()))
	newBot.RegisterCommandCallback("cancel_sender_setting", middleware.AdminMiddleware(b.userService, b.userCallbackHandler.CallbackCancelSenderSetting()))

	newBot.RegisterCommandCallback("role_setting", middleware.SuperAdminMiddleware(b.userService, b.userCallbackHandler.CallbackSuperAdminSetting()))
	newBot.RegisterCommandCallback("create_admin", middleware.SuperAdminMiddleware(b.userService, b.userCallbackHandler.CallbackSetAdmin()))
	newBot.RegisterCommandCallback("create_super_admin", middleware.SuperAdminMiddleware(b.userService, b.userCallbackHandler.CallbackSetSuperAdmin()))
	newBot.RegisterCommandCallback("delete_admin", middleware.SuperAdminMiddleware(b.userService, b.userCallbackHandler.CallbackDeleteAdmin()))
	newBot.RegisterCommandCallback("all_admin", middleware.SuperAdminMiddleware(b.userService, b.userCallbackHandler.CallbackGetAllAdmin()))
	newBot.RegisterCommandCallback("cancel_admin_setting", middleware.SuperAdminMiddleware(b.userService, b.userCallbackHandler.CallbackCancelAdminSetting()))

	newBot.RegisterCommandCallback("bot_spam_settings", middleware.AdminMiddleware(b.userService, b.spammerCallbackHandler.CallbackBotSpammerSetting()))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}
	return nil
}
