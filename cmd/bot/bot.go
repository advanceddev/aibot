package bot

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

// InitBot - инициализация бота и его подключение к телеграму
func InitBot(token string) *tele.Bot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота %v", err)
	}

	return b
}
