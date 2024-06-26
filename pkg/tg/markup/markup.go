package markup

import (
	"fmt"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	StartMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление каналами️", "channel_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление данными пользователей", "user_setting")),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("Управление дополнительными ботами для рассылки", "bot_spam_settings")),
	)

	InfoRequest = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять всех", "approved_all")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отклонить всех", "rejected_all")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Включить/Отключить капчу", "captcha_manager")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Настройка приветственного сообщения", "hello_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Настройка рассылки по базе", "sender_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять через установленное время", "approved_time")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назначить время принятия", "time_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статистика за день", "get_statistic")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "channel_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	UserSettingMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Рассылка по всей базе", "all_db_sender")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Настройка рассылки по всей базе", "global_setting_notification")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Скачать Excel файл с пользователя", "download_excel")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление ролями(ограниченный доступ)", "role_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	HelloMessageSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить сообщение", "add_text_notification"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить сообщение", "delete_text_notification")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить фотографию", "add_photo_notification"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить фотографию", "delete_photo_notification")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить кнопку", "add_button_notification"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить кнопку", "delete_button_notification")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить пример рассылки", "example_notification")),
		tgbotapi.NewInlineKeyboardRow(button.ComebackSetting),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	GlobalHelloMessageSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить сообщение", "global_add_text_notification"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить сообщение", "global_delete_text_notification")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить фотографию", "global_add_photo_notification"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить фотографию", "global_delete_photo_notification")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить кнопку", "global_add_button_notification"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить кнопку", "global_delete_button_notification")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить пример рассылки", "global_example_notification")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "user_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	CancelCommand = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена команды", "cancel_setting")),
	)

	CancelCommandSender = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена команды", "cancel_sender_setting")),
	)

	SenderMessageSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сделать рассылку", "send_message")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить сообщение", "update_sender_message"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить сообщение", "delete_sender_message")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить пример сообщения", "example_sender_message")),
		tgbotapi.NewInlineKeyboardRow(button.ComebackSetting),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	SuperAdminSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назначить администратором", "create_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назначить супер администратором", "create_super_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Забрать права администратора", "delete_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список администраторов", "all_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "user_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	CancelAdminCommand = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена команды", "cancel_admin_setting")),
	)

	SuperAdminComeback = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "role_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	BotSpamSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить бота", "add_spam_bot")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить бота", "delete_spam_bot")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список подключенных ботов", "list_spam_bot")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сделать рассылку", "activate_spam_bots")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	AllDbSender = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Настройка рассылки", "all_db_sender_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Настройка рассылки", "all_db_sender_setting")),
	)

	Captcha = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ДА", "press_captcha")))
)

func InfoRequestV2(time int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять всех", "approved_all")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отклонить всех", "rejected_all")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Включить/Отключить капчу", "captcha_manager")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Включить/Отключить опрос после капчи", "question_handbrake")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать/Изменить опрос после капчи", "question_manager")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить пример опроса после капчи", "question_example")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Настройка приветственного сообщения", "hello_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Настройка рассылки по базе", "sender_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Принять в установленное время: 16:00 и 10:00"), "approved_time")),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("Назначить время принятия", "time_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статистика за день", "get_statistic")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "channel_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)
}
