package chat

import (
	"unrealbot/cmd/bot"
)

// Handler - структура обработчика
type Handler struct {
	bot *bot.UnrealBot
}

// GenAIUserResponse - Структура ответа информации об аккаунте GenAI
type GenAIUserResponse struct {
	Name            string `json:"name"`              // Имя пользователя
	Email           string `json:"email"`             // Электронная почта пользователя
	PhoneNumber     string `json:"phone_number"`      // Номер телефона пользователя
	Balance         string `json:"balance"`           // Баланс пользователя
	EmailVerifiedAt string `json:"email_verified_at"` // Дата и время верификации электронной почты
	CreatedAt       string `json:"created_at"`        // Дата и время создания аккаунта
}

// GenAIApiResponse - Структура ответа GenAI API
type GenAIApiResponse struct {
	RequestID int                     `json:"request_id"` // Идентификатор запроса
	Model     string                  `json:"model"`      // Модель, используемая для запроса
	Cost      float64                 `json:"cost"`       // Стоимость запроса
	Response  []GenAIUserResponseBody `json:"response"`   // Ответы от модели
	Usage     GenAIUsage              `json:"usage"`      // Информация об использовании
}

// GenAIUsage - Структура использования GenAI API
type GenAIUsage struct {
	CompletionTokens int `json:"completion_tokens"` // Количество токенов в ответе
	PromptTokens     int `json:"prompt_tokens"`     // Количество токенов в запросе
	TotalTokens      int `json:"total_tokens"`      // Общее количество токенов
}

// GenAIUserResponseBody - Структура варианта ответа
type GenAIUserResponseBody struct {
	FinishReason string       `json:"finish_reason"` // Причина завершения ответа
	Index        int          `json:"index"`         // Индекс ответа
	Message      GenAIMessage `json:"message"`       // Сообщение
}

// GenAIMessage - Структура сообщения
type GenAIMessage struct {
	Content string `json:"content"` // Содержимое сообщения
	Role    string `json:"role"`    // Роль отправителя сообщения
}
