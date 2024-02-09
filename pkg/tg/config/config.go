package config

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	StartConfigMenu = tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "/start",
			Description: "Начать общение с ботом",
		},
	)
)
