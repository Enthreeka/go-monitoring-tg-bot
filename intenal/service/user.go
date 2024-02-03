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
	GetAllID(ctx context.Context) ([]*int64, error)
	//JoinChannel(ctx context.Context, user *entity.User) error
}

type userService struct {
	userRepo    postgres.UserRepo
	requestRepo postgres.RequestRepo
	log         *logger.Logger
}

func NewUserService(userRepo postgres.UserRepo, requestRepo postgres.RequestRepo, log *logger.Logger) UserService {
	return &userService{
		userRepo:    userRepo,
		requestRepo: requestRepo,
		log:         log,
	}
}

//func (u *userService) CreateUser(ctx context.Context, user *entity.User) error {
//	isExist, err := u.userRepo.IsUserExistByUsernameTg(ctx, user.UsernameTg)
//	if err != nil {
//		u.log.Error("userRepo.IsUserExistByUsernameTg: failed to check user: %v", err)
//		return err
//	}
//
//	if !isExist {
//		u.log.Info("Get user: %s, with request: %s", user.String())
//
//		err := u.userRepo.CreateUser(ctx, user)
//		if err != nil {
//			u.log.Error("userRepo.CreateUser: failed to create user: %v", err)
//			return err
//		}
//	}
//	return nil
//}

func (u *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *userService) UpdateRole(ctx context.Context, role string) error {
	return u.userRepo.UpdateRole(ctx, role)
}

func (u *userService) GetAllID(ctx context.Context) ([]*int64, error) {
	return u.userRepo.GetAllID(ctx)
}

func (u *userService) CreateUser(ctx context.Context, user *entity.User) error {
	u.log.Info("Get user: %s, with request: %s", user.String())

	isExist, err := u.userRepo.IsUserExistByUsernameTg(ctx, user.UsernameTg)
	if err != nil {
		u.log.Error("userRepo.IsUserExistByUsernameTg: failed to check user: %v", err)
		return err
	}

	if !isExist {
		err := u.userRepo.CreateUser(ctx, user)
		if err != nil {
			u.log.Error("userRepo.CreateUser: failed to create user: %v", err)
			return err
		}
	}

	return nil
}
