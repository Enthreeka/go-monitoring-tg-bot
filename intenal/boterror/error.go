package boterror

import (
	"errors"
	"fmt"
)

var (
	ErrIsNotAdmin       = NewError("User is not an admin", errors.New("not_admin"))
	ErrIsNotSuperAdmin  = NewError("User is not super admin", errors.New("not_super_admin"))
	ErrNotFoundUser     = NewError("Not found user", errors.New("not_found"))
	ErrDeleteSuperAdmin = NewError("Delete super admin in tg bot", errors.New("delete_super_admin"))
)

var (
	ErrUniqueViolation     = NewError("Violation must be unique", errors.New("non_unique_value"))
	ErrForeignKeyViolation = NewError("Foreign Key Violation", errors.New("foreign_key_violation "))
	ErrNoRows              = NewError("No rows in result set", errors.New("no_rows"))
)

var (
	ErrNil        = NewError("Nil pointer value", errors.New("nil_pointer"))
	ErrNotFoundID = NewError("ID in callback not found", errors.New("empty_id"))
)

type BotError struct {
	Msg string `json:"message"`
	Err error  `json:"-"`
}

func (a *BotError) Error() string {
	return fmt.Sprintf("%s", a.Msg)
}

func NewError(msg string, err error) *BotError {
	return &BotError{
		Msg: msg,
		Err: err,
	}
}

func ParseErrToText(err error) string {
	switch {
	case errors.Is(err, ErrIsNotAdmin):
		return "Нет администраторских прав доступа"
	case errors.Is(err, ErrIsNotSuperAdmin):
		return "Нет прав доступа супер администратора"
	case errors.Is(err, ErrNotFoundUser):
		return "Пользователь с таким никнеймом не был найден"
	case errors.Is(err, ErrDeleteSuperAdmin):
		return "Нельзя забирать права супер админа через бота, необходимо изменять через базу данных"

	}

	return "Произошла внутрення ошибка на сервере"
}
