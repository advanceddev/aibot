package bot

import tele "gopkg.in/telebot.v4"

// UnrealBot - объект бота
type UnrealBot struct {
	Bot                   *tele.Bot
	BotID                 string
	APIToken              string
	APIUrl                string
	PaymentProviderAPIKey string
	ChannelID             int64
	AdminUserID           int64
}
