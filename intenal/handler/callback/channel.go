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

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, messageShowAllChannel)
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

		channel.WaitingCount, err = c.RequestService.GetCountByStatusRequestAndChannelTgID(ctx, tgbot.RequestInProgress, channel.TelegramID)
		if err != nil {
			c.Log.Error("RequestService.GetCountByStatusRequestAndChannelTgID: failed to get count waiting people")
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			messageGetChannelInfo(channel.ChannelName, channel.WaitingCount))
		msg.ParseMode = tgbotapi.ModeHTML

		msg.ReplyMarkup = &markup.InfoRequest

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}