package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
)

type NotificationService interface {
	Create(ctx context.Context, notification *entity.Notification) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.Notification, error)
	GetByChannelID(ctx context.Context, channelID int64) (*entity.Notification, error)
}

type notificationService struct {
	notificationRepo postgres.NotificationRepo
}

func NewNotificationService(notificationRepo postgres.NotificationRepo) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
	}
}

func (n *notificationService) Create(ctx context.Context, notification *entity.Notification) error {
	return n.notificationRepo.Create(ctx, notification)
}

func (n *notificationService) Delete(ctx context.Context, id int) error {
	return n.notificationRepo.Delete(ctx, id)
}

func (n *notificationService) GetAll(ctx context.Context) ([]entity.Notification, error) {
	return n.notificationRepo.GetAll(ctx)
}

func (n *notificationService) GetByChannelID(ctx context.Context, channelID int64) (*entity.Notification, error) {
	return n.notificationRepo.GetByChannelID(ctx, channelID)
}
