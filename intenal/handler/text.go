package handler

import (
	"fmt"
	"time"
)

const (
	MessageShowAllChannel = `<strong>Ниже представлен список каналов, в которых бот является администратором</strong>`
)

func MessageGetChannelInfo(channel string, waitingCount int) string {
	return fmt.Sprintf("<strong>Управление каналом</strong>\n"+
		"Канал:<i>%s</i> \n\n"+
		"Количество людей, которые ожидают принятия: %d", channel, waitingCount)
}

const (
	GeneralMainBotMenu     = `<b>Главное меню бота</b>`
	GeneralUserSettingMenu = "<b>Взаимодествие с данными пользователей</b>"
)

const (
	UserUpdateSenderText  = "Отправьте сообщение для рассылки базе"
	UserDeleteSenderText  = "Сообщение успешно удалено"
	UserSenderEmpty       = "Сообщение отсутствует"
	UserSenderError       = "Внутренняя ошибка при рассылке пользователям"
	UserSenderErrorEmpty  = "Пользователей по данному каналу в базе не было найдено"
	UserSenderDone        = "Рассылка завершена"
	UserSuperAdminSetting = "<strong>Управление администраторами</strong>\n\nДанная команда доступна только людям с правами " +
		"<i>Супер администратор</i>\n\n" +
		"<i>Супер администратор</i> - права позволяют управлять правами любого участника бота, а также самим ботом\n" +
		"<i>Администратор</i> - права позволяют управлять ботом"
	UserSetAdmin      = "Отправьте никнейм пользователя, которого хотите назначить администратором"
	UserSetSuperAdmin = "Отправьте никнейм пользователя, которого хотите назначить супер администратором"
	UserDeleteAdmin   = "Отправьте никнейм пользователя, у которого хотите забрать админские права"
	UserAdminCancel   = "Команда отменена"
)

func UserExcelFileText() string {
	return fmt.Sprintf("<i>Выгрузка данных на:</i> %v", time.Now().Format("15:04 2006-01-02"))
}

func UserSenderSetting(channel string) string {
	return fmt.Sprintf("<strong>Управление рассылок по базе с пользователями</strong>\n"+
		"Канал:<i>%s</i> \n\n"+
		"Рассылка сообщения производится по пользователям выбранного канала", channel)
}

const (
	RequestApproved = `Все заявки статуса "in progress" были приняты`
	RequestDecline  = `Все заявки статуса "in progress" были отклонены`
	RequestEmpty    = `Запросы отсутствуют`
)

func RequestApproveThroughTime(seconds int) string {
	return fmt.Sprintf("Все заявки статуса \"in progress\" были приняты через заданный промежуток времени: : %d", seconds)
}

func NotificationSettingText(channel string) string {
	return fmt.Sprintf("<strong>Управление рассылок для новых пользователей</strong>\n"+
		"Канал:<i>%s</i> \n\n"+
		"Кнопка `<u>Отправить пример рассылки</u>` отправит вам сообщение такого же вида, как это будут видеть новые пользователи", channel)
}

const (
	NotificationUpdateText   = `Отправьте сообщение, которое будет отправляться новым пользователям`
	NotificationUpdateFile   = "Отправьте файл/фотографию, который будет отправляться новым пользователям\n\nЕсли отправляете фотографию, то поставьте галочку для сжатия изображения"
	NotificationUpdateButton = "Отправьте сообщение и ссылку для создания кнопки, которая будет отправляться новым пользователям. \n" +
		"Пример сообщения: на чем написан бот?|https://go.dev/"
	NotificationCancel       = "Команда была отменена"
	NotificationEmpty        = "Рассылка отсутствует"
	NotificationDeleteText   = "Текст успешно удален"
	NotificationDeleteButton = "Кнопка успешно удалена"
	NotificationDeleteFile   = "Документ/фотография успешно удалена"
	NotificationExampleError = "С кнопкой обязательно должно быть сообщение/файл"
)
