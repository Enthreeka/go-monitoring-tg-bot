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

type SenderRepo interface {
	Create(ctx context.Context, sender *entity.Sender) error
	Update(ctx context.Context, sender *entity.Sender) error
	DeleteByID(ctx context.Context, channelID int64) error
	GetByChannelName(ctx context.Context, channelName string) (*entity.Sender, error)
	IsExistByChannelName(ctx context.Context, channelName string) (bool, error)
}

type senderRepo struct {
	*postgres.Postgres
}

func NewSenderRepo(pg *postgres.Postgres) SenderRepo {
	return &senderRepo{
		pg,
	}
}

func (s *senderRepo) collectRow(row pgx.Row) (*entity.Sender, error) {
	var sender entity.Sender
	err := row.Scan(&sender.ID, &sender.ChannelTelegramID, &sender.Message)
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
	return &sender, err
}

func (s *senderRepo) Create(ctx context.Context, sender *entity.Sender) error {
	query := `insert into sender (channel_tg_id,message) values ($1,$2)`

	_, err := s.Pool.Exec(ctx, query, sender.ChannelTelegramID, sender.Message)
	return err
}

func (s *senderRepo) Update(ctx context.Context, sender *entity.Sender) error {
	query := `update sender set message = $1 where channel_tg_id = $2`

	_, err := s.Pool.Exec(ctx, query, sender.Message, sender.ChannelTelegramID)
	return err
}

func (s *senderRepo) DeleteByID(ctx context.Context, channelID int64) error {
	query := `delete from sender where channel_tg_id = $1`

	_, err := s.Pool.Exec(ctx, query, channelID)
	return err
}

func (s *senderRepo) GetByChannelName(ctx context.Context, channelName string) (*entity.Sender, error) {
	query := `select sender.id, sender.channel_tg_id,sender.message from sender
         join channel on sender.channel_tg_id = channel.tg_id
         where channel.channel_name = $1`
	sender := new(entity.Sender)

	err := s.Pool.QueryRow(ctx, query, channelName).Scan(&sender.ID, &sender.ChannelTelegramID, &sender.Message)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}

	return sender, err
}

func (s *senderRepo) IsExistByChannelName(ctx context.Context, channelName string) (bool, error) {
	query := `select exists (select sender.id from sender
    	join channel on sender.channel_tg_id = channel.tg_id
         where channel.channel_name = $1)`
	var isExist bool

	err := s.Pool.QueryRow(ctx, query, channelName).Scan(&isExist)
	return isExist, err
}
