package chat

import (
	"unrealbot/cmd/bot"
	"unrealbot/internal/utils"

	tele "gopkg.in/telebot.v4"
)

// NewCommandHandler создает новый обработчик сообщений
func NewCommandHandler(bot *bot.UnrealBot) *Handler {
	return &Handler{bot: bot}
}

// StartHandler обрабатывает команду /start
func (h *Handler) StartHandler(ctx tele.Context) error {
	menu := &tele.ReplyMarkup{RemoveKeyboard: true}
	return ctx.Send(utils.SumStrings("Привет, ", ctx.Sender().FirstName, "! 👋"), menu)
}

// RequestSubscribeHandler - обрабатывает запрос на подписку
func (h *Handler) RequestSubscribeHandler(ctx tele.Context) error {
	var senderID = ctx.Sender().Username
	if ctx.Sender().Username == "" || ctx.Sender().Username == " " || ctx.Sender().Username == "null" {
		senderID = string(rune(ctx.Sender().ID))
		ctx.ForwardTo(&tele.Chat{ID: h.bot.AdminUserID})
		return ctx.Send("У вас скрытый профиль или отсутствует имя пользователя (корокое имя) в настройках Telegram.\n\nСвяжитесь с администратором напрямую.")
	}

	_, err := ctx.Bot().Send(&tele.User{ID: h.bot.AdminUserID}, utils.SumStrings("Получен запрос на доступ от пользователя @", senderID))
	if err != nil {
		return ctx.Send(utils.SumStrings("Ошибка при отправке запроса: ", err.Error()))
	}
	ctx.Send("Запрос отправлен администратору.")
	return ctx.Delete()
}
