package tg

import (
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync/atomic"
)

type Sender interface {
	SendMsgToNewUser(userID int64) error
	GetSuccessCounter() int64
}

type sender struct {
	log          *logger.Logger
	notification *entity.Notification
	bot          *tgbotapi.BotAPI

	successCounter int64
}

func NewSender(log *logger.Logger, notification *entity.Notification, bot *tgbotapi.BotAPI) *sender {
	return &sender{
		log:          log,
		notification: notification,
		bot:          bot,
	}
}

func (b *sender) SendMsgToNewUser(userID int64) error {
	var isPhoto bool
	if b.notification.FileType != nil {
		if *b.notification.FileType == "photo" {
			isPhoto = true
		}
	}

	switch {
	case b.notification.FileType == nil && b.notification.NotificationText != nil:
		msg := tgbotapi.NewMessage(userID, "")
		buttonMarkup := b.buttonQualifier(b.notification.ButtonURL, b.notification.ButtonText)
		if buttonMarkup != nil {
			msg.ReplyMarkup = &buttonMarkup
		}
		if b.notification.NotificationText != nil {
			msg.Text = *b.notification.NotificationText
		}

		if _, err := b.bot.Send(msg); err != nil {
			b.log.Error("failed to send message user_id:%d: %v", userID, err)
			return err
		}
		atomic.AddInt64(&b.successCounter, 1)
		return nil

	case isPhoto && b.notification.FileType != nil:
		notificationPhoto := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(*b.notification.FileID))
		msg := tgbotapi.NewPhoto(userID, notificationPhoto.Media)
		buttonMarkup := b.buttonQualifier(b.notification.ButtonURL, b.notification.ButtonText)
		if buttonMarkup != nil {
			msg.ReplyMarkup = &buttonMarkup
		}
		if b.notification.NotificationText != nil {
			msg.Caption = *b.notification.NotificationText
		}

		if _, err := b.bot.Send(msg); err != nil {
			b.log.Error("failed to send message: %v", err)
			return err
		}

		atomic.AddInt64(&b.successCounter, 1)
		return nil

	case !isPhoto && b.notification.FileType != nil:
		msg := tgbotapi.DocumentConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID: userID,
				},
				File: tgbotapi.FileID(*b.notification.FileID),
			},
		}
		buttonMarkup := b.buttonQualifier(b.notification.ButtonURL, b.notification.ButtonText)
		if buttonMarkup != nil {
			msg.ReplyMarkup = &buttonMarkup
		}
		if b.notification.NotificationText != nil {
			msg.Caption = *b.notification.NotificationText
		}

		if _, err := b.bot.Send(msg); err != nil {
			b.log.Error("failed to send message", err)
			return err
		}

		atomic.AddInt64(&b.successCounter, 1)
		return nil
	}

	return nil
}

func (s *sender) buttonQualifier(buttonURL *string, buttonText *string) *tgbotapi.InlineKeyboardMarkup {
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

func (s *sender) GetSuccessCounter() int64 {
	return s.successCounter
}
