package excel

import (
	"fmt"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

const filename = "users.xlsx"

type Excel struct {
	log *logger.Logger
	mu  sync.Mutex
}

func NewExcel(log *logger.Logger) *Excel {
	return &Excel{log: log}
}

func (e *Excel) GenerateExcelFile(users []entity.User, username string) (string, error) {
	start := time.Now()

	e.mu.Lock()
	f := excelize.NewFile()

	defer func() {
		e.mu.Unlock()
		if err := f.Close(); err != nil {
			e.log.Error("failed to close excel: %v", err)
		}
	}()

	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	headers := map[string]string{
		"A1": "ID",
		"B1": "UsernameTg",
		"C1": "Phone",
		"D1": "ChannelFrom",
		"E1": "CreatedAt",
		"F1": "Role",
	}

	for cell, value := range headers {
		f.SetCellValue(sheetName, cell, value)
	}

	for i, user := range users {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), user.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), user.UsernameTg)

		if user.Phone != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), *user.Phone)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "-")
		}

		if user.ChannelFrom != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), *user.ChannelFrom)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "-")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), user.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), user.Role)
	}

	err := f.SaveAs(filename)
	if err != nil {
		e.log.Error("failed to save file: %s", filename)
		return "", err
	}

	end := time.Since(start)
	e.log.Info("[%s] by [%s] Время генерации файла: %f", filename, username, end.Seconds())
	return filename, nil
}

func (e *Excel) GetExcelFile(fileName string) (*[]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		e.log.Error("os.Open: failed to open file: %v", err)
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		e.log.Error("file.Stat: failed to get file stat: %v", err)
		return nil, err
	}

	fileSize := fileInfo.Size()
	fileID := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: make([]byte, fileSize),
	}

	if _, err = file.Read(fileID.Bytes); err != nil {
		e.log.Error("file.Read: failed to get read file: %v", err)
		return nil, err
	}

	return &fileID.Bytes, nil
}
