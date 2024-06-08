package payments

import (
	"encoding/json"
	"fmt"
	"log"
	"unrealbot/cmd/bot"

	tele "gopkg.in/telebot.v3"
)

// NewCheckoutHandler - обработчик платежа
func NewCheckoutHandler(bot *bot.UnrealBot) *Handler {
	return &Handler{bot: bot}
}

// handleCheckoutError - обработка ошибок при проведении платежа
func (h *Handler) handleCheckoutError(err error, pre *tele.PreCheckoutQuery) {
	log.Println(err)
	h.sendErrorMessage(pre)
}

// sendErrorMessage - отправка сообщения об ошибке пользователю
func (h *Handler) sendErrorMessage(pre *tele.PreCheckoutQuery) {
	_, sendErr := h.bot.Bot.Send(
		pre.Sender,
		fmt.Sprintf(
			"Ошибка при проведении платежа. "+
				"Свяжитесь с администратором.\n\n Счет ID: %s",
			pre.ID,
		),
	)
	if sendErr != nil {
		log.Println(sendErr)
	}
}

// HandleCheckout - обработчик платежа
func (h *Handler) HandleCheckout(ctx tele.Context) error {
	pre := ctx.PreCheckoutQuery()
	if pre == nil || pre.ID == "" {
		return nil
	}

	_, err := h.decodeInvoice(pre.Payload)
	if err != nil {
		h.handleCheckoutError(err, pre)
		return nil
	}

	if err := h.bot.Bot.Accept(pre); err != nil {
		h.handleCheckoutError(fmt.Errorf("Не удалось подтвердить платеж: %w", err), pre)
		return nil
	}

	return nil
}

// decodeInvoice - декодирование данных счета из JSON
func (h *Handler) decodeInvoice(payload string) (InvoiceData, error) {
	var inv InvoiceData
	err := json.Unmarshal([]byte(payload), &inv)
	if err != nil {
		return inv, fmt.Errorf("Ошибка декодированя счета: %w", err)
	}
	return inv, nil
}
