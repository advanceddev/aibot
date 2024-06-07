package bot

import (
	"log"

	tele "gopkg.in/telebot.v3"
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

// ContactHandler - обработчик входящих контактных данных пользователя
func (bot *UnrealBot) ContactHandler(ctx tele.Context) error {

	ctx.Notify("typing")
	phone := ctx.Message().Contact.PhoneNumber
	ctx.Send("Записал твой номер: " + phone + "!")

	return nil
}
