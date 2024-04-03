package callback

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
)

type CallbackGeneral struct {
	NotificationService service.NotificationService
	Log                 *logger.Logger
}

func (v *CallbackGeneral) CallbackStart() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, handler.GeneralMainBotMenu)
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
		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, handler.GeneralUserSettingMenu)
		msg.ReplyMarkup = &markup.UserSettingMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			v.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (v *CallbackGeneral) CallbackConfirmCaptcha() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		_, after, exist := strings.Cut(update.CallbackQuery.Message.Text, "к каналу:")
		if !exist {
			v.Log.Error("failed to cut channel in CallbackConfirmCaptcha")
			return nil
		}
		notification, err := v.NotificationService.GetByChannelName(ctx, after[1:len(after)-1])
		if err != nil {
			v.Log.Error("NotificationService.GetByChannelName: failed to get channel by name: %v", err)
			return nil
		}

		sender := tg.NewSender(v.Log, notification, bot)
		if err := sender.SendMsgToNewUser(update.FromChat().ID); err != nil {
			v.Log.Error("sender.SendMsgToNewUser: failed to send msg in CallbackConfirmCaptcha: %v", err)
			return nil
		}

		v.Log.Info("CallbackConfirmCaptcha: notification: %v, channel: %s", notification, after[1:len(after)-1])

		return nil
	}
}
