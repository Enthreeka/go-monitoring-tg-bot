package tgbot

import (
	"errors"
	"strings"
)

var (
	ErrNotFound = errors.New("not found in map")
)

func (b *Bot) CallbackStrings(callbackData string) (error, ViewFunc) {
	switch {

	case strings.HasPrefix(callbackData, "channel_get_"):
		callbackView, ok := b.callbackView["channel_get"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "channel_setting"):
		callbackView, ok := b.callbackView["channel_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "main_menu"):
		callbackView, ok := b.callbackView["main_menu"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "user_setting"):
		callbackView, ok := b.callbackView["user_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "download_excel"):
		callbackView, ok := b.callbackView["download_excel"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	default:
		return nil, nil
	}
}
