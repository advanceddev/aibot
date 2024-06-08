package chat

import (
	"unrealbot/cmd/bot"

	tele "gopkg.in/telebot.v3"
)

// NewCommandHandler создает новый обработчик сообщений
func NewCommandHandler(bot *bot.UnrealBot) *Handler {
	return &Handler{bot: bot}
}

// StartHandler обрабатывает команду /start
func (h *Handler) StartHandler(ctx tele.Context) error {
	menu := &tele.ReplyMarkup{RemoveKeyboard: true}
	return ctx.Send("Привет, "+ctx.Sender().FirstName+"! 👋", menu)
}

// ContactHandler обрабатывает команду /contact
func (h *Handler) ContactHandler(ctx tele.Context) error {
	ctx.Notify("typing")
	phone := ctx.Message().Contact.PhoneNumber
	return ctx.Send("Записал твой номер: " + phone + "!")
}
