package button

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Вернуться в главное меню", "main_menu")

	AddChannelButton = tgbotapi.NewInlineKeyboardButtonData("Добавить канал", "add_channel")

	ComebackSetting = tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "comeback")
)
