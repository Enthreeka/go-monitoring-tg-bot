package callback

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackChannel struct {
	ChannelService service.ChannelService
	RequestService service.RequestService
	UserService    service.UserService
	Log            *logger.Logger
}

func (c *CallbackChannel) CallbackShowAllChannel() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelMarkup, err := c.ChannelService.GetAllAdminChannel(ctx)
		if err != nil {
			c.Log.Error("channelService.GetAllAdminChannel: failed to get channel: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, handler.MessageShowAllChannel)
		msg.ParseMode = tgbotapi.ModeHTML

		msg.ReplyMarkup = channelMarkup

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}

func (c *CallbackChannel) CallbackShowChannelInfo() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelID := entity.GetID(update.CallbackData())
		if channelID == 0 {
			c.Log.Error("entity.GetID: failed to get id from channel button")
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrNotFoundID))
			return nil
		}

		channel, err := c.ChannelService.GetByID(ctx, channelID)
		if err != nil {
			c.Log.Error("ChannelService.GetByID: failed to get channel")
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		userCount, err := c.UserService.GetCountUserByChannelTgID(ctx, channel.TelegramID)
		if err != nil {
			c.Log.Error("UserService.GetCountUserByChannelTgID: failed to get count user in channel: %s", channel.ChannelName)
		}

		channel.WaitingCount, err = c.RequestService.GetCountByStatusRequestAndChannelTgID(ctx, tgbot.RequestInProgress, channel.TelegramID)
		if err != nil {
			c.Log.Error("RequestService.GetCountByStatusRequestAndChannelTgID: failed to get count waiting people")
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			handler.MessageGetChannelInfo(channel.ChannelName, channel.WaitingCount, userCount, channel.NeedCaptcha))
		msg.ParseMode = tgbotapi.ModeHTML

		msg.ReplyMarkup = &markup.InfoRequest

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}

func (c *CallbackChannel) CallbackShowChannelInfoByName() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		c.Log.Info("", channelName)
		channel, err := c.ChannelService.GetByChannelName(ctx, channelName)
		if err != nil {
			c.Log.Error("ChannelService.GetByChannelName: failed to get channel by name: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		userCount, err := c.UserService.GetCountUserByChannelTgID(ctx, channel.TelegramID)
		if err != nil {
			c.Log.Error("UserService.GetCountUserByChannelTgID: failed to get count user in channel: %s", channel.ChannelName)
		}

		channel.WaitingCount, err = c.RequestService.GetCountByStatusRequestAndChannelTgID(ctx, tgbot.RequestInProgress, channel.TelegramID)
		if err != nil {
			c.Log.Error("RequestService.GetCountByStatusRequestAndChannelTgID: failed to get count waiting people:%v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			handler.MessageGetChannelInfo(channel.ChannelName, channel.WaitingCount, userCount, channel.NeedCaptcha))
		msg.ParseMode = tgbotapi.ModeHTML

		msg.ReplyMarkup = &markup.InfoRequest

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}

func (c *CallbackChannel) CallbackCaptchaManager() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		c.Log.Info("", channelName)

		if err := c.ChannelService.UpdateNeedCaptchaByChannelName(ctx, channelName); err != nil {
			c.Log.Error("ChannelService.UpdateNeedCaptchaByChannelName: : %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		callback := c.CallbackShowChannelInfoByName()
		if err := callback(ctx, bot, update); err != nil {
			c.Log.Error("failed to process callback in CallbackCaptchaManager: %v", err)
			return err
		}

		return nil

	}
}
