package middlewares

import (
	"sync"
	"unrealbot/cmd/bot"
	"unrealbot/internal/handlers/chat"
	"unrealbot/internal/utils"

	tele "gopkg.in/telebot.v3"
)

var chatPool = sync.Pool{
	New: func() interface{} {
		return &tele.Chat{}
	},
}

// CheckMembership - мидлвейр, проверяет подписку на канал
func CheckMembership(bot bot.UnrealBot) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {

			isMember, err := checkSubscription(bot.ChannelID, c)
			if err != nil {
				return err
			}

			if !isMember {
				return handleNoAccess(c, bot)
			}

			return next(c)
		}
	}
}

func handleNoAccess(c tele.Context, bot bot.UnrealBot) error {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	btnAccessRequest := menu.Text("🛡️ Запросить доступ")
	menu.Reply(menu.Row(btnAccessRequest))
	cmd := chat.NewCommandHandler(&bot)
	c.Bot().Handle(&btnAccessRequest, cmd.RequestSubscribeHandler)

	// Избегаем лишнего выделения памяти для строки
	msg := "У вас нет доступа к этому боту. Запросите доступ или свяжитесь с администратором."
	return c.Send(msg, menu)
}

// Используем указатели для передачи объектов
func checkSubscription(channelID int64, c tele.Context) (bool, error) {
	user := c.Recipient()

	channel := chatPool.Get().(*tele.Chat)
	channel.ID = channelID

	chatMember, err := c.Bot().ChatMemberOf(channel, user)
	if err != nil {
		errMsg := utils.SumStrings("Ошибка при проверке подписки: ", err.Error())
		return false, c.Send(errMsg)
	}

	// Очищаем объект и возвращаем его в пул
	chatPool.Put(channel)

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
