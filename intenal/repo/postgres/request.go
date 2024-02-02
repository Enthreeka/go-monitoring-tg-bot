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
	Create(ctx context.Context, request *entity.Request) error
	GetAll(ctx context.Context) ([]entity.Request, error)
	GetAllByStatusRequest(ctx context.Context, status string) ([]entity.Request, error)
	UpdateStatusRequestByID(ctx context.Context, status string, id int) error
	DeleteByStatus(ctx context.Context, status string) error
	DeleteByID(ctx context.Context, id int) error
	GetCountByStatusRequestAndChannelTgID(ctx context.Context, status string, channelTgID int64) (int, error)
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
	err := row.Scan(&req.ID, &req.UserID, req.ChannelTelegramID, &req.StatusRequest, &req.DateRequest)
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

func (r *requestRepo) Create(ctx context.Context, request *entity.Request) error {
	query := `insert into request (user_id,status_request,channel_tg_id,date_request) values ($1,$2,$3,$4)`

	_, err := r.Pool.Exec(ctx, query, request.UserID, request.StatusRequest, request.ChannelTelegramID, request.DateRequest)
	return err
}

func (r *requestRepo) GetAll(ctx context.Context) ([]entity.Request, error) {
	query := `select * from request`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return r.collectRows(rows)
}

func (r *requestRepo) GetAllByStatusRequest(ctx context.Context, status string) ([]entity.Request, error) {
	query := `select id, channel_tg_id, user_id, status_request, date_request
		from request
		where status_request = $1`

	rows, err := r.Pool.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	return r.collectRows(rows)
}

func (r *requestRepo) UpdateStatusRequestByID(ctx context.Context, status string, id int) error {
	query := `update request set status_request = $1 where id = $2`

	_, err := r.Pool.Exec(ctx, query, status, id)
	return err
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
