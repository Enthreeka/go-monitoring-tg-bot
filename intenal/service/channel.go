package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
)

type ChannelService interface {
	Create(ctx context.Context, channel *entity.Channel) error
	GetByID(ctx context.Context, id int) (*entity.Channel, error)
	DeleteByID(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.Channel, error)
	ChatMember(ctx context.Context, channel *entity.Channel) error
}

type channelService struct {
	channelRepo postgres.ChannelRepo
	log         *logger.Logger
}

const (
	kicked        = "kicked"
	administrator = "administrator"
	left          = "left"
)

func NewChannelService(channelRepo postgres.ChannelRepo, log *logger.Logger) ChannelService {
	return &channelService{
		channelRepo: channelRepo,
		log:         log,
	}
}

func (c *channelService) Create(ctx context.Context, channel *entity.Channel) error {
	return c.channelRepo.Create(ctx, channel)
}

func (c *channelService) GetByID(ctx context.Context, id int) (*entity.Channel, error) {
	return c.channelRepo.GetByID(ctx, id)
}

func (c *channelService) DeleteByID(ctx context.Context, id int) error {
	return c.channelRepo.DeleteByID(ctx, id)
}

func (c *channelService) GetAll(ctx context.Context) ([]entity.Channel, error) {
	return c.channelRepo.GetAll(ctx)
}

func (c *channelService) ChatMember(ctx context.Context, channel *entity.Channel) error {
	c.log.Info("Get channel: %s", channel.String())

	isExist, err := c.channelRepo.IsChannelExistByTgID(ctx, channel.TelegramID)
	if err != nil {
		c.log.Error("channelRepo.IsChannelExistByTgID: failed to check channel: %v", err)
		return err
	}

	if !isExist {
		err := c.channelRepo.Create(ctx, channel)
		if err != nil {
			c.log.Error("channelRepo.Create: failed to create channel: %v", err)
			return err
		}
		return nil
	}

	err = c.channelRepo.UpdateStatusByTgID(ctx, channel.Status, channel.TelegramID)
	if err != nil {
		c.log.Error("channelRepo.UpdateStatusByTgID: failed to update channel status: %v", err)
		return err
	}
	return nil
}
