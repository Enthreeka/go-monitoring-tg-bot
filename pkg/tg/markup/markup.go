package markup

import (
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	StartMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление каналами️", "channel_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление данными пользователей", "user_setting")),
	)

	RequestSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Настройки принятия заявок", "request_setting")),
	)

	InfoRequest = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять всех", "approved_all"),
			tgbotapi.NewInlineKeyboardButtonData("Принять часть", "approved_part")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отклонить всем", "rejected_all")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять через: 600с", "approved_time")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "channel_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	UserSettingMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Скачать Excel файл с пользователя", "download_excel")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)
)
