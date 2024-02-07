package postgres

import (
	"context"
	pgxError "github.com/Entreeka/monitoring-tg-bot/intenal/boterror/pgx_error"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateRoleByUsername(ctx context.Context, role string, username string) error
	GetAllIDByChannelTgID(ctx context.Context, channelTelegramID int64) ([]int64, error)
	IsUserExistByUsernameTg(ctx context.Context, usernameTg string) (bool, error)
	GetAllAdmin(ctx context.Context) ([]entity.User, error)
	UserChannelRepo
}

type UserChannelRepo interface {
	CreateUserChannel(ctx context.Context, userID int64, channelTelegramID int64) error
	IsExistUserChannel(ctx context.Context, userID int64, channelTelegramID int64) (bool, error)
}

type userRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) UserRepo {
	return &userRepo{
		pg,
	}
}

func (u *userRepo) collectRow(row pgx.Row) (*entity.User, error) {
	var user entity.User
	err := row.Scan(&user.ID, &user.UsernameTg, &user.CreatedAt, &user.Phone, &user.ChannelFrom, &user.Role)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}

	return &user, err
}

func (u *userRepo) collectRows(rows pgx.Rows) ([]entity.User, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.User, error) {
		user, err := u.collectRow(row)
		return *user, err
	})
}

func (u *userRepo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `select * from "user" where tg_username = $1`

	row := u.Pool.QueryRow(ctx, query, username)
	return u.collectRow(row)
}

func (u *userRepo) CreateUser(ctx context.Context, user *entity.User) error {
	query := `insert into "user" (id,tg_username,created_at,phone,channel_from,user_role) values ($1,$2,$3,$4,$5,$6)`

	_, err := u.Pool.Exec(ctx, query, user.ID, user.UsernameTg, user.CreatedAt, user.Phone, user.ChannelFrom, user.Role)
	return err
}

func (u *userRepo) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	query := `select * from "user"`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return u.collectRows(rows)
}

func (u *userRepo) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `select * from "user" where id = $1`

	row := u.Pool.QueryRow(ctx, query, id)
	return u.collectRow(row)
}

func (u *userRepo) UpdateRoleByUsername(ctx context.Context, role string, username string) error {
	query := `update "user" set user_role = $1 where tg_username = $2`

	_, err := u.Pool.Exec(ctx, query, role, username)
	return err
}

func (u *userRepo) GetAllIDByChannelTgID(ctx context.Context, channelTelegramID int64) ([]int64, error) {
	query := `select id from "user"
			join user_channel on user_channel.user_id = "user".id
			where user_channel.channel_tg_id = $1`

	rows, err := u.Pool.Query(ctx, query, channelTelegramID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allID := make([]int64, 0, 256)

	for rows.Next() {
		var id int64

		err := rows.Scan(&id)
		if err != nil {
			if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
				return nil, checkErr
			}
			return nil, err
		}

		allID = append(allID, id)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return allID, nil
}

func (u *userRepo) IsUserExistByUsernameTg(ctx context.Context, usernameTg string) (bool, error) {
	query := `select exists (select id from "user" where tg_username = $1)`
	var isExist bool

	err := u.Pool.QueryRow(ctx, query, usernameTg).Scan(&isExist)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return isExist, checkErr
	}

	return isExist, err
}

func (u *userRepo) CreateUserChannel(ctx context.Context, userID int64, channelTelegramID int64) error {
	query := `insert into user_channel (user_id, channel_tg_id) values ($1,$2)`

	_, err := u.Pool.Exec(ctx, query, userID, channelTelegramID)
	return err
}

func (u *userRepo) IsExistUserChannel(ctx context.Context, userID int64, channelTelegramID int64) (bool, error) {
	query := `select exists (select user_id from user_channel where user_id = $1 and channel_tg_id = $2)`
	var isExist bool

	err := u.Pool.QueryRow(ctx, query, userID, channelTelegramID).Scan(&isExist)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return isExist, checkErr
	}

	return isExist, err
}

func (u *userRepo) GetAllAdmin(ctx context.Context) ([]entity.User, error) {
	query := `select * from "user" where user_role = 'admin' or user_role = 'superAdmin'`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return u.collectRows(rows)
}
