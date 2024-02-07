package callback

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/excel"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"sync"
)

type CallbackUser struct {
	UserService   service.UserService
	SenderService service.SenderService
	Log           *logger.Logger
	Excel         *excel.Excel
	Store         *stateful.Store

	mu sync.Mutex
}

func (c *CallbackUser) CallbackGetExcelFile() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		users, err := c.UserService.GetAllUsers(ctx)
		if err != nil {
			c.Log.Error("userService.GetAllUsers: failed to get all users: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		c.mu.Lock()
		fileName, err := c.Excel.GenerateExcelFile(users, update.CallbackQuery.From.UserName)
		if err != nil {
			c.Log.Error("Excel.GenerateExcelFile: failed to generate excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		fileIDBytes, err := c.Excel.GetExcelFile(fileName)
		if err != nil {
			c.Log.Error("Excel.GetExcelFile: failed to get excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}
		c.mu.Unlock()

		if fileIDBytes == nil {
			c.Log.Error("fileIDBytes: %v", boterror.ErrNil)
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrNil))
			return nil
		}

		msg := tgbotapi.NewDocument(update.FromChat().ID, tgbotapi.FileBytes{
			Name:  fileName,
			Bytes: *fileIDBytes,
		})
		msg.ParseMode = tgbotapi.ModeHTML
		msg.Caption = handler.UserExcelFileText()

		if _, err = bot.Send(msg); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}

func (c *CallbackUser) CallbackGetUserSenderSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		msg := tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID, handler.UserSenderSetting(channelName))
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = &markup.SenderMessageSetting

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackUser) CallbackPostMessageToUser() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		sender, err := c.SenderService.GetSender(ctx, channelName)
		if err != nil {
			c.Log.Error("SenderService.GetSender: failed to get sender: %v", err)
			if errors.Is(err, boterror.ErrNoRows) {
				if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.UserSenderEmpty)); err != nil {
					c.Log.Error("failed to send message", zap.Error(err))
					return err
				}
				return nil
			}
			return err
		}

		usersID, err := c.UserService.GetAllIDByChannelID(ctx, channelName)
		if err != nil {
			c.Log.Error("UserService.GetAllIDByChannelID: failed to get users id: %v", err)
			if errors.Is(err, boterror.ErrNoRows) {
				if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.UserSenderErrorEmpty)); err != nil {
					c.Log.Error("failed to send message", zap.Error(err))
					return err
				}
				return nil
			}
			return err
		}

		for _, id := range usersID {
			msg := tgbotapi.NewMessage(id, sender.Message)
			if _, err := bot.Send(msg); err != nil {
				c.Log.Error("failed to send message to user:%d err:%v", id, zap.Error(err))

				if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.UserSenderError)); err != nil {
					c.Log.Error("failed to send message", zap.Error(err))
					return err
				}
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.UserSenderDone)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}

func (c *CallbackUser) CallbackUpdateUserSenderMessage() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		msg := tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID, handler.UserUpdateSenderText)
		msg.ReplyMarkup = &markup.CancelCommandSender
		msg.ParseMode = tgbotapi.ModeHTML

		sendMsg, err := bot.Send(msg)
		if err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Sender: &stateful.Sender{
				ChannelName: channelName,
				MessageID:   sendMsg.MessageID,
			},
		}, userID)

		return nil
	}
}

func (c *CallbackUser) CallbackDeleteUserSenderMessage() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		if err := c.SenderService.DeleteSender(ctx, channelName); err != nil {
			c.Log.Error("SenderService.DeleteSender: failed to delete sender: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if _, err := bot.Send(tgbotapi.NewMessage(userID, handler.UserDeleteSenderText)); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackUser) CallbackGetExampleUserSenderMessage() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelName := findTitle(update.CallbackQuery.Message.Text)
		userID := update.FromChat().ID

		sender, err := c.SenderService.GetSender(ctx, channelName)
		if err != nil {
			c.Log.Error("SenderService.GetSender: failed to get sender: %v", err)
			if errors.Is(err, boterror.ErrNoRows) {
				if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, handler.UserSenderEmpty)); err != nil {
					c.Log.Error("failed to send message", zap.Error(err))
					return err
				}
				return nil
			}
			return err
		}

		msg := tgbotapi.NewMessage(userID, sender.Message)

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackUser) CallbackCancelSenderSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		data, _ := c.Store.Read(userID)
		c.Store.Delete(userID)

		channelName := data.Sender.ChannelName

		msg := tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID,
			handler.UserSenderSetting(channelName))
		msg.ReplyMarkup = &markup.SenderMessageSetting
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackUser) CallbackSuperAdminSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, handler.UserSuperAdminSetting)
		msg.ReplyMarkup = &markup.SuperAdminSetting
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}

func (c *CallbackUser) CallbackSetAdmin() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageID := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageID, handler.UserSetAdmin)
		msg.ReplyMarkup = &markup.CancelAdminCommand
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Admin: &stateful.Admin{
				OperationType: stateful.OperationAddAdmin,
				MessageID:     messageID,
			},
		}, userID)

		return nil
	}
}

func (c *CallbackUser) CallbackSetSuperAdmin() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageID := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageID, handler.UserSetSuperAdmin)
		msg.ReplyMarkup = &markup.CancelAdminCommand
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Admin: &stateful.Admin{
				OperationType: stateful.OperationAddSuperAdmin,
				MessageID:     messageID,
			},
		}, userID)

		return nil
	}
}

func (c *CallbackUser) CallbackDeleteAdmin() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		messageID := update.CallbackQuery.Message.MessageID

		msg := tgbotapi.NewEditMessageText(userID, messageID, handler.UserDeleteAdmin)
		msg.ReplyMarkup = &markup.CancelAdminCommand
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		c.Store.Delete(userID)
		c.Store.Set(&stateful.StoreData{
			Admin: &stateful.Admin{
				OperationType: stateful.OperationDeleteAdmin,
				MessageID:     messageID,
			},
		}, userID)

		return nil
	}
}

func (c *CallbackUser) CallbackGetAllAdmin() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		admin, err := c.UserService.GetAllAdmin(ctx)
		if err != nil {
			c.Log.Error("UserService.GetAllAdmin: failed to get admin: %v", err)
			return err
		}

		adminBytes, _ := json.MarshalIndent(admin, "", " ")

		msg := tgbotapi.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			string(adminBytes))
		msg.ReplyMarkup = &markup.SuperAdminComeback
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}

func (c *CallbackUser) CallbackCancelAdminSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		userID := update.FromChat().ID
		c.Store.Delete(userID)

		msg := tgbotapi.NewEditMessageText(userID, update.CallbackQuery.Message.MessageID,
			handler.UserSuperAdminSetting)
		msg.ReplyMarkup = &markup.SuperAdminSetting
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.Log.Error("failed to send message", zap.Error(err))
			return err
		}
		return nil
	}
}
