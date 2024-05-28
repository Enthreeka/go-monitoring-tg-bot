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
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/spam"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"os/signal"
	"syscall"
)

type Bot struct {
	psql           *postgres.Postgres
	excel          *excel.Excel
	store          *stateful.Store
	spammerStorage *spam.SpammerBots

	userService         service.UserService
	requestService      service.RequestService
	notificationService service.NotificationService
	channelService      service.ChannelService
	senderService       service.SenderService
	spamBotService      service.SpamBotService

	generalViewHandler view.ViewGeneral

	channelCallbackHandler      callback.CallbackChannel
	generalCallbackHandler      callback.CallbackGeneral
	userCallbackHandler         callback.CallbackUser
	requestCallbackHandler      callback.CallbackRequest
	notificationCallbackHandler callback.CallbackNotification
	spamBotCallbackHandler      callback.CallbackSpamBot
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
	spamBotRepo := pgRepo.NewSpamBotRepo(psql)

	b.userService = service.NewUserService(userRepo, requestRepo, channelRepo, log)
	b.requestService = service.NewRequestService(requestRepo, log)
	b.notificationService = service.NewNotificationService(notificationRepo, channelRepo, log)
	b.channelService = service.NewChannelService(channelRepo, log, b.store)
	b.senderService = service.NewSenderService(senderRepo, channelRepo)
	b.spamBotService = service.NewSpamBotService(userRepo, spamBotRepo, b.spammerStorage, log)
}

func (b *Bot) initHandlers(log *logger.Logger) {
	b.generalViewHandler = view.ViewGeneral{
		UserService:         b.userService,
		ChannelService:      b.channelService,
		NotificationService: b.notificationService,
		Log:                 log,
		Store:               b.store,
	}
	b.channelCallbackHandler = callback.CallbackChannel{
		ChannelService: b.channelService,
		RequestService: b.requestService,
		UserService:    b.userService,
		Log:            log,
		Store:          b.store,
	}
	b.generalCallbackHandler = callback.CallbackGeneral{
		NotificationService: b.notificationService,
		Log:                 log,
	}
	b.userCallbackHandler = callback.CallbackUser{
		UserService:         b.userService,
		SenderService:       b.senderService,
		NotificationService: b.notificationService,
		Log:                 log,
		Excel:               b.excel,
		Store:               b.store,
	}
	b.requestCallbackHandler = callback.CallbackRequest{
		RequestService:      b.requestService,
		NotificationService: b.notificationService,
		ChannelService:      b.channelService,
		Log:                 log,
		Store:               b.store,
	}
	b.notificationCallbackHandler = callback.CallbackNotification{
		ChannelService:      b.channelService,
		NotificationService: b.notificationService,
		Log:                 log,
		Store:               b.store,
	}
	b.spamBotCallbackHandler = callback.CallbackSpamBot{
		NotificationService: b.notificationService,
		UserService:         b.userService,
		SpammerStorage:      b.spammerStorage,
		SpamBot:             b.spamBotService,
		Store:               b.store,
		Log:                 log,
	}
}

func (b *Bot) initExcel(log *logger.Logger) {
	b.excel = excel.NewExcel(log)
}

func (b *Bot) initialize(ctx context.Context, log *logger.Logger) {
	b.initStore(ctx)
	b.initExcel(log)
	b.initServices(b.psql, log)
	b.initHandlers(log)
}

func (b *Bot) initStore(ctx context.Context) {
	b.store = stateful.NewStore()

	go b.store.Worker(ctx, 60)
}

func (b *Bot) Run(log *logger.Logger, cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("failed to load token %v", err)
	}
	bot.Debug = false

	log.Info("Authorized on account %s", bot.Self.UserName)

	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}
	defer psql.Close()
	b.psql = psql

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	b.initialize(ctx, log)

	newBot := tgbot.NewBot(bot, log, b.store, b.spammerStorage, b.requestService, b.userService, b.channelService, b.notificationService, b.senderService, b.spamBotService)

	newBot.RegisterCommandView("start", middleware.AdminMiddleware(b.userService, b.generalViewHandler.ViewStart()))
	newBot.RegisterCommandView("confirm", b.generalViewHandler.ViewConfirmCaptcha())

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

	newBot.RegisterCommandCallback("get_statistic", middleware.AdminMiddleware(b.userService, b.requestCallbackHandler.CallbackRequestStatisticForToday()))

	newBot.RegisterCommandCallback("all_db_sender", middleware.AdminMiddleware(b.userService, b.userCallbackHandler.CallbackAllUserSender()))

	newBot.RegisterCommandCallback("global_setting_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGetSettingGlobalNotification()))
	newBot.RegisterCommandCallback("global_add_text_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGlobalUpdateTextNotification()))
	newBot.RegisterCommandCallback("global_example_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGetGlobalExampleNotification()))
	newBot.RegisterCommandCallback("global_add_photo_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGlobalUpdateFileNotification()))
	newBot.RegisterCommandCallback("global_add_button_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGlobalUpdateButtonNotification()))
	newBot.RegisterCommandCallback("global_delete_photo_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGlobalDeleteFileNotification()))
	newBot.RegisterCommandCallback("global_delete_text_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGlobalDeleteTextNotification()))
	newBot.RegisterCommandCallback("global_delete_button_notification", middleware.AdminMiddleware(b.userService, b.notificationCallbackHandler.CallbackGlobalDeleteButtonNotification()))

	newBot.RegisterCommandCallback("captcha_manager", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackCaptchaManager()))
	newBot.RegisterCommandCallback("time_setting", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackTimerSetting()))
	newBot.RegisterCommandCallback("question_example", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackGetQuestionExample()))
	newBot.RegisterCommandCallback("question_manager", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackQuestionManager()))
	newBot.RegisterCommandCallback("answer", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackGetAnswer()))
	newBot.RegisterCommandCallback("question_handbrake", middleware.AdminMiddleware(b.userService, b.channelCallbackHandler.CallbackQuestionHandbrake()))

	//newBot.RegisterCommandCallback("press_captcha", b.generalCallbackHandler.CallbackConfirmCaptcha())
	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}
	return nil
}
