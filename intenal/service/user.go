package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateRole(ctx context.Context, role string) error
	GetAllID(ctx context.Context) ([]*int64, error)
}

type userService struct {
	userRepo postgres.UserRepo
}

func NewUserService(userRepo postgres.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUser(ctx context.Context, user *entity.User) error {
	return u.userRepo.CreateUser(ctx, user)
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

func (u *userService) GetAllID(ctx context.Context) ([]*int64, error) {
	return u.userRepo.GetAllID(ctx)
}
