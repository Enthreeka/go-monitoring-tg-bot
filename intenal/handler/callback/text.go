package callback

import "fmt"

const (
	messageShowAllChannel = `<strong>Ниже представлен список каналов, в которых бот является администратором</strong>`
)

func messageGetChannelInfo(channel string, waitingCount int) string {
	return fmt.Sprintf("<strong>Управление каналом</strong>\n\n"+
		"Канал: <i>%s</i> \n"+
		"Количество людей, которые ожидают принятия: %d", channel, waitingCount)
}

const (
	generalMainBotMenu     = `<b>Главное меню бота</b>`
	generalUserSettingMenu = "<b>Взаимодествие с данными пользователей</b>"
)
