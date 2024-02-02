package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
)

type RequestService interface {
	Create(ctx context.Context, request *entity.Request) error
	GetAll(ctx context.Context) ([]entity.Request, error)
	GetAllByStatusRequest(ctx context.Context, status string) ([]entity.Request, error)
	UpdateStatusRequestByID(ctx context.Context, status string, id int) error
	DeleteByStatus(ctx context.Context, status string) error
	DeleteByID(ctx context.Context, id int) error
	GetCountByStatusRequestAndChannelTgID(ctx context.Context, status string, channelTgID int64) (int, error)
}

type requestService struct {
	requestRepo postgres.RequestRepo
}

func NewRequestService(requestRepo postgres.RequestRepo) RequestService {
	return &requestService{
		requestRepo: requestRepo,
	}
}

func (r *requestService) Create(ctx context.Context, request *entity.Request) error {
	return r.requestRepo.Create(ctx, request)
}

func (r *requestService) GetAll(ctx context.Context) ([]entity.Request, error) {
	return r.requestRepo.GetAll(ctx)
}

func (r *requestService) GetAllByStatusRequest(ctx context.Context, status string) ([]entity.Request, error) {
	return r.requestRepo.GetAllByStatusRequest(ctx, status)
}

func (r *requestService) UpdateStatusRequestByID(ctx context.Context, status string, id int) error {
	return r.requestRepo.UpdateStatusRequestByID(ctx, status, id)
}

func (r *requestService) DeleteByStatus(ctx context.Context, status string) error {
	return r.requestRepo.DeleteByStatus(ctx, status)
}

func (r *requestService) DeleteByID(ctx context.Context, id int) error {
	return r.requestRepo.DeleteByID(ctx, id)
}

func (r *requestService) GetCountByStatusRequestAndChannelTgID(ctx context.Context, status string, channelTgID int64) (int, error) {
	return r.requestRepo.GetCountByStatusRequestAndChannelTgID(ctx, status, channelTgID)
}
