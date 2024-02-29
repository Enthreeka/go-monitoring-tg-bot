package callback

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type CallbackSpamBot struct {
	SpamBot service.SpamBotService
	Log     *logger.Logger
	Store   *stateful.Store
}

const (
	AddBot    = "add"
	DeleteBot = "delete"
	GetBot    = "get"
)

func (c *CallbackSpamBot) CallbackBotSpammerSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "TODO")
		msg.ReplyMarkup = &markup.BotSpamSetting
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackSpamBot) CallbackAddBotSpammer() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		_, err := c.SpamBot.GetAllBots(ctx, AddBot)
		if err != nil {
			c.Log.Error("SpamBot.GetAllBots: failed to get bots: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.SpamBotAdd)
		//msg.ReplyMarkup = botsMarkup
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		// Delete all past state and set new with stateful.OperationAddBot
		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			SpamBot: &stateful.SpamBot{
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationAddBot,
			}}, userID)

		return nil
	}
}

func (c *CallbackSpamBot) CallbackDeleteBotSpammer() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		botsMarkup, err := c.SpamBot.GetAllBots(ctx, DeleteBot)
		if err != nil {
			c.Log.Error("SpamBot.GetAllBots: failed to get bots: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.SpamBotDelete)
		msg.ReplyMarkup = botsMarkup
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		// Delete all past state and set new with stateful.OperationDeleteBot
		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			SpamBot: &stateful.SpamBot{
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationDeleteBot,
			}}, userID)

		return nil
	}
}

func (c *CallbackSpamBot) CallbackShowAllBotSpammer() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		botsMarkup, err := c.SpamBot.GetAllBots(ctx, GetBot)
		if err != nil {
			c.Log.Error("SpamBot.GetAllBots: failed to get bots: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.SpamBotGet)
		msg.ReplyMarkup = botsMarkup
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}
