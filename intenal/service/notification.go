package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
)

type NotificationService interface {
	createNotificationIfNotExist(ctx context.Context, notification *entity.Notification) (int64, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.Notification, error)
	GetByChannelName(ctx context.Context, channelName string) (*entity.Notification, error)
	UpdateTextNotification(ctx context.Context, notification *entity.Notification) error
	UpdateFileNotification(ctx context.Context, notification *entity.Notification) error
	UpdateButtonNotification(ctx context.Context, notification *entity.Notification) error
	GetByChannelTelegramID(ctx context.Context, channelTelegramID int64) (*entity.Notification, error)
}

type notificationService struct {
	notificationRepo postgres.NotificationRepo
	channelRepo      postgres.ChannelRepo
	log              *logger.Logger
}

func NewNotificationService(notificationRepo postgres.NotificationRepo, channelRepo postgres.ChannelRepo, log *logger.Logger) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		channelRepo:      channelRepo,
		log:              log,
	}
}

func (n *notificationService) createNotificationIfNotExist(ctx context.Context, notification *entity.Notification) (int64, error) {
	n.log.Info("Get notification: %s", notification.String())

	channelID, err := n.channelRepo.GetChannelIDByChannelName(ctx, notification.ChannelName)
	if err != nil {
		n.log.Error("channelRepo.GetChannelIDByChannelName: failed to get channel id: %v", err)
		return 0, err
	}

	isExist, err := n.notificationRepo.IsExistNotificationByChannelID(ctx, channelID)
	if err != nil {
		n.log.Error("notificationRepo.IsExistNotificationByChannelID: failed to check notification: %v", err)
		return 0, err
	}

	if !isExist {
		notification.ChannelID = channelID
		err := n.notificationRepo.Create(ctx, notification)
		if err != nil {
			n.log.Error("notificationRepo.Create: failed to create notification: %v", err)
			return 0, err
		}
		return 0, nil
	}
	return channelID, nil
}

func (n *notificationService) Delete(ctx context.Context, id int) error {
	return n.notificationRepo.Delete(ctx, id)
}

func (n *notificationService) GetAll(ctx context.Context) ([]entity.Notification, error) {
	return n.notificationRepo.GetAll(ctx)
}

func (n *notificationService) GetByChannelName(ctx context.Context, channelName string) (*entity.Notification, error) {
	channelID, err := n.channelRepo.GetChannelIDByChannelName(ctx, channelName)
	if err != nil {
		n.log.Error("channelRepo.GetChannelIDByChannelName: failed to get channel id: %v", err)
		return nil, err
	}

	notification, err := n.notificationRepo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (n *notificationService) UpdateTextNotification(ctx context.Context, notification *entity.Notification) error {
	channelID, err := n.createNotificationIfNotExist(ctx, notification)
	if err != nil {
		n.log.Error("createNotificationIfNotExist: %v", err)
		return err
	}

	if channelID == 0 && err == nil {
		return nil
	}

	return n.notificationRepo.UpdateTextByChannelID(ctx, notification.NotificationText, channelID)
}

func (n *notificationService) UpdateFileNotification(ctx context.Context, notification *entity.Notification) error {
	channelID, err := n.createNotificationIfNotExist(ctx, notification)
	if err != nil {
		n.log.Error("createNotificationIfNotExist: %v", err)
		return err
	}

	if channelID == 0 && err == nil {
		return nil
	}

	err = n.notificationRepo.UpdateFileByChannelID(ctx, notification.FileID, notification.FileType, channelID)
	if err != nil {
		n.log.Error("notificationRepo.UpdateFileByChannelID: failed to update file in notification: %v", err)
		return err
	}
	return nil
}

func (n *notificationService) UpdateButtonNotification(ctx context.Context, notification *entity.Notification) error {
	channelID, err := n.createNotificationIfNotExist(ctx, notification)
	if err != nil {
		n.log.Error("createNotificationIfNotExist: %v", err)
		return err
	}

	if channelID == 0 && err == nil {
		return nil
	}

	err = n.notificationRepo.UpdateButtonByChannelID(ctx, notification.ButtonURL, notification.ButtonText, channelID)
	if err != nil {
		n.log.Error("notificationRepo.UpdateButtonByChannelID: failed to update button in notification: %v", err)
		return err
	}
	return nil
}

func (n *notificationService) GetByChannelTelegramID(ctx context.Context, channelTelegramID int64) (*entity.Notification, error) {
	return n.notificationRepo.GetByChannelID(ctx, channelTelegramID)
}
