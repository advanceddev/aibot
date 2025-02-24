package bot

import (
	"fmt"
	"time"
	"unrealbot/internal/config"

	tele "gopkg.in/telebot.v4"
)

// InitBot - инициализация подключения к телеграм-боту
func InitBot(cfg *config.Config) (*UnrealBot, error) {
	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать инстанс бота: %w", err)
	}

	return &UnrealBot{
		APIToken:          cfg.APIToken,
		APIUrl:            cfg.APIUrl,
		Bot:               b,
		BotID:             cfg.BotID,
		ChannelID:         cfg.ChannelID,
		AdminUserID:       cfg.AdminUserID,
		AiModelIdentifier: cfg.AiModelIdentifier,
	}, nil
}
