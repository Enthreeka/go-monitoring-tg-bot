package view

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

type ViewGeneral struct {
	UserService         service.UserService
	ChannelService      service.ChannelService
	NotificationService service.NotificationService
	Log                 *logger.Logger
	Store               *stateful.Store
}

func (v *ViewGeneral) ViewStart() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "<b>Главное меню бота</b>")
		msg.ReplyMarkup = &markup.StartMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			v.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}

func (v *ViewGeneral) ViewConfirmCaptcha() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		isPassed, err := v.UserService.IsPassedCaptchaByUserID(ctx, userID)
		if err != nil {
			v.Log.Error("UserService.IsPassedCaptchaByUserID", zap.Error(err))
			return nil
		}

		if !isPassed {
			if err := v.UserService.UpdateIsPassedCaptcha(ctx, true, userID); err != nil {
				v.Log.Error("UserService.UpdateIsPassedCaptcha", zap.Error(err))
				return nil
			}
		}

		channelName, exist := v.Store.ReadCaptcha(userID)
		if !exist {
			v.Log.Error("failed to get channel name in Store.ReadCaptcha")
			return nil
		}
		defer v.Store.DeleteCaptcha(userID)

		v.Log.Info("getting data from store with captcha: %v", channelName)

		var channel *entity.Channel
		channel, err = v.ChannelService.GetByChannelName(ctx, channelName.ChannelName)
		if err != nil {
			v.Log.Error("ChannelService.GetByChannelName", zap.Error(err))
			if channel == nil {
				channel = &entity.Channel{}
			}

			channel.ChannelName, err = v.ChannelService.GetChannelByUserID(ctx, userID)
			if err != nil {
				v.Log.Error("ChannelService.GetChannelByUserID", zap.Error(err))
				return nil
			}
		}

		noritifcation, err := v.NotificationService.GetByChannelName(ctx, channel.ChannelName)
		if err != nil || noritifcation == nil {
			v.Log.Error("NotificationService.GetByChannelName", zap.Error(err), channel.ChannelName)
			return nil
		}

		sender := tg.NewSender(v.Log, noritifcation, bot)
		err = sender.SendMsgToNewUser(userID)
		if err != nil {
			v.Log.Error("sender.SendMsgToNewUser", zap.Error(err))
			return nil
		}

		if channel.QuestionEnabled && channel.Question != nil {
			time.AfterFunc(10*time.Second, func() {
				question, mrk, err := v.ChannelService.GetQuestion(context.Background(), channel.ChannelName)
				if err != nil {
					v.Log.Error("failed to get question", zap.Error(err))
					return
				}

				newMsg := tgbotapi.NewMessage(userID, question)
				newMsg.ParseMode = tgbotapi.ModeHTML
				newMsg.ReplyMarkup = mrk

				if _, err := bot.Send(newMsg); err != nil {
					v.Log.Error("failed to send msg: %v", err)
					return
				}
			})
		}

		return nil
	}
}
