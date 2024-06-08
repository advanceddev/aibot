package chat

import (
	"encoding/json"
	"net/http"
	"time"
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

// SubscribeHandler генерирует уникальную пригласительную ссылку на группу
func (h *Handler) SubscribeHandler(ctx tele.Context) error {
	channel := &tele.Chat{ID: h.bot.ChannelID, Type: "privatechannel"}
	link, err := ctx.Bot().InviteLink(channel)
	if err != nil {
		return ctx.Send("Произошла ошибка при формировании пригласительной ссылки.")
	}
	return ctx.Send(link)
}

// BalanceHandler обрабатывает команду /balance
func (h *Handler) BalanceHandler(ctx tele.Context) error {
	url := h.bot.APIUrl + "/user"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ctx.Send("Произошла ошибка при создании запроса: ", err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+h.bot.APIToken)

	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return ctx.Send("Произошла ошибка при запросе данных: ", err.Error())
	}
	defer res.Body.Close()

	var message GenAIUserResponse
	if err := json.NewDecoder(res.Body).Decode(&message); err != nil {
		return ctx.Send("Произошла ошибка при декодировании JSON: ", err.Error())
	}

	return ctx.Send("💰 Баланс GenAPI: " + message.Balance + "₽\n\nПополнить: https://gen-api.ru/account/billing")
}
