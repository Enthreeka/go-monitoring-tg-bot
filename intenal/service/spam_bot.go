package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/button"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/spam"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SpamBotService interface {
	Create(ctx context.Context, bot *entity.SpamBot) error
	GetAllBots(ctx context.Context, command string) (*tgbotapi.InlineKeyboardMarkup, error)
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*entity.SpamBot, error)
	GetSpamBotsFromDBToCache(ctx context.Context)
}

type spamBotService struct {
	userRepo       postgres.UserRepo
	spamBotRepo    postgres.SpamBotRepo
	spammerStorage spam.SpamBot
	log            *logger.Logger
}

func NewSpamBotService(userRepo postgres.UserRepo, spamBotRepo postgres.SpamBotRepo, spammerStorage spam.SpamBot, log *logger.Logger) SpamBotService {
	return &spamBotService{
		userRepo:       userRepo,
		spamBotRepo:    spamBotRepo,
		spammerStorage: spammerStorage,
		log:            log,
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

func (s *spamBotService) GetByID(ctx context.Context, id int) (*entity.SpamBot, error) {
	return s.spamBotRepo.GetByID(ctx, id)
}

func (s *spamBotService) createBotMarkup(bot []entity.SpamBot, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range bot {
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.BotName),
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

func (s *spamBotService) GetSpamBotsFromDBToCache(ctx context.Context) {
	bots, err := s.spamBotRepo.GetAll(ctx)
	if err != nil {
		if errors.Is(err, boterror.ErrNoRows) {
			return
		}
		s.log.Fatal("failed to get bots token from postgres")
		return
	}

	if len(bots) != 0 {
		for _, bot := range bots {
			_, err := s.spammerStorage.InitializeBot(bot.Token)
			if err != nil {
				s.log.Error("failed to initialize bot with start service: %v", err)
			}
		}
	}
}
