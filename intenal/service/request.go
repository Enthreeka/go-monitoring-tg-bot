package service

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/repo/postgres"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
)

type RequestService interface {
	CreateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error)
	GetAll(ctx context.Context) ([]entity.Request, error)
	GetAllByStatusRequest(ctx context.Context, status string, channelName string) ([]entity.Request, error)
	DeleteByStatus(ctx context.Context, status string) error
	DeleteByID(ctx context.Context, id int) error
	GetCountByStatusRequestAndChannelTgID(ctx context.Context, status string, channelTgID int64) (int, error)
	UpdateStatusRequestByID(ctx context.Context, status string, id int) error
	GetCountRequestTodayByChannelID(ctx context.Context, id int64) (int, error)
}

type requestService struct {
	requestRepo postgres.RequestRepo
	log         *logger.Logger
}

func NewRequestService(requestRepo postgres.RequestRepo, log *logger.Logger) RequestService {
	return &requestService{
		requestRepo: requestRepo,
		log:         log,
	}
}

func (r *requestService) CreateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	r.log.Info("Get request: %s", request.String())
	var (
		req *entity.Request
	)

	isExist, err := r.requestRepo.IsExistByUserIDAndChannelID(ctx, request.UserID, request.ChannelTelegramID)
	if err != nil {
		r.log.Error("requestRepo.IsExistByUserID: failed to check user in requests: %v", err)
		return nil, err
	}

	if isExist {
		req, err = r.requestRepo.UpdateStatusRequestByUserID(ctx, request)
		if err != nil {
			r.log.Error("requestRepo.UpdateStatusRequestByID: failed to update request")
			return nil, err
		}
		return req, nil
	}

	// create only `in progress`
	req, err = r.requestRepo.Create(ctx, request)
	if err != nil {
		r.log.Error("requestRepo.Create: failed to create request")
		return nil, err
	}

	return req, nil
}

func (r *requestService) GetAll(ctx context.Context) ([]entity.Request, error) {
	return r.requestRepo.GetAll(ctx)
}

func (r *requestService) GetAllByStatusRequest(ctx context.Context, status string, channelName string) ([]entity.Request, error) {
	return r.requestRepo.GetAllByStatusRequestAndChannelName(ctx, status, channelName)
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

func (r *requestService) UpdateStatusRequestByID(ctx context.Context, status string, id int) error {
	return r.requestRepo.UpdateStatusRequestByID(ctx, status, id)
}

func (r *requestService) GetCountRequestTodayByChannelID(ctx context.Context, id int64) (int, error) {
	return r.requestRepo.GetCountRequestTodayByChannelID(ctx, id)
}
