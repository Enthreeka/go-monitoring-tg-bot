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

type NotificationRepo interface {
	Create(ctx context.Context, notification *entity.Notification) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.Notification, error)
	GetByChannelID(ctx context.Context, channelID int64) (*entity.Notification, error)
	UpdateTextByChannelID(ctx context.Context, text string, channelID int64) error
	UpdateFileByChannelID(ctx context.Context, fileID string, fileType string, channelID int64) error
	UpdateButtonByChannelID(ctx context.Context, button string, channelID int64) error
	IsExistNotificationByChannelID(ctx context.Context, channelID int64) (bool, error)
}

type notificationRepo struct {
	*postgres.Postgres
}

func NewNotificationRepo(pg *postgres.Postgres) NotificationRepo {
	return &notificationRepo{
		pg,
	}
}

func (n *notificationRepo) collectRow(row pgx.Row) (*entity.Notification, error) {
	var notify entity.Notification
	err := row.Scan(&notify.ID, &notify.ChannelID, &notify.NotificationText, &notify.FileID, &notify.FileType, &notify.ButtonURL)
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
	return &notify, err
}

func (n *notificationRepo) collectRows(rows pgx.Rows) ([]entity.Notification, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Notification, error) {
		notification, err := n.collectRow(row)
		return *notification, err
	})
}

func (n *notificationRepo) Create(ctx context.Context, notification *entity.Notification) error {
	query := `insert into notification (channel_id,notification_text,file_id,file_type,button_url) values ($1,$2,$3,$4,$5)`

	_, err := n.Pool.Exec(ctx, query, notification.ChannelID,
		notification.NotificationText,
		notification.FileID,
		notification.FileType,
		notification.ButtonURL,
	)
	return err
}

func (n *notificationRepo) Delete(ctx context.Context, id int) error {
	query := `delete from notification where id = $1`

	_, err := n.Pool.Exec(ctx, query, id)
	return err
}

func (n *notificationRepo) GetAll(ctx context.Context) ([]entity.Notification, error) {
	query := `select * from notification`

	rows, err := n.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return n.collectRows(rows)
}

func (n *notificationRepo) GetByChannelID(ctx context.Context, channelID int64) (*entity.Notification, error) {
	query := `select * from notification where channel_id = $1`

	row := n.Pool.QueryRow(ctx, query, channelID)
	return n.collectRow(row)
}

func (n *notificationRepo) UpdateTextByChannelID(ctx context.Context, text string, channelID int64) error {
	query := `update notification set notification_text = $1 where channel_id = $2`

	_, err := n.Pool.Exec(ctx, query, text, channelID)
	return err
}

func (n *notificationRepo) UpdateFileByChannelID(ctx context.Context, fileID string, fileType string, channelID int64) error {
	query := `update notification set file_id = $1, file_type = $2 where channel_id = $3`

	_, err := n.Pool.Exec(ctx, query, fileID, fileType, channelID)
	return err
}

func (n *notificationRepo) UpdateButtonByChannelID(ctx context.Context, button string, channelID int64) error {
	query := `update notification set button_url = $1 where channel_id = $2`

	_, err := n.Pool.Exec(ctx, query, button, channelID)
	return err
}

func (n *notificationRepo) IsExistNotificationByChannelID(ctx context.Context, channelID int64) (bool, error) {
	query := `select exists (select id from notification where channel_id = $1)`
	var isExist bool

	err := n.Pool.QueryRow(ctx, query, channelID).Scan(&isExist)
	return isExist, err
}
