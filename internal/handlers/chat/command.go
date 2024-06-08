package chat

import (
	"encoding/json"
	"net/http"
	"time"
	"unrealbot/cmd/bot"
	"unrealbot/internal/utils"

	tele "gopkg.in/telebot.v3"
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

// ContactHandler обрабатывает команду /contact
func (h *Handler) ContactHandler(ctx tele.Context) error {
	ctx.Notify("typing")
	phone := ctx.Message().Contact.PhoneNumber
	return ctx.Send(utils.SumStrings("Записал твой номер: ", phone, "!"))
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
	url := utils.SumStrings(h.bot.APIUrl, "/user")
	parsedURL, err := utils.SanitizeURL(url)
	req, err := http.NewRequest("GET", parsedURL, nil)
	if err != nil {
		return ctx.Send("Произошла ошибка при создании запроса: ", err.Error())
	}

	req.Header.Add("Authorization", utils.SumStrings("Bearer ", h.bot.APIToken))

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

	btnRefill := tele.InlineButton{
		Text: "Пополнить",
		URL:  "https://gen-api.ru/account/billing",
	}

	menu := &tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{{btnRefill}}}
	return ctx.Send(utils.SumStrings("💰 Баланс GenAPI: ", message.Balance, "₽"), menu)
}
