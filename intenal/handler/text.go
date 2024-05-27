package handler

import (
	"fmt"
	"time"
)

const (
	MessageShowAllChannel = `<strong>Ниже представлен список каналов, в которых бот является администратором</strong>`
)

func MessageGetChannelInfo(channel string, waitingCount int, userCount int, isExistCaptcha bool, questionEnabled bool) string {
	var isExistCaptchaStr string
	if isExistCaptcha == true {
		isExistCaptchaStr = "Да"
	} else {
		isExistCaptchaStr = "Нет"
	}
	var questionEnabledStr string
	if questionEnabled == true {
		questionEnabledStr = "Да"
	} else {
		questionEnabledStr = "Нет"
	}

	return fmt.Sprintf("<strong>Управление каналом</strong>\n"+
		"Канал:<i>%s</i> \n"+
		"Капча включена: %s\n"+
		"Опрос после капчи включен: %s\n\n"+
		"Количество людей, которые ожидают принятия: %d\n\n"+
		"Количество людей в базе бота по данному каналу: %d", channel, isExistCaptchaStr, questionEnabledStr, waitingCount, userCount)
}

const (
	GeneralMainBotMenu     = `<b>Главное меню бота</b>`
	GeneralUserSettingMenu = "<b>Взаимодействие с данными пользователей</b>"
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
	RequestEmpty = `Запросы отсутствуют`
)

func RequestDeclineText(countRejected int) string {
	return fmt.Sprintf("Людей было отклонено: %d", countRejected)

}

func RequestApprovedText(countApproved int) string {
	return fmt.Sprintf("Людей было принято: %d", countApproved)
}

func RequestError(countErr int) string {
	return fmt.Sprintf("Со стороны ограничений телеграмма не удалость обработать %d людей", countErr)
}

func RequestApproveThroughTime(seconds int, countApproved int) string {
	return fmt.Sprintf("Было принято %d людей через заданный промежуток времени: %d", countApproved, seconds)
}

func RequestStatistic(day int, countRequest int, countSentMsg int64, channelName string) string {
	return fmt.Sprintf("За число: %d, было подано заявок: %d, по каналу: %s. Успешно отправленных сообщений"+
		" %d", day, countRequest, channelName, countSentMsg)
}

func NotificationSettingText(channel string) string {
	return fmt.Sprintf("<strong>Управление рассылками для новых пользователей</strong>\n"+
		"Канал:<i>%s</i> \n\n"+
		"Кнопка `<u>Отправить пример рассылки</u>` отправит вам сообщение такого же вида, как это будут видеть новые пользователи", channel)
}

func NotificationGlobalSendingStat(value int64) string {
	return fmt.Sprintf("Рассылка по всей базе завершена, число успешно отправленных сообщений: %d", value)
}

const (
	NotificationUpdateText   = `Отправьте сообщение, которое будет отправляться новым пользователям`
	NotificationUpdateFile   = "Отправьте файл/фотографию, который будет отправляться новым пользователям\n\nЕсли отправляете фотографию, то поставьте галочку для сжатия изображения"
	NotificationUpdateButton = "Отправьте сообщение и ссылку для создания кнопки, которая будет отправляться новым пользователям. \n" +
		"Пример сообщения: на чем написан бот?|https://go.dev/"
	NotificationEmpty        = "Рассылка отсутствует"
	NotificationDeleteText   = "Текст успешно удален"
	NotificationDeleteButton = "Кнопка успешно удалена"
	NotificationDeleteFile   = "Документ/фотография успешно удалена"

	NotificationGlobalSetting      = "<strong>Управление рассылкой для всех пользователей</strong>\n"
	NotificationGlobalUpdateText   = `Отправьте сообщение, которое отправится пользователям`
	NotificationGlobalUpdateFile   = "Отправьте файл/фотографию, которое отправится пользователям\n\nЕсли отправляете фотографию, то поставьте галочку для сжатия изображения"
	NotificationGlobalUpdateButton = "Отправьте сообщение и ссылку для создания кнопки, которое отправится пользователям. \n" +
		"Пример сообщения: на чем написан бот?|https://go.dev/"
)

const (
	SpamBotAdd    = `Отправьте токен бота`
	SpamBotDelete = `Выберите бота, которого хотите удалить из базы`
	SpamBotGet    = `Список всех доступных ботов для рассылок`
)

func BotCaptcha(channel string) string {
	return fmt.Sprintf("Вы действительно хотите присоединиться к каналу: %s?\nДля подтверждения нажмите на /confirm", channel)
}

const (
	ChannelSetTimer = "Отправьте время в минутах, для авто принятия пользователей"
)

func ChannelUpdateQuestion(question string, channel string) string {
	return fmt.Sprintf("Канал:<i>%s</i> \n"+
		"Ниже представлены поля данного опроса:\n\n%s"+
		"\n\nЕсли хотите изменить опрос, отправьте его с измененными полями в формате JSON. Если не хотите менять, нажмите"+
		" кнопку отменить команду", channel, question)
}
