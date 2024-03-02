package balancer

import (
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync/atomic"
)

func (b *BotPool) sendMsgToNewUser(notification *entity.Notification, userID int64, bot *tgbotapi.BotAPI) error {
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
			b.log.Error("failed to send message user_id:%d: %v", userID, err)
			return err
		}
		atomic.AddInt64(&b.successCounter, 1)
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
			b.log.Error("failed to send message: %v", err)
			return err
		}

		atomic.AddInt64(&b.successCounter, 1)
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
			b.log.Error("failed to send message", err)
			return err
		}

		atomic.AddInt64(&b.successCounter, 1)
		return nil
	}

	return nil
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
