package callback

import (
	"context"
	"errors"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

type CallbackRequest struct {
	RequestService      service.RequestService
	NotificationService service.NotificationService
	ChannelService      service.ChannelService
	Log                 *logger.Logger
	Store               *stateful.Store
}

func (c *CallbackRequest) CallbackApproveAllRequest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		var (
			channelName   = findTitle(update.CallbackQuery.Message.Text)
			countErr      int
			countApproved int
		)

		request, err := c.RequestService.GetAllByStatusRequest(ctx, tgbot.RequestInProgress, channelName)
		if err != nil {
			if errors.Is(err, boterror.ErrNoRows) {
				if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestEmpty)); err != nil {
					c.Log.Error("failed to send message", zap.Error(err))
					return err
				}
				return nil
			}
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

			if _, err := bot.Request(approveRequest); err != nil {
				c.Log.Error("failed to approve requests: %v, request: %v", err, req)
				countErr++

				if err = c.RequestService.UpdateStatusRequestByID(ctx, tgbot.RequestRejected, req.ID); err != nil {
					c.Log.Error("RequestService.UpdateStatusRequestByID: failed to update status request:%s: %v, request:%v",
						channelName, err, req)
				}
				continue
			}

			if err = c.RequestService.UpdateStatusRequestByID(ctx, tgbot.RequestApproved, req.ID); err != nil {
				c.Log.Error("RequestService.UpdateStatusRequestByID: failed to update status request:%s: %v, request:%v",
					channelName, err, req)
			}
			countApproved++
		}

		if countErr > 0 {
			if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestError(countErr))); err != nil {
				c.Log.Error("failed to send msg: %v", err)
				return err
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestApprovedText(countApproved))); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}

func (c *CallbackRequest) sendMsgToNewUser(ctx context.Context, userID int64, channelID int64, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	notification, err := c.NotificationService.GetByChannelTelegramID(ctx, channelID)
	if err != nil {
		if errors.Is(err, boterror.ErrNoRows) {
			if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.NotificationEmpty)); err != nil {
				c.Log.Error("failed to send message", zap.Error(err))
				return err
			}
			return nil
		}
		c.Log.Error("NotificationService.GetByChannelName: failed to get channel: %v", err)
		return err
	}
	var isPhoto bool
	if notification.FileType != nil {
		if *notification.FileType == "photo" {
			isPhoto = true
		}
	}

	switch {
	case notification.FileType == nil && notification.NotificationText != nil:
		msg := tgbotapi.NewMessage(userID, "")
		buttonMarkup := buttonQualifier(notification.ButtonURL, notification.ButtonText)
		if buttonMarkup != nil {
			msg.ReplyMarkup = &buttonMarkup
		}
		if notification.NotificationText != nil {
			msg.Text = *notification.NotificationText
		}

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil

	case isPhoto && notification.FileType != nil:
		notificationPhoto := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(*notification.FileID))
		msg := tgbotapi.NewPhoto(userID, notificationPhoto.Media)
		buttonMarkup := buttonQualifier(notification.ButtonURL, notification.ButtonText)
		if buttonMarkup != nil {
			msg.ReplyMarkup = &buttonMarkup
		}
		if notification.NotificationText != nil {
			msg.Caption = *notification.NotificationText
		}

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil

	case !isPhoto && notification.FileType != nil:
		msg := tgbotapi.DocumentConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID: userID,
				},
				File: tgbotapi.FileID(*notification.FileID),
			},
		}
		buttonMarkup := buttonQualifier(notification.ButtonURL, notification.ButtonText)
		if buttonMarkup != nil {
			msg.ReplyMarkup = &buttonMarkup
		}
		if notification.NotificationText != nil {
			msg.Caption = *notification.NotificationText
		}

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}

	return nil
}

func (c *CallbackRequest) CallbackRejectAllRequest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		var (
			channelName   = findTitle(update.CallbackQuery.Message.Text)
			countErr      int
			countRejected int
		)

		request, err := c.RequestService.GetAllByStatusRequest(ctx, tgbot.RequestInProgress, channelName)
		if err != nil {
			if errors.Is(err, boterror.ErrNoRows) {
				if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestEmpty)); err != nil {
					c.Log.Error("failed to send message", zap.Error(err))
					return err
				}
				return nil
			}
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

			if _, err := bot.Request(declineRequest); err != nil {
				c.Log.Error("failed to reject requests: %v", err)
				countErr++
			} else {
				countRejected++
			}

			if err = c.RequestService.UpdateStatusRequestByID(ctx, tgbot.RequestRejected, req.ID); err != nil {
				c.Log.Error("RequestService.UpdateStatusRequestByID: failed to update status request:%s: %v, request:%v",
					channelName, err, req)
			}
		}

		if countErr > 0 {
			if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestError(countErr))); err != nil {
				c.Log.Error("failed to send msg: %v", err)
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestDeclineText(countRejected))); err != nil {
			c.Log.Error("failed to send msg: %v", err)
		}
		return nil
	}
}

func (c *CallbackRequest) CallbackApproveAllThroughTime() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		go func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
			seconds := 600
			time.Sleep(time.Duration(seconds) * time.Second)

			var (
				channelName   = findTitle(update.CallbackQuery.Message.Text)
				countErr      int
				countApproved int
			)

			request, err := c.RequestService.GetAllByStatusRequest(context.Background(), tgbot.RequestInProgress, channelName)
			if err != nil {
				if errors.Is(err, boterror.ErrNoRows) {
					if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestEmpty)); err != nil {
						c.Log.Error("failed to send message", zap.Error(err))
						return
					}
					return
				}
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

				if _, err := bot.Request(approveRequest); err != nil {
					c.Log.Error("failed to approve requests: %v, request: %v", err, req)
					countErr++

					if err = c.RequestService.UpdateStatusRequestByID(ctx, tgbot.RequestRejected, req.ID); err != nil {
						c.Log.Error("RequestService.UpdateStatusRequestByID: failed to update status request:%s: %v, request:%v",
							channelName, err, req)
					}

					continue
				}

				if err = c.RequestService.UpdateStatusRequestByID(context.Background(), tgbot.RequestApproved, req.ID); err != nil {
					c.Log.Error("RequestService.UpdateStatusRequestByID: failed to update status request:%s: %v", channelName, err)
				}
				countApproved++
			}

			if countErr > 0 {
				if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestError(countErr))); err != nil {
					c.Log.Error("failed to send msg: %v", err)
				}
			}

			if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestApproveThroughTime(seconds, countApproved))); err != nil {
				c.Log.Error("failed to send msg: %v", err)
				return
			}
		}(update, bot)

		return nil
	}
}

func findTitle(caption string) string {
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
			if string(el) == "\n" {
				return string(channelName[:len(channelName)-2])
			}
			i += 1
			channelName = append(channelName, captionRune[i])
		}
	}

	return ""
}

func (c *CallbackRequest) CallbackRequestStatisticForToday() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)

		channel, err := c.ChannelService.GetByChannelName(ctx, channelName)
		if err != nil {
			c.Log.Error("ChannelService.GetByChannelName: failed to get channel by name: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		countRequest, err := c.RequestService.GetCountRequestTodayByChannelID(ctx, channel.ID)
		if err != nil {
			c.Log.Error("RequestService.GetCountRequestToday: failed to get count for today: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		day, countSentMsg := c.Store.GetSuccessfulSentMsg(channel.TelegramID)

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.RequestStatistic(day, countRequest,
			countSentMsg, channelName))); err != nil {
			c.Log.Error("failed to send msg: %v", err)
		}

		return nil
	}
}
