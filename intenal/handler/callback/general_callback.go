package callback

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type CallbackGeneral struct {
	Log *logger.Logger
}

func (v *CallbackGeneral) CallbackStart() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, generalMainBotMenu)
		msg.ReplyMarkup = &markup.StartMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			v.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (v *CallbackGeneral) CallbackGetUserSettingMenu() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, generalUserSettingMenu)
		msg.ReplyMarkup = &markup.UserSettingMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			v.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}
