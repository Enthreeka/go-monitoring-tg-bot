package callback

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
)

type CallbackChannel struct {
	ChannelService service.ChannelService
	RequestService service.RequestService
	UserService    service.UserService
	Log            *logger.Logger
	Store          *stateful.Store
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
			handler.MessageGetChannelInfo(channel.ChannelName, channel.WaitingCount, userCount, channel.NeedCaptcha, channel.QuestionEnabled))
		msg.ParseMode = tgbotapi.ModeHTML
		InfoRequestV2Mrk := markup.InfoRequestV2(channel.AcceptTimer)

		msg.ReplyMarkup = &InfoRequestV2Mrk

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
			handler.MessageGetChannelInfo(channel.ChannelName, channel.WaitingCount, userCount, channel.NeedCaptcha, channel.QuestionEnabled))
		msg.ParseMode = tgbotapi.ModeHTML
		InfoRequestV2Mrk := markup.InfoRequestV2(channel.AcceptTimer)

		msg.ReplyMarkup = &InfoRequestV2Mrk

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

func (c *CallbackChannel) CallbackQuestionHandbrake() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		c.Log.Info("", channelName)

		if err := c.ChannelService.UpdateQuestionEnabledByChannelName(ctx, channelName); err != nil {
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

func (c *CallbackChannel) CallbackGetQuestionExample() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)

		question, mrk, err := c.ChannelService.GetQuestion(ctx, channelName)
		if err != nil {
			c.Log.Error("ChannelService.GetQuestion: failed to get question: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, question)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = mrk

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}

func (c *CallbackChannel) CallbackQuestionManager() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		c.Log.Info("", channelName)
		var arg any

		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		channel, err := c.ChannelService.GetByChannelName(ctx, channelName)
		if err != nil {
			c.Log.Error("channelRepo.GetByChannelName: failed to query channel: %v", err)
			return err
		}

		if channel.Question == nil {
			exampleModel := entity.QuestionModel{
				Question: "Вопрос отсутствует",
				Answer: []entity.Answer{
					{
						ID:              0,
						AnswerVariation: "Варианта ответа отсутствует v1",
						Url:             "Ссылка отсутствует v1",
						TextResult:      "Текст после ответа отсутствует v1",
					},
					{
						ID:              0,
						AnswerVariation: "Варианта ответа отсутствует v2",
						Url:             "Ссылка отсутствует v2",
						TextResult:      "Текст после ответа отсутствует v2",
					},
				},
			}
			arg = exampleModel
		} else {
			arg = channel.Question
		}

		channel.Question, err = json.MarshalIndent(arg, "", "  ")
		if err != nil {
			c.Log.Error("channelRepo.GetByChannelName: failed to marshal exampleModel: %v", err)
			return err
		}

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.ChannelUpdateQuestion(string(channel.Question), channelName))
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Channel: &stateful.Channel{
				ChannelName:   channelName,
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationUpdateQuestion,
			},
		}, userID)
		c.Log.Info(channelName, msgSend.MessageID, stateful.OperationUpdateQuestion)

		return nil
	}
}

func (c *CallbackChannel) CallbackTimerSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		c.Log.Info("", channelName)

		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.ChannelSetTimer)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Channel: &stateful.Channel{
				ChannelName:   channelName,
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationSetTimer,
			},
		}, userID)
		c.Log.Info(channelName, msgSend.MessageID, stateful.OperationSetTimer)

		return nil
	}
}

func (c *CallbackChannel) CallbackGetAnswer() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelBase64, answerID := entity.ExtractValues(update.CallbackData())
		userID := update.FromChat().ID

		channelByte, err := base64.StdEncoding.DecodeString(channelBase64)
		if err != nil {
			c.Log.Error("failed to base64 decode channel answer", zap.Error(err))
			return nil
		}

		questionByte, err := c.ChannelService.GetQuestionByChannelName(ctx, string(channelByte))
		if err != nil {
			c.Log.Error("failed to query question", zap.Error(err))
			return nil
		}

		answerIDInt, err := strconv.Atoi(answerID)
		if err != nil {
			c.Log.Error("failed to convert answer to int", zap.Error(err))
			return nil
		}
		var model entity.QuestionModel

		if err := json.Unmarshal(questionByte, &model); err != nil {
			c.Log.Error("failed to unmarshal question", zap.Error(err))
			return nil
		}

		for _, answer := range model.Answer {
			if answer.ID == answerIDInt {

				msg := tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID, answer.TextResult)
				if answer.Url != "" && entity.IsValidURL(answer.Url) {
					urlButton := tgbotapi.NewInlineKeyboardButtonURL("ССЫЛКА", answer.Url)
					row := []tgbotapi.InlineKeyboardButton{urlButton}
					keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
					msg.ReplyMarkup = &keyboard
				}

				_, err := bot.Send(msg)
				if err != nil {
					c.Log.Error("failed to send message", zap.Error(err))
					return nil
				}

				return nil
			}
		}

		return nil
	}
}
