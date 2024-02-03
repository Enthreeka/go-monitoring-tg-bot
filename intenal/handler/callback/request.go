package callback

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type CallbackRequest struct {
	RequestService service.RequestService
	Log            *logger.Logger
}

func (c *CallbackRequest) CallbackApproveAllRequest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)

		request, err := c.RequestService.GetAllByStatusRequest(ctx, tgbot.RequestInProgress, channelName)
		if err != nil {
			c.Log.Error("requestService.GetAllByStatusRequest: failed to get all in progress request: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		for _, req := range request {
			approveRequest := tgbotapi.ApproveChatJoinRequestConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: req.ChannelTelegramID,
				},
				UserID: req.UserID,
			}

			_, err := bot.Request(approveRequest)
			if err != nil {
				c.Log.Error("failed to approve requests: %v", err)
				return err
			}

			err = c.RequestService.UpdateStatusRequestByID(ctx, tgbot.RequestApproved, req.ID)
			if err != nil {
				c.Log.Error("RequestService.UpdateStatusRequestByID: failed to update status request:%s: %v", channelName, err)
				handler.HandleError(bot, update, boterror.ParseErrToText(err))
				return nil
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, requestApproved)); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}

func (c *CallbackRequest) CallbackRejectAllRequest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)

		request, err := c.RequestService.GetAllByStatusRequest(ctx, tgbot.RequestInProgress, channelName)
		if err != nil {
			c.Log.Error("requestService.GetAllByStatusRequest: failed to get all in progress request: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		for _, req := range request {
			declineRequest := tgbotapi.DeclineChatJoinRequest{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: req.ChannelTelegramID,
				},
				UserID: req.UserID,
			}

			_, err := bot.Request(declineRequest)
			if err != nil {
				c.Log.Error("failed to approve requests: %v", err)
				return err
			}

			err = c.RequestService.UpdateStatusRequestByID(ctx, tgbot.RequestRejected, req.ID)
			if err != nil {
				c.Log.Error("RequestService.UpdateStatusRequestByID: failed to update status request:%s: %v", channelName, err)
				handler.HandleError(bot, update, boterror.ParseErrToText(err))
				return nil
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, requestDecline)); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}

func (c *CallbackRequest) CallbackApproveAllThroughTime() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		go func(update *tgbotapi.Update) {
			seconds := 10
			time.Sleep(time.Duration(seconds) * time.Second)

			channelName := findTitle(update.CallbackQuery.Message.Text)

			request, err := c.RequestService.GetAllByStatusRequest(context.Background(), tgbot.RequestInProgress, channelName)
			if err != nil {
				c.Log.Error("requestService.GetAllByStatusRequest: failed to get all in progress request: %v", err)
				handler.HandleError(bot, update, boterror.ParseErrToText(err))
				return
			}

			for _, req := range request {
				approveRequest := tgbotapi.ApproveChatJoinRequestConfig{
					ChatConfig: tgbotapi.ChatConfig{
						ChatID: req.ChannelTelegramID,
					},
					UserID: req.UserID,
				}

				_, err := bot.Request(approveRequest)
				if err != nil {
					c.Log.Error("failed to approve requests: %v", err)
					return
				}

				err = c.RequestService.UpdateStatusRequestByID(context.Background(), tgbot.RequestApproved, req.ID)
				if err != nil {
					c.Log.Error("RequestService.UpdateStatusRequestByID: failed to update status request:%s: %v", channelName, err)
					handler.HandleError(bot, update, boterror.ParseErrToText(err))
					return
				}
			}

			if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, requestApproveThroughTime(seconds))); err != nil {
				c.Log.Error("failed to send msg: %v", err)
				return
			}
		}(update)

		return nil
	}
}

func findTitle(caption string) string { // переделать
	var (
		captionRune = []rune(caption)
		word        = "Канал:"
		wordRune    = []rune(word)
		wordRuneLen = len(wordRune)
		tempLen     = 0
		channelName = []rune("")
	)
	for i, el := range captionRune {
		if wordRuneLen != tempLen {
			if el == wordRune[0] {
				wordRune = wordRune[1:]
				tempLen++
			} else {
				wordRune = []rune(word)
				tempLen = 0
			}
		}

		if wordRuneLen == tempLen {
			if string(el) == " " {
				return string(channelName[:len(channelName)-1])
			}
			i += 1
			channelName = append(channelName, captionRune[i])
		}
	}

	return ""
}
