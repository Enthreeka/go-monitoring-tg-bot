package tgbot

import (
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	internalServerError = "internal server error"
)

func (b *Bot) HandleError(update *tgbotapi.Update, messageError string) error {
	if update == nil {
		return boterror.ErrNil
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messageError)
	_, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error("failed to send message: %v", err)
	}

	return err
}
