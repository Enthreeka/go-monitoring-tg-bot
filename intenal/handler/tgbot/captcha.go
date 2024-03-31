package tgbot

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *Bot) sendCaptcha(ctx context.Context, userID int64, channel string) error {
	msg := tgbotapi.NewMessage(userID, handler.BotCaptcha(channel))
	//msg.ReplyMarkup = &markup.Captcha

	if _, err := b.bot.Send(msg); err != nil {
		b.log.Error("failed to send message", zap.Error(err))
		return err
	}
	return nil
}
