package markup

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	StartMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление каналами️", "channel_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление данными пользователей", "user_setting")),
	)
)
