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

	case strings.HasPrefix(callbackData, "approved_all"):
		callbackView, ok := b.callbackView["approved_all"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "rejected_all"):
		callbackView, ok := b.callbackView["rejected_all"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "approved_time"):
		callbackView, ok := b.callbackView["approved_time"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "hello_setting"):
		callbackView, ok := b.callbackView["hello_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "add_text_notification"):
		callbackView, ok := b.callbackView["add_text_notification"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "add_photo_notification"):
		callbackView, ok := b.callbackView["add_photo_notification"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "add_button_notification"):
		callbackView, ok := b.callbackView["add_button_notification"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "example_notification"):
		callbackView, ok := b.callbackView["example_notification"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "cancel_setting"):
		callbackView, ok := b.callbackView["cancel_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "delete_text_notification"):
		callbackView, ok := b.callbackView["delete_text_notification"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "delete_photo_notification"):
		callbackView, ok := b.callbackView["delete_photo_notification"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "delete_button_notification"):
		callbackView, ok := b.callbackView["delete_button_notification"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "sender_setting"):
		callbackView, ok := b.callbackView["sender_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "send_message"):
		callbackView, ok := b.callbackView["send_message"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "update_sender_message"):
		callbackView, ok := b.callbackView["update_sender_message"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "delete_sender_message"):
		callbackView, ok := b.callbackView["delete_sender_message"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "example_sender_message"):
		callbackView, ok := b.callbackView["example_sender_message"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "comeback"):
		callbackView, ok := b.callbackView["comeback"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "cancel_sender_setting"):
		callbackView, ok := b.callbackView["cancel_sender_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "role_setting"):
		callbackView, ok := b.callbackView["role_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "create_admin"):
		callbackView, ok := b.callbackView["create_admin"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "create_super_admin"):
		callbackView, ok := b.callbackView["create_super_admin"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "delete_admin"):
		callbackView, ok := b.callbackView["delete_admin"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "all_admin"):
		callbackView, ok := b.callbackView["all_admin"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "cancel_admin_setting"):
		callbackView, ok := b.callbackView["cancel_admin_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "bot_spam_settings"):
		callbackView, ok := b.callbackView["bot_spam_settings"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	default:
		return nil, nil
	}
}
