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

type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateRole(ctx context.Context, role string) error
	GetAllID(ctx context.Context) ([]*int64, error)
	IsUserExistByUsernameTg(ctx context.Context, usernameTg string) (bool, error)
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
	return &user, err
}

func (u *userRepo) collectRows(rows pgx.Rows) ([]entity.User, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.User, error) {
		user, err := u.collectRow(row)
		return *user, err
	})
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

func (u *userRepo) UpdateRole(ctx context.Context, role string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) GetAllID(ctx context.Context) ([]*int64, error) {
	query := `select id from "user"`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allID := make([]*int64, 0, 256)

	for rows.Next() {
		var id int64

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		allID = append(allID, &id)
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
	return isExist, err
}
