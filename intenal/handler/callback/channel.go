package callback

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackChannel struct {
	ChannelService service.ChannelService
	Log            *logger.Logger
}

//func NewViewChannel(channelService service.ChannelService, log *logger.Logger) *ViewChannel {
//	return &ViewChannel{
//		channelService: channelService,
//		log:            log,
//	}
//}

func (v *CallbackChannel) CallbackShowAllChannel() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelMarkup, err := v.ChannelService.GetAllAdminChannel(ctx)
		if err != nil {
			v.Log.Error("channelService.GetAllAdminChannel: failed to get channel: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, messageShowAllChannel)
		msg.ParseMode = tgbotapi.ModeHTML

		msg.ReplyMarkup = channelMarkup

		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
