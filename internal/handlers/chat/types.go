package chat

import (
	"unrealbot/cmd/bot"
)

// Handler - структура обработчика
type Handler struct {
	bot *bot.UnrealBot
}

// GenAIApiResponse - Структура ответа GenAI API
type GenAIApiResponse struct {
	RequestID int                    `json:"request_id"` // Идентификатор запроса
	Model     string                 `json:"model"`      // Модель, используемая для запроса
	Cost      float64                `json:"cost"`       // Стоимость запроса
	Response  []GenAIApiResponseBody `json:"response"`   // Ответы от модели
}

// GenAIApiResponseBody - Структура тела ответа
type GenAIApiResponseBody struct {
	FinishReason string       `json:"finish_reason"` // Причина завершения ответа
	Index        int          `json:"index"`         // Индекс ответа
	Message      GenAIMessage `json:"message"`       // Сообщение
}

// GenAIMessage - Структура сообщения
type GenAIMessage struct {
	Content string `json:"content"` // Содержимое сообщения
	Role    string `json:"role"`    // Роль отправителя сообщения
}
