package middlewares

import (
	tele "gopkg.in/telebot.v3"
)

// CheckMembership - мидлвейр, проверяет подписку на канал
func CheckMembership(channelID, adminUserID int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {

			isMember, err := checkSubscription(channelID, c)
			if err != nil {
				return err
			}

			if !isMember {
				return handleNoAccess(c, adminUserID)
			}

			return next(c)
		}
	}
}

func handleNoAccess(c tele.Context, adminUserID int64) error {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	btnAccessRequest := menu.Text("🛡️ Запросить доступ")
	menu.Reply(menu.Row(btnAccessRequest))

	c.Bot().Handle(&btnAccessRequest, func(c tele.Context) error {

		var senderID = c.Sender().Username
		if c.Sender().Username == "" || c.Sender().Username == " " || c.Sender().Username == "null" {
			senderID = string(rune(c.Sender().ID))
			c.ForwardTo(&tele.Chat{ID: adminUserID})
			return c.Send("У вас скрытый профиль или отсутствует имя пользователя (корокое имя) в настройках Telegram.\n\nСвяжитесь с администратором @frntbck напрямую.")
		}

		_, err := c.Bot().Send(&tele.User{ID: adminUserID}, "Получен запрос на доступ от пользователя @"+senderID)
		if err != nil {
			return c.Send("Ошибка при отправке запроса: " + err.Error())
		}
		c.Send("Запрос отправлен администратору.")
		return c.Delete()
	})

	return c.Send("У вас нет доступа к этому боту. Запросите доступ или свяжитесь с администратором: @frntbck", menu)
}

func checkSubscription(channelID int64, c tele.Context) (bool, error) {
	user := c.Recipient()
	channel := &tele.Chat{ID: channelID}

	chatMember, err := c.Bot().ChatMemberOf(channel, user)
	if err != nil {
		return false, c.Send("Ошибка при проверке подписки: " + err.Error())
	}

	return isMember(chatMember.Role), nil
}

func isMember(role tele.MemberStatus) bool {
	validRoles := map[tele.MemberStatus]bool{
		tele.Member:        true,
		tele.Administrator: true,
		tele.Creator:       true,
	}
	return validRoles[role]
}
