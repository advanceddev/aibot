package payments

import (
	"encoding/json"
	"fmt"
	"unrealbot/cmd/bot"

	tele "gopkg.in/telebot.v3"
)

// NewInvoiceHandler - конструктор обработчика инвойса
func NewInvoiceHandler(bot *bot.UnrealBot) *Handler {
	return &Handler{bot: bot}
}

// HandlePaymentSuccess - обработчик успешного платежа
func (h *Handler) HandlePaymentSuccess(ctx tele.Context) error {
	if ctx.Message().Payment != nil {
		channel := &tele.Chat{ID: h.bot.ChannelID, Type: "privatechannel"}
		link, err := ctx.Bot().InviteLink(channel)
		if err != nil {
			return ctx.Send("Произошла ошибка при формировании пригласительной ссылки.")
		}
		return ctx.Send(link)
	}
	return nil
}

// HandleInvoice - создание и отправка инвойса
func (h *Handler) HandleInvoice(ctx tele.Context) error {
	paymentID := int64(1) // here you need to assign a payment ID
	telegramUserID := ctx.Sender().ID

	_, err := h.createAndSendInvoice(paymentID, telegramUserID, ctx.Sender())
	if err != nil {
		return fmt.Errorf("handle invoice: %w", err)
	}

	return nil
}

func (h *Handler) createAndSendInvoice(paymentID, userTelegramID int64, sender *tele.User) (*tele.Invoice, error) {
	invoice, err := CreateInvoice(paymentID, userTelegramID, h.bot.APIToken)
	if err != nil {
		return nil, fmt.Errorf("create invoice: %w", err)
	}

	_, err = invoice.Send(h.bot.Bot, sender, nil)
	if err != nil {
		return nil, fmt.Errorf("send invoice: %w", err)
	}

	return invoice, nil
}

// CreateInvoice - формирование инвойса
func CreateInvoice(paymentID, userTelegramID int64, apiKey string) (*tele.Invoice, error) {
	payload := InvoiceData{
		SystemPaymentID: paymentID,
		UserTelegramID:  userTelegramID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Ошибка маршалинга: %w", err)
	}

	const totalAmount = invoiceAmount * 100
	invoice := tele.Invoice{
		Title:       invoiceTitle,
		Description: invoiceDescription,
		Currency:    invoiceCurrency,
		Prices: []tele.Price{
			{
				Label:  pricesLabel,
				Amount: totalAmount,
			},
		},
		Token:   apiKey,
		Total:   totalAmount,
		Payload: string(payloadBytes),
	}
	return &invoice, nil
}
