package middleware

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AdminMiddleware(channelID []int64, next tgbot.ViewFunc) tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		for _, chatID := range channelID {
			admins, err := bot.GetChatAdministrators(
				tgbotapi.ChatAdministratorsConfig{
					ChatConfig: tgbotapi.ChatConfig{
						ChatID: chatID,
					},
				})

			if err != nil {
				return err
			}

			for _, admin := range admins {
				if admin.User.ID == update.Message.From.ID {
					return next(ctx, bot, update)
				}
			}
		}
		return boterror.ErrIsNotAdmin
	}
}
