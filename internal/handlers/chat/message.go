package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unrealbot/cmd/bot"
	"unrealbot/internal/utils"

	tele "gopkg.in/telebot.v3"
)

// NewMessageHandler создает новый обработчик сообщений
func NewMessageHandler(bot *bot.UnrealBot) *Handler {
	return &Handler{bot: bot}
}

// HandleMessage обрабатывает текстовое сообщение пользователя
func (h *Handler) HandleMessage(ctx tele.Context) error {
	message := ctx.Message().Text

	payload := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": message,
					},
				},
			},
		},
		"is_sync": true,
		"model":   "gpt-4o-2024-05-13",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Ошибка маршалинга: %w", err)
	}

	req, err := createPostRequest(h.bot.APIUrl, h.bot.APIToken, payloadBytes)
	if err != nil {
		return err
	}
	ctx.Notify("typing")
	res, err := sendRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Ответ с ошибкой: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Не удалось прочитать ответ: %w", err)
	}

	apiResponse, err := parseAPIResponse(body)
	if err != nil {
		return err
	}

	answerContent := findAssistantMessage(apiResponse)
	if answerContent == "" {
		return handleNoAnswer(ctx)
	}

	// Выбираем случайную рекламную строку
	promoMessage := GetPromoString()

	// Отправляем сообщение и рекламу
	return ctx.Send(answerContent + promoMessage)
}

// createPostRequest создает POST запрос к API
func createPostRequest(apiURL, apiKey string, payload []byte) (*http.Request, error) {
	parsedURL, err := utils.SanitizeURL(apiURL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", parsedURL+"/networks/chat-gpt-4-turbo", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать запрос: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return req, nil
}

// sendRequest отправляет запрос и возвращает ответ
func sendRequest(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

// parseApiResponse декодирует ответ API в структуру GenAIApiResponse
func parseAPIResponse(body []byte) (*GenAIApiResponse, error) {
	var apiResponse GenAIApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("Не удалось декодировать запрос: %w", err)
	}
	return &apiResponse, nil
}

// findAssistantMessage находит сообщение от ассистента в API-ответе
func findAssistantMessage(apiResponse *GenAIApiResponse) string {
	for _, choice := range apiResponse.Choices {
		if choice.Message.Role == "assistant" {
			return choice.Message.Content
		}
	}
	return ""
}

// handleNoAnswer обрабатывает ситуацию, когда ответ не найден
func handleNoAnswer(ctx tele.Context) error {
	if err := ctx.Send("К сожалению, я не знаю, что ответить... :("); err != nil {
		return fmt.Errorf("Ответа нет и не удалось отправить: %w", err)
	}

	if err := ctx.Send(ctx.Sender().ID); err != nil {
		return fmt.Errorf("failed to send user ID: %w", err)
	}

	return nil
}
