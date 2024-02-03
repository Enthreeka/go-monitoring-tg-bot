package callback

import (
	"fmt"
	"time"
)

const (
	messageShowAllChannel = `<strong>Ниже представлен список каналов, в которых бот является администратором</strong>`
)

func messageGetChannelInfo(channel string, waitingCount int) string {
	return fmt.Sprintf("<strong>Управление каналом</strong>\n\n"+
		"Канал:<i>%s</i> \n"+
		"Количество людей, которые ожидают принятия: %d", channel, waitingCount)
}

const (
	generalMainBotMenu     = `<b>Главное меню бота</b>`
	generalUserSettingMenu = "<b>Взаимодествие с данными пользователей</b>"
)

func userExcelFileText() string {
	return fmt.Sprintf("<i>Выгрузка данных на:</i> %v", time.Now().Format("15:04 2006-01-02"))
}

const (
	requestApproved = `Все заявки статуса "in progress" были приняты`
	requestDecline  = `Все заявки статуса "in progress" были отклонены`
)
