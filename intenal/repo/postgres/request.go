package postgres

import (
	"context"
	"errors"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	pgxError "github.com/Entreeka/monitoring-tg-bot/intenal/boterror/pgx_error"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type RequestRepo interface {
	Create(ctx context.Context, request *entity.Request) (*entity.Request, error)
	GetAll(ctx context.Context) ([]entity.Request, error)
	GetAllByStatusRequestAndChannelName(ctx context.Context, status string, channelName string) ([]entity.Request, error)
	UpdateStatusRequestByUserID(ctx context.Context, request *entity.Request) (*entity.Request, error)
	DeleteByStatus(ctx context.Context, status string) error
	DeleteByID(ctx context.Context, id int) error
	GetCountByStatusRequestAndChannelTgID(ctx context.Context, status string, channelTgID int64) (int, error)
	IsExistByUserID(ctx context.Context, userID int64) (bool, error)
	UpdateStatusRequestByID(ctx context.Context, status string, id int) error
	GetCountRequestTodayByChannelID(ctx context.Context, id int64) (int, error)
}

type requestRepo struct {
	*postgres.Postgres
}

func NewRequestRepo(pg *postgres.Postgres) RequestRepo {
	return &requestRepo{
		pg,
	}
}

func (r *requestRepo) collectRow(row pgx.Row) (*entity.Request, error) {
	var req entity.Request
	err := row.Scan(&req.ID, &req.ChannelTelegramID, &req.UserID, &req.StatusRequest, &req.DateRequest)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, boterror.ErrNoRows
	}
	errCode := pgxError.ErrorCode(err)
	if errCode == pgxError.ForeignKeyViolation {
		return nil, boterror.ErrForeignKeyViolation
	}
	if errCode == pgxError.UniqueViolation {
		return nil, boterror.ErrUniqueViolation
	}
	return &req, err
}

func (r *requestRepo) collectRows(rows pgx.Rows) ([]entity.Request, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Request, error) {
		request, err := r.collectRow(row)
		return *request, err
	})
}

func (r *requestRepo) Create(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	query := `insert into request (user_id,status_request,channel_tg_id,date_request) values ($1,$2,$3,$4) returning *`
	req := new(entity.Request)

	err := r.Pool.QueryRow(ctx, query, request.UserID, request.StatusRequest, request.ChannelTelegramID, request.DateRequest).Scan(&req.ID,
		&req.ChannelTelegramID, &req.UserID, &req.StatusRequest, &req.DateRequest)
	return req, err
}

func (r *requestRepo) GetAll(ctx context.Context) ([]entity.Request, error) {
	query := `select * from request`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return r.collectRows(rows)
}

func (r *requestRepo) GetAllByStatusRequestAndChannelName(ctx context.Context, status string, channelName string) ([]entity.Request, error) {
	query := `select request.id, request.channel_tg_id, request.user_id, request.status_request, request.date_request
		from request
		join channel on request.channel_tg_id = channel.tg_id
		where request.status_request = $1 and channel.channel_name = $2`

	rows, err := r.Pool.Query(ctx, query, status, channelName)
	if err != nil {
		return nil, err
	}
	return r.collectRows(rows)
}

func (r *requestRepo) UpdateStatusRequestByUserID(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	query := `update request set status_request = $1, date_request = $2 where user_id = $3 returning *`
	req := new(entity.Request)

	err := r.Pool.QueryRow(ctx, query, request.StatusRequest, request.DateRequest, request.UserID).Scan(&req.ID,
		&req.ChannelTelegramID, &req.UserID, &req.StatusRequest, &req.DateRequest)

	return req, err
}

func (r *requestRepo) DeleteByStatus(ctx context.Context, status string) error {
	query := `delete from request where status_request = $1`

	_, err := r.Pool.Exec(ctx, query, status)
	return err
}

func (r *requestRepo) DeleteByID(ctx context.Context, id int) error {
	query := `delete from request where id = $1`

	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

func (r *requestRepo) GetCountByStatusRequestAndChannelTgID(ctx context.Context, status string, channelTgID int64) (int, error) {
	query := `select count(*) from request where status_request = $1 and channel_tg_id = $2`
	var waitingCount int

	err := r.Pool.QueryRow(ctx, query, status, channelTgID).Scan(&waitingCount)
	return waitingCount, err
}

func (r *requestRepo) IsExistByUserID(ctx context.Context, userID int64) (bool, error) {
	query := `select exists (select id from request where user_id = $1)`
	var isExist bool

	err := r.Pool.QueryRow(ctx, query, userID).Scan(&isExist)
	return isExist, err
}

func (r *requestRepo) UpdateStatusRequestByID(ctx context.Context, status string, id int) error {
	query := `update request set status_request = $1 where id = $2`

	_, err := r.Pool.Exec(ctx, query, status, id)
	return err
}

func (r *requestRepo) GetCountRequestTodayByChannelID(ctx context.Context, id int64) (int, error) {
	query := `select count(*) from request where date_request::date = current_date::date and channel_tg_id = $1`
	var countRequestToday int

	err := r.Pool.QueryRow(ctx, query, id).Scan(&countRequestToday)
	return countRequestToday, err
}
