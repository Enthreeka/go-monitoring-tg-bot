package callback

import (
	"context"
	"errors"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
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
			notificationSettingText(channelName))
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

		// Delete all past state and set new with stateful.OperationUpdateText
		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   channelName,
				OperationType: stateful.OperationUpdateText,
			},
		}, userID)

		msg := tgbotapi.NewMessage(userID, notificationUpdateText)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackNotification) CallbackUpdateFileNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		// Delete all past state and set new with stateful.OperationUpdateFile
		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   channelName,
				OperationType: stateful.OperationUpdateFile,
			},
		}, userID)

		msg := tgbotapi.NewMessage(userID, notificationUpdateFile)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackNotification) CallbackUpdateButtonNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		// Delete all past state and set new with stateful.OperationUpdateButton
		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   channelName,
				OperationType: stateful.OperationUpdateButton,
			},
		}, userID)

		msg := tgbotapi.NewMessage(userID, notificationUpdateButton)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackNotification) CallbackGetExampleNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		notification, err := c.NotificationService.GetByChannelName(ctx, channelName)
		if err != nil {
			if errors.Is(err, boterror.ErrNoRows) {
				if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, notificationEmpty)); err != nil {
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
			if _, err := bot.Send(tgbotapi.NewMessage(userID, notificationExampleError)); err != nil {
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
		c.Store.Delete(userID)

		msg := tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID, notificationCancel)
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

		if _, err := bot.Send(tgbotapi.NewMessage(userID, notificationDeleteButton)); err != nil {
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
			c.Log.Error("NotificationService.UpdateButtonNotification: failed to delete button: %v", err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, notificationDeleteText)); err != nil {
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
			c.Log.Error("NotificationService.UpdateButtonNotification: failed to delete button: %v", err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, notificationDeleteFile)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}
