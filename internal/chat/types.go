package chat

import (
	"unrealbot/cmd/bot"
)

// Handler - структура обработчика
type Handler struct {
	bot *bot.UnrealBot
}

// GenAIApiResponse - структура ответа GenAI API
type GenAIApiResponse struct {
	RequestID           int                  `json:"request_id"`
	Model               string               `json:"model"`
	Choices             []GenAIApiChoice     `json:"choices"`
	PromptFilterResults []PromptFilterResult `json:"prompt_filter_results"`
	Usage               struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// PromptFilterResult - структура поля prompt_filter_results
type PromptFilterResult struct {
	PromptIndex          int `json:"prompt_index"`
	ContentFilterResults struct {
		Hate struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"hate"`
		SelfHarm struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"self_harm"`
		Sexual struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"sexual"`
		Violence struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"violence"`
	} `json:"content_filter_results"`
}

// GenAIApiChoice - структура варианта ответа
type GenAIApiChoice struct {
	ContentFilterResults struct {
		Hate struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"hate"`
		SelfHarm struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"self_harm"`
		Sexual struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"sexual"`
		Violence struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"violence"`
	} `json:"content_filter_results"`
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	Message      struct {
		Content string `json:"content"`
		Role    string `json:"role"`
	} `json:"message"`
}

// CozeResponseMessage - структура сообщения в ответе Coze API
type CozeResponseMessage struct {
	Role        string `json:"role"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

// CozeAPIResponse - структура ответа Coze API
type CozeAPIResponse struct {
	Messages       []CozeResponseMessage `json:"messages"`
	ConversationID string                `json:"conversation_id"`
	Code           int                   `json:"code"`
	Msg            string                `json:"msg"`
}