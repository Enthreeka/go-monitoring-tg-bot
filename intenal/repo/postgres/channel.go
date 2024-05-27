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
	GetAllAdminChannel(ctx context.Context) ([]entity.Channel, error)
	GetChannelIDByChannelName(ctx context.Context, channelName string) (int64, error)
	GetByChannelName(ctx context.Context, channelName string) (*entity.Channel, error)
	UpdateNeedCaptchaByChannelName(ctx context.Context, channelName string) error
	GetChannelByUserID(ctx context.Context, userID int64) (string, error)
	GetChannelByChannelTgID(ctx context.Context, channelTgID int64) (*entity.Channel, error)
	SetAcceptTimer(ctx context.Context, channelName string, timer int) error
	UpdateQuestion(ctx context.Context, channelName string, question []byte) error
	GetQuestionByChannelName(ctx context.Context, channelName string) ([]byte, error)
	UpdateQuestionEnabledByChannelName(ctx context.Context, channelName string) error
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
	err := row.Scan(&channel.ID, &channel.TelegramID, &channel.ChannelName,
		&channel.ChannelURL, &channel.Status, &channel.NeedCaptcha, &channel.AcceptTimer, &channel.Question, &channel.QuestionEnabled)
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
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return checkErr
	}

	return err
}

func (u *channelRepo) IsChannelExistByTgID(ctx context.Context, telegramID int64) (bool, error) {
	query := `select exists (select id from channel where tg_id = $1)`
	var isExist bool

	err := u.Pool.QueryRow(ctx, query, telegramID).Scan(&isExist)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return isExist, checkErr
	}

	return isExist, err
}

func (u *channelRepo) GetAllAdminChannel(ctx context.Context) ([]entity.Channel, error) {
	query := `select * from channel where channel_status = 'administrator'`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return u.collectRows(rows)
}

func (u *channelRepo) GetChannelIDByChannelName(ctx context.Context, channelName string) (int64, error) {
	query := `select tg_id from channel where channel_name = $1`
	var ChannelTelegramID int64

	err := u.Pool.QueryRow(ctx, query, channelName).Scan(&ChannelTelegramID)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return 0, checkErr
	}

	return ChannelTelegramID, err
}

func (u *channelRepo) GetByChannelName(ctx context.Context, channelName string) (*entity.Channel, error) {
	query := `select * from channel where channel_name = $1`
	channel := new(entity.Channel)

	err := u.Pool.QueryRow(ctx, query, channelName).Scan(&channel.ID,
		&channel.TelegramID,
		&channel.ChannelName,
		&channel.ChannelURL,
		&channel.Status,
		&channel.NeedCaptcha,
		&channel.AcceptTimer,
		&channel.Question,
		&channel.QuestionEnabled,
	)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}

	return channel, err
}

func (u *channelRepo) UpdateQuestionEnabledByChannelName(ctx context.Context, channelName string) error {
	query := `update channel set question_enabled = not question_enabled where channel_name = $1;`

	_, err := u.Pool.Exec(ctx, query, channelName)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return checkErr
	}

	return err
}

func (u *channelRepo) UpdateNeedCaptchaByChannelName(ctx context.Context, channelName string) error {
	query := `update channel set need_captcha = not need_captcha where channel_name = $1;`

	_, err := u.Pool.Exec(ctx, query, channelName)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return checkErr
	}

	return err
}

func (u *channelRepo) GetChannelByUserID(ctx context.Context, userID int64) (string, error) {
	query := `select channel_name from channel
				join user_channel on  user_channel.channel_tg_id = channel.tg_id
				where user_channel.user_id = $1 `

	var channel string

	err := u.Pool.QueryRow(ctx, query, userID).Scan(&channel)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return channel, checkErr
	}

	return channel, err
}

func (u *channelRepo) GetChannelByChannelTgID(ctx context.Context, channelTgID int64) (*entity.Channel, error) {
	query := `select * from channel where tg_id = $1`
	channel := new(entity.Channel)

	err := u.Pool.QueryRow(ctx, query, channelTgID).Scan(&channel.ID,
		&channel.TelegramID,
		&channel.ChannelName,
		&channel.ChannelURL,
		&channel.Status,
		&channel.NeedCaptcha,
		&channel.AcceptTimer,
		&channel.Question,
		&channel.QuestionEnabled,
	)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}

	return channel, err
}

func (u *channelRepo) SetAcceptTimer(ctx context.Context, channelName string, timer int) error {
	query := `update channel set accept_timer = $1 where channel_name = $2`

	_, err := u.Pool.Exec(ctx, query, timer, channelName)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return checkErr
	}

	return err
}

func (u *channelRepo) UpdateQuestion(ctx context.Context, channelName string, question []byte) error {
	query := `update channel set question = $1 where channel_name = $2`

	_, err := u.Pool.Exec(ctx, query, question, channelName)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return checkErr
	}

	return err
}

func (u *channelRepo) GetQuestionByChannelName(ctx context.Context, channelName string) ([]byte, error) {
	query := `select question from channel where channel_name = $1`
	var question []byte
	err := u.Pool.QueryRow(ctx, query, channelName).Scan(&question)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}
	return question, err
}
