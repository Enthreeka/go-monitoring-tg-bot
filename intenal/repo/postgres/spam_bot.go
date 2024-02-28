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

type SpamBotRepo interface {
	Create(ctx context.Context, bot *entity.SpamBot) error
	GetAll(ctx context.Context) ([]entity.SpamBot, error)
	Delete(ctx context.Context, id int) error
}

type spamBotRepo struct {
	*postgres.Postgres
}

func NewSpamBotRepo(pg *postgres.Postgres) SpamBotRepo {
	return &spamBotRepo{
		pg,
	}
}

func (r *spamBotRepo) collectRow(row pgx.Row) (*entity.SpamBot, error) {
	var bot entity.SpamBot
	err := row.Scan(&bot.ID, &bot.Token, &bot.ChannelName)
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
	return &bot, err
}

func (r *spamBotRepo) collectRows(rows pgx.Rows) ([]entity.SpamBot, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.SpamBot, error) {
		bot, err := r.collectRow(row)
		return *bot, err
	})
}

func (s *spamBotRepo) Create(ctx context.Context, bot *entity.SpamBot) error {
	query := `insert into spam_bot (token,channel_name) values ($1,$2)`

	_, err := s.Pool.Exec(ctx, query, bot.Token, bot.ChannelName)
	return err
}

func (s *spamBotRepo) GetAll(ctx context.Context) ([]entity.SpamBot, error) {
	query := `select * from spam_bot`

	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return s.collectRows(rows)
}

func (s *spamBotRepo) Delete(ctx context.Context, id int) error {
	query := `delete from spam_bot where id = $1`

	_, err := s.Pool.Exec(ctx, query, id)
	return err
}
