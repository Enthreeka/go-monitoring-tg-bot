package callback

import (
	"context"
	"errors"
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
)

type CallbackNotification struct {
	NotificationService service.NotificationService
	Log                 *logger.Logger
	Store               *stateful.Store
}

func (c *CallbackNotification) CallbackGetSettingNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			handler.NotificationSettingText(channelName))
		msg.ReplyMarkup = &markup.HelloMessageSetting
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackNotification) CallbackUpdateTextNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.NotificationUpdateText)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		// Delete all past state and set new with stateful.OperationUpdateText
		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   channelName,
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationUpdateText,
			},
		}, userID)

		return nil
	}
}

func (c *CallbackNotification) CallbackUpdateFileNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.NotificationUpdateFile)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		// Delete all past state and set new with stateful.OperationUpdateFile
		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   channelName,
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationUpdateFile,
			},
		}, userID)

		return nil
	}
}

func (c *CallbackNotification) CallbackUpdateButtonNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.NotificationUpdateButton)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		// Delete all past state and set new with stateful.OperationUpdateButton
		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   channelName,
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationUpdateButton,
			},
		}, userID)

		return nil
	}
}

func (c *CallbackNotification) CallbackGetExampleNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		notification, err := c.NotificationService.GetByChannelName(ctx, channelName)
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

		userID := update.FromChat().ID
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
		default:
			if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.NotificationExampleError)); err != nil {
				c.Log.Error("failed to send message", zap.Error(err))
				return err
			}
		}

		return nil
	}
}

func buttonQualifier(buttonURL *string, buttonText *string) *tgbotapi.InlineKeyboardMarkup {
	if buttonURL != nil && buttonText != nil {
		var (
			btnText string
			btnURL  string
		)

		btnText = *buttonText
		btnURL = *buttonURL

		button := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(btnText, btnURL)),
		)
		return &button
	}
	return nil
}

func (c *CallbackNotification) CallbackCancelNotificationSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		data, _ := c.Store.Read(userID)
		c.Store.Delete(userID)

		channelName := data.Notification.ChannelName

		msg := tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID,
			handler.NotificationSettingText(channelName))
		msg.ReplyMarkup = &markup.HelloMessageSetting
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackNotification) CallbackDeleteButtonNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		if err := c.NotificationService.UpdateButtonNotification(ctx, &entity.Notification{
			ChannelName: channelName,
			ButtonText:  nil,
			ButtonURL:   nil,
		}); err != nil {
			c.Log.Error("NotificationService.UpdateButtonNotification: failed to delete button: %v", err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.NotificationDeleteButton)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackNotification) CallbackDeleteTextNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		if err := c.NotificationService.UpdateTextNotification(ctx, &entity.Notification{
			ChannelName:      channelName,
			NotificationText: nil,
		}); err != nil {
			c.Log.Error("NotificationService.UpdateTextNotification: failed to delete text: %v", err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.NotificationDeleteText)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackNotification) CallbackDeleteFileNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		if err := c.NotificationService.UpdateFileNotification(ctx, &entity.Notification{
			ChannelName: channelName,
			FileID:      nil,
			FileType:    nil,
		}); err != nil {
			c.Log.Error("NotificationService.UpdateButtonNotification: failed to delete file: %v", err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.NotificationDeleteFile)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}
