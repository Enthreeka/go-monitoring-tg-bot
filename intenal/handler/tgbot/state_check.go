package tgbot

import (
	"context"
	"errors"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrOperationType  = errors.New("operation type not found")
	ErrEmptyStoreData = errors.New("all values is nil")
)

const (
	success = "Успешно добавлено"
)

func (b *Bot) isStateExist(userID int64) (*stateful.StoreData, bool) {
	data, exist := b.store.Read(userID)
	return data, exist
}

func (b *Bot) getState(ctx context.Context, update *tgbotapi.Message) (bool, error) {
	storeData, isExist := b.isStateExist(update.From.ID)
	if isExist {

		typeData := getStoreData(storeData)
		if typeData == nil {
			b.log.Error("failed to get data: typeData == nil")
			return true, ErrEmptyStoreData
		}

		switch typeData {
		case typeData.(*stateful.Notification):
			if err := b.storeDataNotificationOperationType(ctx, storeData, update); err != nil {
				b.log.Error("storeDataNotificationOperationType: %v", err)
				return true, err
			}
			return true, nil
		case typeData.(*stateful.Channel):

		}

		return true, nil
	}
	return false, nil
}

func (b *Bot) storeDataNotificationOperationType(ctx context.Context, storeData *stateful.StoreData, update *tgbotapi.Message) error {
	userID := update.From.ID

	switch storeData.Notification.OperationType {
	case stateful.OperationUpdateText:
		defer b.store.Delete(userID)
		notification := &entity.Notification{
			NotificationText: &update.Text,
			FileID:           nil,
			FileType:         nil,
			ButtonURL:        nil,
			ChannelName:      storeData.Notification.ChannelName,
		}
		err := b.notificationService.UpdateTextNotification(ctx, notification)
		if err != nil {
			b.log.Error("notificationService.UpdateTextNotification: failed to work with notification: %v", err)
			return err
		}

		if _, err := b.bot.Send(tgbotapi.NewMessage(userID, success)); err != nil {
			b.log.Error("failed to send msg: %v", err)
			return err
		}
		return nil
	case stateful.OperationUpdateFile:
	case stateful.OperationUpdateButton:
	}

	return ErrOperationType
}

func getStoreData(storeData *stateful.StoreData) any {
	switch {
	case storeData.Notification != nil:
		return storeData.Notification
	case storeData.Channel != nil:
		return storeData.Channel
	}
	return nil
}
