package tgbot

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Entreeka/monitoring-tg-bot/intenal/entity"
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler"
	"github.com/Entreeka/monitoring-tg-bot/pkg/stateful"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/button"
	"github.com/Entreeka/monitoring-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"net/url"
	"strconv"
)

var (
	ErrOperationType   = errors.New("operation type not found")
	ErrEmptyStoreData  = errors.New("all values is nil")
	ErrEmptyButtonData = errors.New("button is incorrect")
	ErrEmptyFile       = errors.New("file is empty")
	ErrUrl             = errors.New("link not valid")
)

const (
	success = "Успешно выполнено"
)

func ParseJSON[T any](src string) (T, error) {
	var args T
	if err := json.Unmarshal([]byte(src), &args); err != nil {
		return *(new(T)), err
	}
	return args, nil
}

func (b *Bot) isStateExist(userID int64) (*stateful.StoreData, bool) {
	data, exist := b.store.Read(userID)
	return data, exist
}

func (b *Bot) getState(ctx context.Context, update *tgbotapi.Update) (bool, error) {
	storeData, isExist := b.isStateExist(update.Message.From.ID)
	if isExist {

		typeData := getStoreData(storeData)
		if typeData == nil {
			b.log.Error("failed to get data: typeData == nil")
			return true, ErrEmptyStoreData
		}

		switch typeData.(type) {
		case *stateful.Channel:
			if err := b.storeDataChannelOperationType(ctx, storeData, update.Message); err != nil {
				b.log.Error("storeDataNotificationOperationType: %v", err)
				return true, err
			}
			return true, nil

		case *stateful.Notification:
			if err := b.storeDataNotificationOperationType(ctx, storeData, update.Message); err != nil {
				b.log.Error("storeDataNotificationOperationType: %v", err)
				return true, err
			}
			return true, nil

		case *stateful.Sender:
			if err := b.createSenderMessage(ctx, storeData, update.Message); err != nil {
				b.log.Error("createSenderMessage: %v", err)
				return true, err
			}
			return true, nil

		case *stateful.Admin:
			if err := b.storeDataUserOperationType(ctx, storeData, update.Message); err != nil {
				b.log.Error("storeDataUserOperationType: %v", err)
				return true, err
			}
			return true, nil

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

func (b *Bot) storeDataSpamBotOperationType(ctx context.Context, storeData *stateful.StoreData, update *tgbotapi.Update) error {
	switch storeData.SpamBot.OperationType {
	case stateful.OperationDeleteBot:
		if update.CallbackQuery != nil {
			userID := update.CallbackQuery.From.ID

			botID := entity.GetID(update.CallbackData())

			bot, err := b.spamBotService.GetByID(ctx, botID)
			if err != nil {
				b.log.Error("spamBotService.GetByID: %v", err)
				return err
			}

			defer func() {
				b.store.Delete(userID)
				b.spammerStorage.Delete(bot.BotName)
			}()

			err = b.spamBotService.Delete(ctx, botID)
			if err != nil {
				b.log.Error("spamBotService.Delete: %v", err)
				return err
			}

			if _, err := b.bot.Send(tgbotapi.NewMessage(userID, "Удален успешно")); err != nil {
				b.log.Error("failed to send message", zap.Error(err))
				return err
			}
		}
		return nil
	case stateful.OperationAddBot:
		userID := update.Message.From.ID

		defer b.store.Delete(userID)

		botName, err := b.spammerStorage.InitializeBot(update.Message.Text)
		if err != nil {
			b.log.Error("spammerStorage.InitializeBot: %v", err)
			return err
		}

		if err = b.spamBotService.Create(ctx, &entity.SpamBot{
			Token:   update.Message.Text,
			BotName: botName,
		}); err != nil {
			b.log.Error("spamBotService.Create: %v", err)
			return err
		}

		if _, err := b.bot.Send(tgbotapi.NewMessage(userID, "Добавлен успешно")); err != nil {
			b.log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}

	return ErrOperationType
}

func (b *Bot) storeDataChannelOperationType(ctx context.Context, storeData *stateful.StoreData, update *tgbotapi.Message) error {
	userID := update.From.ID
	defer b.store.Delete(userID)
	b.log.Info("storeDataChannelOperationType: %v", storeData)

	switch storeData.Channel.OperationType {
	case stateful.OperationSetTimer:
		timer, err := strconv.Atoi(update.Text)
		if err != nil {
			return errors.New("send not int type")
		}
		if err := b.channelService.SetAcceptTimer(ctx, storeData.Channel.ChannelName, timer); err != nil {
			b.log.Error("channelService.SetAcceptTimer: %v", err)
			return err
		}

		if err := b.requestTimer(userID, storeData, update); err != nil {
			b.log.Error("requestTimer: %v", err)
			return err
		}

		return nil

	case stateful.OperationUpdateQuestion:
		args, err := ParseJSON[entity.QuestionModel](update.Text)
		if err != nil {
			b.log.Error("ParseJSON: %v", err)
			return err
		}

		index := 1
		for i, _ := range args.Answer {
			args.Answer[i].ID = index
			index++
		}

		questionModelByte, err := json.Marshal(args)
		if err != nil {
			b.log.Error("json.Marshal: %v", err)
			return err
		}

		if err := b.channelService.UpdateQuestion(ctx, storeData.Channel.ChannelName, questionModelByte); err != nil {
			b.log.Error("channelService.UpdateQuestion: %v", err)
			return err
		}

		if err := b.requestChannel(userID, storeData, update, &args); err != nil {
			b.log.Error("requestChannel: %v", err)
			return err
		}

		return nil
	}

	return ErrOperationType
}

func (b *Bot) storeDataNotificationOperationType(ctx context.Context, storeData *stateful.StoreData, update *tgbotapi.Message) error {
	userID := update.From.ID

	switch storeData.Notification.OperationType {
	case stateful.OperationUpdateText:
		defer b.store.Delete(userID)

		notification := &entity.Notification{
			NotificationText: &update.Text,
			FileID:           nil,
			FileType:         nil,
			ButtonURL:        nil,
			ButtonText:       nil,
			ChannelName:      storeData.Notification.ChannelName,
		}
		if err := b.notificationService.UpdateTextNotification(ctx, notification); err != nil {
			b.log.Error("notificationService.UpdateTextNotification: %v", err)
			return err
		}

		if err := b.requestNotification(userID, storeData, update); err != nil {
			return err
		}

		return nil
	case stateful.OperationUpdateFile:
		defer b.store.Delete(userID)

		notification := &entity.Notification{
			NotificationText: nil,
			ButtonURL:        nil,
			ButtonText:       nil,
			ChannelName:      storeData.Notification.ChannelName,
		}
		ft := fileType(update)
		if ft == "" {
			return ErrEmptyFile
		}
		notification.FileType = &ft

		if ft == "photo" {
			largestPhoto := update.Photo[len(update.Photo)-1]
			fileID := largestPhoto.FileID
			notification.FileID = &fileID
		} else {
			fileID := update.Document.FileID
			notification.FileID = &fileID
		}

		if err := b.notificationService.UpdateFileNotification(ctx, notification); err != nil {
			b.log.Error("notificationService.UpdateFileNotification: %v", err)
			return err
		}

		if err := b.requestNotification(userID, storeData, update); err != nil {
			return err
		}

		return nil
	case stateful.OperationUpdateButton:
		defer b.store.Delete(userID)

		btnUrl, btnText := entity.GetButtonData(update.Text)
		if btnUrl == "" || btnText == "" {
			return ErrEmptyButtonData
		}
		if !isUrl(btnUrl) {
			return ErrUrl
		}

		notification := &entity.Notification{
			NotificationText: nil,
			FileID:           nil,
			FileType:         nil,
			ButtonURL:        &btnUrl,
			ButtonText:       &btnText,
			ChannelName:      storeData.Notification.ChannelName,
		}

		if err := b.notificationService.UpdateButtonNotification(ctx, notification); err != nil {
			b.log.Error("notificationService.UpdateButtonNotification: failed to work with notification: %v", err)
			return err
		}

		if err := b.requestNotification(userID, storeData, update); err != nil {
			return err
		}

		return nil
	}

	return ErrOperationType
}

func (b *Bot) createSenderMessage(ctx context.Context, storeData *stateful.StoreData, update *tgbotapi.Message) error {
	userID := update.From.ID
	defer b.store.Delete(userID)

	sender := &entity.Sender{
		Message:     update.Text,
		ChannelName: storeData.Sender.ChannelName,
	}

	if err := b.senderService.CreateSender(ctx, sender); err != nil {
		b.log.Error("senderService.CreateSender: failed to create/update sender message: %v", err)
		return err
	}

	_, err := b.bot.Send(tgbotapi.NewMessage(userID, success))
	if err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}

	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		storeData.Sender.MessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", storeData.Sender.MessageID, string(resp.Result), err)
	}

	if err := b.sendSenderSetting(userID, storeData.Sender.ChannelName); err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}

	return nil
}

func getStoreData(storeData *stateful.StoreData) any {
	switch {
	case storeData.Notification != nil:
		return storeData.Notification
	case storeData.Sender != nil:
		return storeData.Sender
	case storeData.Admin != nil:
		return storeData.Admin
	case storeData.SpamBot != nil:
		return storeData.SpamBot
	case storeData.Channel != nil:
		return storeData.Channel

	}
	return nil
}

func fileType(update *tgbotapi.Message) string {
	switch {
	case update.Document != nil:
		return update.Document.MimeType
	case update.Photo != nil:
		return "photo"
	}
	return ""
}

func isUrl(str string) bool {
	parsedUrl, _ := url.Parse(str)
	return parsedUrl.Scheme == "" || parsedUrl.Host == ""
}

func (b *Bot) sendHelloSetting(userID int64, channelName string) error {
	var (
		msg tgbotapi.MessageConfig
	)

	if channelName != "" {
		msg = tgbotapi.NewMessage(userID, handler.NotificationSettingText(channelName))
		msg.ReplyMarkup = &markup.HelloMessageSetting
	} else {
		msg = tgbotapi.NewMessage(userID, handler.NotificationGlobalSetting)
		msg.ReplyMarkup = &markup.GlobalHelloMessageSetting
	}
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err := b.bot.Send(msg); err != nil {
		b.log.Error("failed to send message", zap.Error(err))
		return err
	}
	return nil
}

func (b *Bot) sendSenderSetting(userID int64, channelName string) error {
	msg := tgbotapi.NewMessage(userID, handler.UserSenderSetting(channelName))
	msg.ReplyMarkup = &markup.SenderMessageSetting
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err := b.bot.Send(msg); err != nil {
		b.log.Error("failed to send message", zap.Error(err))
		return err
	}
	return nil
}

func (b *Bot) storeDataUserOperationType(ctx context.Context, storeData *stateful.StoreData, update *tgbotapi.Message) error {
	userID := update.From.ID

	switch storeData.Admin.OperationType {
	case stateful.OperationAddAdmin:
		defer b.store.Delete(userID)

		if err := b.userService.UpdateRoleByUsername(ctx, stateful.OperationAddAdmin, update.Text); err != nil {
			b.log.Error("userService.UpdateRoleByUsername: %v", err)
			return err
		}

		if err := b.requestAdmin(userID, storeData, update); err != nil {
			return err
		}

		return nil
	case stateful.OperationAddSuperAdmin:
		defer b.store.Delete(userID)

		if err := b.userService.UpdateRoleByUsername(ctx, stateful.OperationAddSuperAdmin, update.Text); err != nil {
			b.log.Error("userService.UpdateRoleByUsername: %v", err)
			return err
		}

		if err := b.requestAdmin(userID, storeData, update); err != nil {
			return err
		}

		return nil
	case stateful.OperationDeleteAdmin:
		defer b.store.Delete(userID)

		if err := b.userService.UpdateRoleByUsername(ctx, stateful.OperationDeleteAdmin, update.Text); err != nil {
			b.log.Error("userService.UpdateRoleByUsername: %v", err)
			return err
		}

		if err := b.requestAdmin(userID, storeData, update); err != nil {
			return err
		}

		return nil
	}

	return ErrOperationType
}

func (b *Bot) sendAdminSetting(userID int64) error {
	msg := tgbotapi.NewMessage(userID, handler.UserSuperAdminSetting)
	msg.ReplyMarkup = &markup.SuperAdminSetting
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err := b.bot.Send(msg); err != nil {
		b.log.Error("failed to send message", zap.Error(err))
		return err
	}
	return nil
}

func (b *Bot) sendQuestionSetting(userID int64, question *entity.QuestionModel, channelName string) error {
	newQ, err := json.MarshalIndent(question, "", "  ")
	if err != nil {
		b.log.Error("channelRepo.GetByChannelName: failed to marshal exampleModel: %v", err)
		return err
	}

	b.log.Info("question: %s", string(newQ))
	msg := tgbotapi.NewMessage(userID, handler.ChannelUpdateQuestion(string(newQ), channelName))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button.ComebackSetting))
	msg.ParseMode = tgbotapi.ModeHTML

	_, err = b.bot.Send(msg)
	if err != nil {
		b.log.Error("failed to send message", zap.Error(err))
		return err
	}
	return nil
}

func (b *Bot) requestNotification(userID int64, storeData *stateful.StoreData, update *tgbotapi.Message) error {
	_, err := b.bot.Send(tgbotapi.NewMessage(userID, success))
	if err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}

	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		storeData.Notification.MessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", storeData.Notification.MessageID, string(resp.Result), err)
	}

	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		update.MessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", storeData.Notification.MessageID, string(resp.Result), err)
	}

	if err := b.sendHelloSetting(userID, storeData.Notification.ChannelName); err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}
	return nil
}

func (b *Bot) requestAdmin(userID int64, storeData *stateful.StoreData, update *tgbotapi.Message) error {
	_, err := b.bot.Send(tgbotapi.NewMessage(userID, success))
	if err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}

	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		storeData.Admin.MessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", storeData.Admin.MessageID, string(resp.Result), err)
	}

	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		update.MessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", storeData.Admin.MessageID, string(resp.Result), err)
	}

	if err := b.sendAdminSetting(userID); err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}
	return nil
}

func (b *Bot) requestChannel(userID int64, storeData *stateful.StoreData, update *tgbotapi.Message, question *entity.QuestionModel) error {
	_, err := b.bot.Send(tgbotapi.NewMessage(userID, success))
	if err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}

	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		storeData.Channel.MessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", storeData.Admin.MessageID, string(resp.Result), err)
	}

	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		update.MessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", storeData.Admin.MessageID, string(resp.Result), err)
	}

	if err := b.sendQuestionSetting(userID, question, storeData.Channel.ChannelName); err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}
	return nil
}

func (b *Bot) requestTimer(userID int64, storeData *stateful.StoreData, update *tgbotapi.Message) error {
	_, err := b.bot.Send(tgbotapi.NewMessage(userID, success))
	if err != nil {
		b.log.Error("failed to send msg: %v", err)
		return err
	}

	return nil
}
