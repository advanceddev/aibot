package middlewares

import (
	"unrealbot/config"

	tele "gopkg.in/telebot.v3"
)

// CheckMembership - мидлвейр, проверяет подписку на канал
func CheckMembership(next tele.HandlerFunc) tele.HandlerFunc {
	cfg := config.MustLoad()

	return func(c tele.Context) error {
		if err := checkSubscription(cfg, c); err != nil {
			return err
		}

		return next(c)
	}
}

func checkSubscription(cfg *config.Config, c tele.Context) error {
	user := c.Recipient()
	channel := &tele.Chat{ID: cfg.ChannelID}
	menu := &tele.ReplyMarkup{RemoveKeyboard: true}

	chatMember, err := c.Bot().ChatMemberOf(channel, user)
	if err != nil {
		return c.Send("Ошибка при проверке подписки: " + err.Error())
	}

	if isMember(chatMember) {
		return nil
	}

	return c.Send("У вас нет доступа к этому боту. Свяжитесь с администратором: @frntbck", menu)
}

func isMember(chatMember *tele.ChatMember) bool {
	return chatMember.Role == tele.Member || chatMember.Role == tele.Administrator || chatMember.Role == tele.Creator
}
