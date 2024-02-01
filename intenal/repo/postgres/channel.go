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

type ChannelRepo interface {
	Create(ctx context.Context, channel *entity.Channel) error
	GetByID(ctx context.Context, id int) (*entity.Channel, error)
	DeleteByID(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.Channel, error)
	UpdateStatusByTgID(ctx context.Context, status string, telegramID int64) error
	IsChannelExistByTgID(ctx context.Context, telegramID int64) (bool, error)
}

type channelRepo struct {
	*postgres.Postgres
}

func NewChannelRepo(pg *postgres.Postgres) ChannelRepo {
	return &channelRepo{
		pg,
	}
}

func (u *channelRepo) collectRow(row pgx.Row) (*entity.Channel, error) {
	var channel entity.Channel
	err := row.Scan(&channel.ID, &channel.TelegramID, &channel.ChannelName, &channel.ChannelURL, &channel.Status)
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
	return &channel, err
}

func (u *channelRepo) collectRows(rows pgx.Rows) ([]entity.Channel, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Channel, error) {
		channel, err := u.collectRow(row)
		return *channel, err
	})
}

func (u *channelRepo) Create(ctx context.Context, channel *entity.Channel) error {
	query := `insert into channel (tg_id,channel_name,channel_url,channel_status) values ($1,$2,$3,$4)`

	_, err := u.Pool.Exec(ctx, query, channel.TelegramID, channel.ChannelName, channel.ChannelURL, channel.Status)
	return err
}

func (u *channelRepo) GetByID(ctx context.Context, id int) (*entity.Channel, error) {
	query := `select * from channel where id = $1`

	row := u.Pool.QueryRow(ctx, query, id)
	return u.collectRow(row)
}

func (u *channelRepo) DeleteByID(ctx context.Context, id int) error {
	query := `delete from channel where id = $1`

	_, err := u.Pool.Exec(ctx, query, id)
	return err
}

func (u *channelRepo) GetAll(ctx context.Context) ([]entity.Channel, error) {
	query := `select * from channel`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return u.collectRows(rows)
}

func (u *channelRepo) UpdateStatusByTgID(ctx context.Context, status string, telegramID int64) error {
	query := `update channel set channel_status = $1 where tg_id = $2`

	_, err := u.Pool.Exec(ctx, query, status, telegramID)
	return err
}

func (u *channelRepo) IsChannelExistByTgID(ctx context.Context, telegramID int64) (bool, error) {
	query := `select exists (select id from channel where tg_id = $1)`
	var isExist bool

	err := u.Pool.QueryRow(ctx, query, telegramID).Scan(&isExist)
	return isExist, err
}
