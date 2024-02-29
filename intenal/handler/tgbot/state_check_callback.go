package tgbot

import (
	"context"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) getStateCallback(ctx context.Context, update *tgbotapi.Update) (bool, error) {
	storeData, isExist := b.isStateExist(update.SentFrom().ID)
	if isExist {

		typeData := getStoreData(storeData)
		if typeData == nil {
			b.log.Error("failed to get data: typeData == nil")
			return true, ErrEmptyStoreData
		}

		switch typeData.(type) {
		case *stateful.SpamBot:
			if err := b.storeDataSpamBotOperationType(ctx, storeData, update); err != nil {
				b.log.Error("storeDataSpamBotOperationType: %v", err)
				return true, err
			}
			return true, nil

		}
		return true, nil
	}
	return false, nil
}
