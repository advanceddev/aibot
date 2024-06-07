package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"unrealbot/cmd/bot"

	tele "gopkg.in/telebot.v3"
)

// Строки с рекламными сообщениями
var promoStrings = [10]string{
	"\n\n------\n\n💥📢 The Absolute Basstards - Drum&Bass label.\nhttps://t.me/+bS_eIEhkLuZkMDYy",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
}

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

	req, err := createPostRequest(h.bot.GenAPIUrl, h.bot.GenAPIKey, payloadBytes)
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
	randomIndex := rand.Intn(len(promoStrings))
	promoMessage := promoStrings[randomIndex]

	// Отправляем сообщение и рекламу
	return ctx.Send(answerContent + promoMessage)
}

// createPostRequest создает POST запрос к API
func createPostRequest(url, apiKey string, payload []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
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

// CozeMessageHandler - обработчик текстовых сообщений для Coze
// func (bot *UnrealBot) CozeMessageHandler(ctx tele.Context) error {
// 	ctx.Notify("typing")
// 	message := ctx.Message().Text
// 	userID := ctx.Sender().ID
// 	url := bot.APIUrl

// 	payload := map[string]interface{}{
// 		"bot_id": bot.BotID,
// 		"user":   string(rune(userID)),
// 		"query":  message,
// 		"stream": false,
// 	}

// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return fmt.Errorf("Ошибка маршалинга: %w", err)
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return fmt.Errorf("Не удалось создать запрос: %w", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("User-Agent", "insomnia/9.2.0")
// 	req.Header.Set("Authorization", "Bearer "+bot.APIToken)

// 	ctx.Notify("typing")

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("Не удалось выполнить запрос: %w", err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return fmt.Errorf("Ответ с ошибкой: %s", res.Status)
// 	}

// 	ctx.Notify("typing")

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return fmt.Errorf("Не удалось прочитать ответ: %w", err)
// 	}

// 	var apiResponse CozeAPIResponse
// 	if err := json.Unmarshal(body, &apiResponse); err != nil {
// 		return fmt.Errorf("Не удалось декодировать запрос: %w", err)
// 	}

// 	ctx.Notify("typing")

// 	answerContent := ""
// 	for _, msg := range apiResponse.Messages {
// 		if msg.Type == "answer" {
// 			answerContent = msg.Content
// 			break
// 		}
// 	}

// 	if answerContent == "" {
// 		if err := ctx.Send("К сожалению, я не знаю, что ответить... :("); err != nil {
// 			return fmt.Errorf("Ответа нет и не удалось отправить: %w", err)
// 		}
// 		if err := ctx.Send(ctx.Sender().ID); err != nil {
// 			return fmt.Errorf("failed to send user ID: %w", err)
// 		}
// 		return nil
// 	}

// 	if err := ctx.Send(answerContent); err != nil {
// 		return fmt.Errorf("Ответ есть, но не удалось отправить: %w", err)
// 	}

// 	return nil
// }
