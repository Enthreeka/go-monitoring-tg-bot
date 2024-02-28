package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
)

type SpamBotService interface {
	Create(ctx context.Context, bot *entity.SpamBot) error
	GetAll(ctx context.Context) ([]entity.SpamBot, error)
	Delete(ctx context.Context, id int) error
}

type spamBotService struct {
	userRepo    postgres.UserRepo
	spamBotRepo postgres.SpamBotRepo
}

func NewSpamBotService(userRepo postgres.UserRepo, spamBotRepo postgres.SpamBotRepo) SpamBotService {
	return &spamBotService{
		userRepo:    userRepo,
		spamBotRepo: spamBotRepo,
	}
}

func (s *spamBotService) Create(ctx context.Context, bot *entity.SpamBot) error {
	return s.spamBotRepo.Create(ctx, bot)
}

func (s *spamBotService) GetAll(ctx context.Context) ([]entity.SpamBot, error) {
	return s.spamBotRepo.GetAll(ctx)
}

func (s *spamBotService) Delete(ctx context.Context, id int) error {
	return s.spamBotRepo.Delete(ctx, id)
}
