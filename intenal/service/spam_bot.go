package service

import (
	"context"
	"fmt"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SpamBotService interface {
	Create(ctx context.Context, bot *entity.SpamBot) error
	GetAllBots(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error)
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

func (s *spamBotService) GetAllBots(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	bots, err := s.spamBotRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return s.createBotMarkup(bots, command)
}

func (s *spamBotService) Delete(ctx context.Context, id int) error {
	return s.spamBotRepo.Delete(ctx, id)
}

func (s *spamBotService) createBotMarkup(bot []entity.SpamBot, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range bot {
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.ChannelName),
			fmt.Sprintf("bot_%s_%d", command, el.ID))

		row = append(row, btn)

		if (i+1)%buttonsPerRow == 0 || i == len(bot)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.MainMenuButton})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}
