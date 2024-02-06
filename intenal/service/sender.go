package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
)

type SenderService interface {
	CreateSender(ctx context.Context, sender *entity.Sender) error
	DeleteSender(ctx context.Context, channelName string) error
	GetSender(ctx context.Context, channelName string) (*entity.Sender, error)
}

type senderService struct {
	senderRepo  postgres.SenderRepo
	channelRepo postgres.ChannelRepo
}

func NewSenderService(senderRepo postgres.SenderRepo, channelRepo postgres.ChannelRepo) SenderService {
	return &senderService{
		senderRepo:  senderRepo,
		channelRepo: channelRepo,
	}
}

func (s *senderService) CreateSender(ctx context.Context, sender *entity.Sender) error {
	isExist, err := s.senderRepo.IsExistByChannelName(ctx, sender.ChannelName)
	if err != nil {
		return err
	}

	sender.ChannelTelegramID, err = s.channelRepo.GetChannelIDByChannelName(ctx, sender.ChannelName)
	if err != nil {
		return err
	}

	if !isExist {
		err := s.senderRepo.Create(ctx, sender)
		if err != nil {
			return err
		}
		return nil
	}

	err = s.senderRepo.Update(ctx, sender)
	if err != nil {
		return err
	}
	return nil
}

func (s *senderService) DeleteSender(ctx context.Context, channelName string) error {
	isExist, err := s.senderRepo.IsExistByChannelName(ctx, channelName)
	if err != nil {
		return err
	}

	if isExist {
		sender, err := s.senderRepo.GetByChannelName(ctx, channelName)
		if err != nil {
			return err
		}

		err = s.senderRepo.DeleteByID(ctx, sender.ChannelTelegramID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *senderService) GetSender(ctx context.Context, channelName string) (*entity.Sender, error) {
	return s.senderRepo.GetByChannelName(ctx, channelName)
}
