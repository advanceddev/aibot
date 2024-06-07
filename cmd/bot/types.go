package bot

import tele "gopkg.in/telebot.v3"

// UnrealBot - объект бота
type UnrealBot struct {
	Bot                   *tele.Bot
	BotID                 string
	APIToken              string
	APIUrl                string
	PaymentProviderAPIKey string
}



