package callback

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/service"
	"github.com/Entreeka/monitoring-tg-bot/pkg/excel"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackUser struct {
	UserService service.UserService
	Log         *logger.Logger
	Excel       *excel.Excel
}

func (c *CallbackUser) CallbackGetExcelFile() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		users, err := c.UserService.GetAllUsers(ctx)
		if err != nil {
			c.Log.Error("userService.GetAllUsers: failed to get all users: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

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
		msg.Caption = userExcelFileText()

		if _, err = bot.Send(msg); err != nil {
			c.Log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	}
}
