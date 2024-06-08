package payments

import (
	"unrealbot/cmd/bot"
)

// Handler - структура обработчика
type Handler struct {
	bot *bot.UnrealBot
}

// InvoiceData - данные инвойса
type InvoiceData struct {
	SystemPaymentID int64 `json:"paymentID"`
	UserTelegramID  int64 `json:"tid"`
}
