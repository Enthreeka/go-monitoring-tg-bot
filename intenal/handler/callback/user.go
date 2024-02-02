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
	"log"
	"os"
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

		fileName, err := c.Excel.GenerateExcelFile(users)
		if err != nil {
			c.Log.Error("Excel.GenerateExcelFile: failed to generate excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		file, err := os.Open(fileName)
		if err != nil {
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {

		}

		fileSize := fileInfo.Size()
		fileID := tgbotapi.FileBytes{Name: "file.txt", Bytes: make([]byte, fileSize)}
		if _, err = file.Read(fileID.Bytes); err != nil {

		}

		msg := tgbotapi.NewDocument(update.FromChat().ID, tgbotapi.FileBytes{Name: fileName, Bytes: fileID.Bytes})
		// Отправка сообщения
		_, err = bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}

		return nil
	}
}
