package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChannelService interface {
	Create(ctx context.Context, channel *entity.Channel) error
	GetByID(ctx context.Context, id int) (*entity.Channel, error)
	DeleteByID(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.Channel, error)
	ChatMember(ctx context.Context, channel *entity.Channel) error
	GetAllAdminChannel(ctx context.Context) (*tgbotapi.InlineKeyboardMarkup, error)
	GetByChannelName(ctx context.Context, channelName string) (*entity.Channel, error)
	UpdateNeedCaptchaByChannelName(ctx context.Context, channelName string) error
	GetChannelByUserID(ctx context.Context, userID int64) (string, error)
	GetChannelByChannelTgID(ctx context.Context, channelTgID int64) (*entity.Channel, error)
	SetAcceptTimer(ctx context.Context, channelName string, timer int) error
	GetQuestion(ctx context.Context, channelName string) (string, *tgbotapi.InlineKeyboardMarkup, error)
	UpdateQuestion(ctx context.Context, channelName string, question []byte) error
	GetQuestionByChannelName(ctx context.Context, channelName string) ([]byte, error)
	UpdateQuestionEnabledByChannelName(ctx context.Context, channelName string) error
	GetChannelAfterConfirm(ctx context.Context, userID int64) (*entity.Channel, error)
}

type channelService struct {
	channelRepo postgres.ChannelRepo
	log         *logger.Logger
	store       *stateful.Store
}

const (
	kicked        = "kicked"
	administrator = "administrator"
	left          = "left"
)

func NewChannelService(channelRepo postgres.ChannelRepo, log *logger.Logger, store *stateful.Store) ChannelService {
	return &channelService{
		channelRepo: channelRepo,
		log:         log,
		store:       store,
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

func (c *channelService) GetAllAdminChannel(ctx context.Context) (*tgbotapi.InlineKeyboardMarkup, error) {
	channel, err := c.channelRepo.GetAllAdminChannel(ctx)
	if err != nil {
		return nil, err
	}

	return c.createChannelMarkup(channel, "get")
}

func (c *channelService) createChannelMarkup(channel []entity.Channel, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range channel {
		if el.TelegramID != 0 { // check for channel for global notification
			btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.ChannelName),
				fmt.Sprintf("channel_%s_%d", command, el.ID))

			row = append(row, btn)

			if (i+1)%buttonsPerRow == 0 || i == len(channel)-1 {
				rows = append(rows, row)
				row = []tgbotapi.InlineKeyboardButton{}
			}
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.MainMenuButton})
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}

func (c *channelService) GetByChannelName(ctx context.Context, channelName string) (*entity.Channel, error) {
	return c.channelRepo.GetByChannelName(ctx, channelName)
}

func (c *channelService) UpdateNeedCaptchaByChannelName(ctx context.Context, channelName string) error {
	return c.channelRepo.UpdateNeedCaptchaByChannelName(ctx, channelName)
}

func (c *channelService) GetChannelByUserID(ctx context.Context, userID int64) (string, error) {
	return c.channelRepo.GetChannelByUserID(ctx, userID)
}

func (c *channelService) GetChannelByChannelTgID(ctx context.Context, channelTgID int64) (*entity.Channel, error) {
	return c.channelRepo.GetChannelByChannelTgID(ctx, channelTgID)
}

func (c *channelService) SetAcceptTimer(ctx context.Context, channelName string, timer int) error {
	return c.channelRepo.SetAcceptTimer(ctx, channelName, timer)
}

func (c *channelService) GetQuestion(ctx context.Context, channelName string) (string, *tgbotapi.InlineKeyboardMarkup, error) {
	channel, err := c.channelRepo.GetByChannelName(ctx, channelName)
	if err != nil {
		c.log.Error("channelRepo.GetByChannelName: failed to query channel: %v", err)
		return "", nil, err
	}

	baseChannel := base64.StdEncoding.EncodeToString([]byte(channelName))

	questionModel := new(entity.QuestionModel)
	if err := json.Unmarshal(channel.Question, questionModel); err != nil {
		c.log.Error("channelRepo.GetByChannelName: failed to unmarshal question: %v", err)
		return "", nil, err
	}
	questionModel.ChanelNameBase64 = baseChannel

	mrk, err := c.createQuestionMarkup(questionModel)
	if err != nil {
		c.log.Error("channelRepo.GetByChannelName: failed to create question markup: %v", err)
		return "", nil, err
	}

	return questionModel.Question, mrk, nil
}

func (c *channelService) createQuestionMarkup(questionModel *entity.QuestionModel) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range questionModel.Answer {
		btn := tgbotapi.NewInlineKeyboardButtonData(el.AnswerVariation,
			fmt.Sprintf("answer_%s_%d", questionModel.ChanelNameBase64, el.ID))

		row = append(row, btn)

		if (i+1)%buttonsPerRow == 0 || i == len(questionModel.Answer)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}

func (c *channelService) UpdateQuestion(ctx context.Context, channelName string, question []byte) error {
	return c.channelRepo.UpdateQuestion(ctx, channelName, question)
}

func (c *channelService) GetQuestionByChannelName(ctx context.Context, channelName string) ([]byte, error) {
	return c.channelRepo.GetQuestionByChannelName(ctx, channelName)
}

func (c *channelService) UpdateQuestionEnabledByChannelName(ctx context.Context, channelName string) error {
	return c.channelRepo.UpdateQuestionEnabledByChannelName(ctx, channelName)
}

func (c *channelService) GetChannelAfterConfirm(ctx context.Context, userID int64) (*entity.Channel, error) {
	return c.channelRepo.GetChannelAfterConfirm(ctx, userID)
}
