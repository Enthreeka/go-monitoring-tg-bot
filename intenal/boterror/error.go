package boterror

import (
	"errors"
	"fmt"
)

var (
	ErrIsNotAdmin = NewError("User is not an admin", errors.New("not_admin"))
)

var (
	ErrUniqueViolation     = NewError("Violation must be unique", errors.New("non_unique_value"))
	ErrForeignKeyViolation = NewError("Foreign Key Violation", errors.New("foreign_key_violation "))
	ErrNoRows              = NewError("No rows in result set", errors.New("no_rows"))
	ErrNotFound            = NewError("Tasks not found", errors.New("not_found"))
)

var (
	ErrNil = NewError("Nil pointer value", errors.New("nil_pointer"))
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

	}

	return "Произошла внутрення ошибка на сервере"
}
