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
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type CallbackNotification struct {
	NotificationService service.NotificationService
	ChannelService      service.ChannelService
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

		userID := update.FromChat().ID

		sender := tg.NewSender(c.Log, notification, bot)

		if err := sender.SendMsgToNewUser(userID); err != nil {
			c.Log.Error("sender.SendMsgToNewUser: failed to send example notification: %v", err)
			return err
		}

		return nil
	}
}

func (c *CallbackNotification) CallbackCancelNotificationSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		data, _ := c.Store.Read(userID)
		c.Store.Delete(userID)

		var (
			msg tgbotapi.EditMessageTextConfig
		)

		switch {
		case data.Channel != nil && data.Channel.OperationType == stateful.OperationUpdateQuestion:
			channelMarkup, err := c.ChannelService.GetAllAdminChannel(ctx)
			if err != nil {
				c.Log.Error("channelService.GetAllAdminChannel: failed to get channel: %v", err)
				handler.HandleError(bot, update, boterror.ParseErrToText(err))
				return nil
			}

			msg = tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, handler.MessageShowAllChannel)
			msg.ParseMode = tgbotapi.ModeHTML
			msg.ReplyMarkup = channelMarkup

		case data.Notification != nil && data.Notification.ChannelName != "":
			channelName := data.Notification.ChannelName

			msg = tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID,
				handler.NotificationSettingText(channelName))
			msg.ReplyMarkup = &markup.HelloMessageSetting
			msg.ParseMode = tgbotapi.ModeHTML

		default:
			msg = tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID,
				handler.NotificationGlobalSetting)
			msg.ReplyMarkup = &markup.GlobalHelloMessageSetting
			msg.ParseMode = tgbotapi.ModeHTML
		}

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

// CallbackGetSettingGlobalNotification callback = global_setting_notification
func (c *CallbackNotification) CallbackGetSettingGlobalNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			handler.NotificationGlobalSetting)
		msg.ReplyMarkup = &markup.GlobalHelloMessageSetting
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

// CallbackGlobalUpdateTextNotification callback = global_add_text_notification
func (c *CallbackNotification) CallbackGlobalUpdateTextNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.NotificationGlobalUpdateText)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   "",
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationUpdateText,
			},
		}, userID)

		return nil
	}
}

// CallbackGlobalUpdateFileNotification callback = global_add_photo_notification
func (c *CallbackNotification) CallbackGlobalUpdateFileNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.NotificationGlobalUpdateFile)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   "",
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationUpdateFile,
			},
		}, userID)

		return nil
	}
}

// CallbackGlobalUpdateButtonNotification callback = global_add_button_notification
func (c *CallbackNotification) CallbackGlobalUpdateButtonNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageId := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageId, handler.NotificationGlobalUpdateButton)
		msg.ReplyMarkup = &markup.CancelCommand
		msg.ParseMode = tgbotapi.ModeHTML

		msgSend, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Notification: &stateful.Notification{
				ChannelName:   "",
				MessageID:     msgSend.MessageID,
				OperationType: stateful.OperationUpdateButton,
			},
		}, userID)

		return nil
	}
}

// CallbackGetGlobalExampleNotification callback = global_example_notification
func (c *CallbackNotification) CallbackGetGlobalExampleNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		notification, err := c.NotificationService.GetByChannelName(ctx, "")
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

		userID := update.FromChat().ID

		sender := tg.NewSender(c.Log, notification, bot)

		if err := sender.SendMsgToNewUser(userID); err != nil {
			c.Log.Error("sender.SendMsgToNewUser: failed to send example notification: %v", err)
			return err
		}

		return nil
	}
}

// CallbackGlobalDeleteButtonNotification callback = global_delete_button_notification
func (c *CallbackNotification) CallbackGlobalDeleteButtonNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID

		if err := c.NotificationService.UpdateButtonNotification(ctx, &entity.Notification{
			ChannelName: "",
			ButtonText:  nil,
			ButtonURL:   nil,
		}); err != nil {
			c.Log.Error("NotificationService.UpdateButtonNotification: failed to delete button in global notification: %v", err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.NotificationDeleteButton)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

// CallbackGlobalDeleteTextNotification callback = global_delete_text_notification
func (c *CallbackNotification) CallbackGlobalDeleteTextNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID

		if err := c.NotificationService.UpdateTextNotification(ctx, &entity.Notification{
			ChannelName:      "",
			NotificationText: nil,
		}); err != nil {
			c.Log.Error("NotificationService.UpdateTextNotification: failed to delete text in global notification: %v", err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.NotificationDeleteText)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

// CallbackGlobalDeleteFileNotification callback = global_delete_photo_notification
func (c *CallbackNotification) CallbackGlobalDeleteFileNotification() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID

		if err := c.NotificationService.UpdateFileNotification(ctx, &entity.Notification{
			ChannelName: "",
			FileID:      nil,
			FileType:    nil,
		}); err != nil {
			c.Log.Error("NotificationService.UpdateButtonNotification: failed to delete file in global notification: %v", err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.NotificationDeleteFile)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}
