package callback

import (
	"github.com/Entreeka/monitoring-tg-bot/intenal/handler/tgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CallbackStrings(update *tgbotapi.Update, bot *tgbot.Bot) (error, tgbot.ViewFunc) {
	//callbackData := update.CallbackData()

	//switch {
	//case strings.HasPrefix():
	//
	//default:
	//	return nil, nil
	//}
	return nil, nil
}
