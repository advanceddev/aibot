package middlewares

import (
	"unrealbot/config"

	tele "gopkg.in/telebot.v3"
)

// CheckMembership - мидлвейр, проверяет подписку на канал
func CheckMembership(next tele.HandlerFunc) tele.HandlerFunc {
	cfg := config.MustLoad()

	return func(c tele.Context) error {
		isMember, err := checkSubscription(cfg, c)
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

func checkSubscription(cfg *config.Config, c tele.Context) (bool, error) {
	user := c.Recipient()
	channel := &tele.Chat{ID: cfg.ChannelID}

	chatMember, err := c.Bot().ChatMemberOf(channel, user)
	if err != nil {
		return false, c.Send("Ошибка при проверке подписки: " + err.Error())
	}

	return isMember(chatMember), nil
}

func isMember(chatMember *tele.ChatMember) bool {
	return chatMember.Role == tele.Member || chatMember.Role == tele.Administrator || chatMember.Role == tele.Creator
}
