package pgx_error

import (
	"errors"
	"github.com/Entreeka/monitoring-tg-bot/intenal/boterror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func ErrorHandler(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return boterror.ErrNoRows
	}
	errCode := ErrorCode(err)
	if errCode == ForeignKeyViolation {
		return boterror.ErrForeignKeyViolation
	}
	if errCode == UniqueViolation {
		return boterror.ErrUniqueViolation
	}

	return nil
}
