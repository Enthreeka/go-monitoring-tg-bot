package tgbot

import (
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func channelUpdateToModel(update *tgbotapi.Update) *entity.Channel {
	channel := &entity.Channel{
		TelegramID:  update.MyChatMember.Chat.ID,
		ChannelName: update.MyChatMember.Chat.Title,
		Status:      update.MyChatMember.NewChatMember.Status,
	}

	if update.MyChatMember.Chat.UserName != "" {
		url := "t.me/" + update.MyChatMember.Chat.UserName
		channel.ChannelURL = &url
	}

	return channel
}

func userUpdateToModel(update *tgbotapi.Update) *entity.User {
	user := new(entity.User)

	if update.Message != nil {
		user.ID = update.Message.From.ID
		user.UsernameTg = update.Message.From.UserName
		user.CreatedAt = time.Now().Local()
		user.Role = roleUser
	}

	if update.ChatJoinRequest != nil {
		user.ID = update.ChatJoinRequest.From.ID
		user.UsernameTg = update.ChatJoinRequest.From.UserName
		user.ChannelFrom = &update.ChatJoinRequest.InviteLink.InviteLink
		user.CreatedAt = time.Now().Local()
		user.Role = roleUser
	}
	return user
}

func requestUpdateToModel(update *tgbotapi.Update) *entity.Request {
	return &entity.Request{
		UserID:            update.ChatJoinRequest.From.ID,
		ChannelTelegramID: update.ChatJoinRequest.Chat.ID,
		StatusRequest:     RequestInProgress,
		DateRequest:       time.Now().Local(),
	}
}
