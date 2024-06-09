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
	Name            string `json:"name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Balance         string `json:"balance"`
	EmailVerifiedAt string `json:"email_verified_at"`
	CreatedAt       string `json:"created_at"`
}

// GenAIApiResponse - структура ответа GenAI API
type GenAIApiResponse struct {
	RequestID           int                  `json:"request_id"`
	Model               string               `json:"model"`
	Choices             []GenAIApiChoice     `json:"choices"`
	PromptFilterResults []PromptFilterResult `json:"prompt_filter_results"`
	Usage               GenAIUsage           `json:"usage"`
}

// GenAIUsage - структура использования GenAI API
type GenAIUsage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// PromptFilterResult - структура поля prompt_filter_results
type PromptFilterResult struct {
	PromptIndex          int                       `json:"prompt_index"`
	ContentFilterResults GenAIContentFilterResults `json:"content_filter_results"`
}

// GenAIApiChoice - структура варианта ответа
type GenAIApiChoice struct {
	ContentFilterResults GenAIContentFilterResults `json:"content_filter_results"`
	FinishReason         string                    `json:"finish_reason"`
	Index                int                       `json:"index"`
	Logprobs             interface{}               `json:"logprobs"`
	Message              GenAIMessage              `json:"message"`
}

// GenAIContentFilterResults - структура результатов фильтрации контента
type GenAIContentFilterResults struct {
	Hate     GenAIFilterType `json:"hate"`
	SelfHarm GenAIFilterType `json:"self_harm"`
	Sexual   GenAIFilterType `json:"sexual"`
	Violence GenAIFilterType `json:"violence"`
}

// GenAIFilterType - структура фильтрации типа контента
type GenAIFilterType struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity"`
}

// GenAIMessage - структура сообщения
type GenAIMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}
