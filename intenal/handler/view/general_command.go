package view

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type viewGeneral struct {
	log *logger.Logger
}

func NewViewGeneral(log *logger.Logger) *viewGeneral {
	return &viewGeneral{
		log: log,
	}
}

func (v *viewGeneral) ViewStart() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "<b>Главное меню бота</b>")
		msg.ReplyMarkup = &markup.StartMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			v.log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}
