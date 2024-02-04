package callback

import (
	"fmt"
	"time"
)

const (
	messageShowAllChannel = `<strong>Ниже представлен список каналов, в которых бот является администратором</strong>`
)

func messageGetChannelInfo(channel string, waitingCount int) string {
	return fmt.Sprintf("<strong>Управление каналом</strong>\n"+
		"Канал:<i>%s</i> \n\n"+
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

func requestApproveThroughTime(seconds int) string {
	return fmt.Sprintf("Все заявки статуса \"in progress\" были приняты через заданный промежуток времени: : %d", seconds)
}

func notificationSettingText(channel string) string {
	return fmt.Sprintf("<strong>Управление рассылок для новых пользователей</strong>\n"+
		"Канал:<i>%s</i> \n\n"+
		"Последняя кнопка отправит вам сообщение такого же вида, как это будут видеть новые пользователи", channel)
}

const (
	notificationUpdateText   = `Отправьте сообщение, которое будет отправляться новым пользователям`
	notificationUpdateFile   = `Отправьте файл/фотографию, который будет отправляться новым пользователям`
	notificationUpdateButton = "Отправьте сообщение и ссылку для создания кнопки, которая будет отправляться новым пользователям. \n" +
		"Пример сообщения: на чем написан бот?|https://go.dev/"
	notificationCancel = "Команда была отменена"
	notificationEmpty  = "Рассылка отсутствует"
)
