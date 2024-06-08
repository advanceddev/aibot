package middlewares

import (
	tele "gopkg.in/telebot.v3"
)

// CheckMembership - мидлвейр, проверяет подписку на канал
func CheckMembership(channelID int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			isMember, err := checkSubscription(channelID, c)
			if err != nil {
				return err
			}
			if !isMember {
				menu := &tele.ReplyMarkup{RemoveKeyboard: true}
				return c.Send("У вас нет доступа к этому боту. Свяжитесь с администратором: @frntbck", menu)
			}
			return next(c)
		}
	}
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
