package boterror

import (
	"errors"
	"fmt"
)

var (
	ErrIsNotAdmin = NewError("User is not an admin", errors.New("not_admin"))
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

	return "Произошла внутренняя ошибка на сервере"
}
