package callback

import (
	"context"
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
			c.Log.Error("NotificationService.GetByChannelName: failed to get channel: %v", err)
			return err
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, *notification.NotificationText)
		if notification.FileID != nil {
			//fileID := tgbotapi.FileID(*notification.FileID)
			//msg := tgbotapi.DocumentConfig{
			//	ParseMode: tgbotapi.ModeHTML,
			//	BaseFile: tgbotapi.BaseFile{
			//		BaseChat: tgbotapi.BaseChat{
			//			ChatID:      update.CallbackQuery.Message.Chat.ID,
			//			ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tg.MainMenuButton)),
			//		},
			//		File: fileID,
			//	},
			//}
		}

		if notification.ButtonURL != nil {
			msg.ReplyMarkup = notification.ButtonURL
		}

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
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
