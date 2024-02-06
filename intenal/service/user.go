package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateRole(ctx context.Context, role string) error
	GetAllIDByChannelID(ctx context.Context, channelName string) ([]int64, error)
	CreateUserChannel(ctx context.Context, userID int64, channelTelegramID int64) error
}

type userService struct {
	userRepo    postgres.UserRepo
	requestRepo postgres.RequestRepo
	channelRepo postgres.ChannelRepo
	log         *logger.Logger
}

func NewUserService(userRepo postgres.UserRepo, requestRepo postgres.RequestRepo, channelRepo postgres.ChannelRepo, log *logger.Logger) UserService {
	return &userService{
		userRepo:    userRepo,
		requestRepo: requestRepo,
		channelRepo: channelRepo,
		log:         log,
	}
}

func (u *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *userService) UpdateRole(ctx context.Context, role string) error {
	return u.userRepo.UpdateRole(ctx, role)
}

func (u *userService) GetAllIDByChannelID(ctx context.Context, channelName string) ([]int64, error) {
	channelID, err := u.channelRepo.GetChannelIDByChannelName(ctx, channelName)
	if err != nil {
		u.log.Error("channelRepo.GetChannelIDByChannelName: failed to get channel id: %v", err)
		return nil, err
	}
	return u.userRepo.GetAllIDByChannelTgID(ctx, channelID)
}

func (u *userService) CreateUser(ctx context.Context, user *entity.User) error {
	isExist, err := u.userRepo.IsUserExistByUsernameTg(ctx, user.UsernameTg)
	if err != nil {
		u.log.Error("userRepo.IsUserExistByUsernameTg: failed to check user: %v", err)
		return err
	}

	if !isExist {
		u.log.Info("Get user: %s, with request: %s", user.String())

		err := u.userRepo.CreateUser(ctx, user)
		if err != nil {
			u.log.Error("userRepo.CreateUser: failed to create user: %v", err)
			return err
		}
	}

	return nil
}

func (u *userService) CreateUserChannel(ctx context.Context, userID int64, channelTelegramID int64) error {
	isExist, err := u.userRepo.IsExistUserChannel(ctx, userID, channelTelegramID)
	if err != nil {
		u.log.Error("userRepo.IsExistUserChannel: failed to check user channel: %v", err)
		return err
	}
	if !isExist {
		err := u.userRepo.CreateUserChannel(ctx, userID, channelTelegramID)
		if err != nil {
			u.log.Error("userRepo.CreateUserChannel: failed to create user channel: %v", err)
			return err
		}
	}
	return nil
}
