package main

import (
	"unrealbot/cmd/bot"
	"unrealbot/config"
	"unrealbot/internal/chat"
	"unrealbot/internal/middlewares"
	"unrealbot/internal/payments"

	tele "gopkg.in/telebot.v3"
)

func main() {

	cfg := config.MustLoad()

	unrealBot := bot.UnrealBot{
		APIToken:              cfg.APIToken,
		APIUrl:                cfg.APIUrl,
		Bot:                   bot.InitBot(cfg.BotToken),
		BotID:                 cfg.BotID,
		PaymentProviderAPIKey: cfg.PaymentProviderAPIKey,
		ChannelID:             cfg.ChannelID,
		AdminUserID:           cfg.AdminUserID,
	}

	defer unrealBot.Bot.Stop()

	registerHandlers(unrealBot)
	unrealBot.Bot.Start()
}

func registerHandlers(unrealBot bot.UnrealBot) {

	checkoutHandler := payments.NewCheckoutHandler(&unrealBot)
	invoiceHandler := payments.NewInvoiceHandler(&unrealBot)
	messageHandler := chat.NewMessageHandler(&unrealBot)
	commandHandler := chat.NewCommandHandler(&unrealBot)

	// Создаем группу хэндлеров и добавляем мидлвэйр
	memberOnly := unrealBot.Bot.Group()
	memberOnly.Use(middlewares.CheckMembership(unrealBot.ChannelID, unrealBot.AdminUserID))

	// Хэндлеры группы membersOnly
	memberOnly.Handle("/start", commandHandler.StartHandler)
	memberOnly.Handle(tele.OnContact, commandHandler.ContactHandler)
	memberOnly.Handle(tele.OnText, messageHandler.HandleMessage)

	// Публичные хэндлеры
	unrealBot.Bot.Handle(tele.OnCheckout, checkoutHandler.HandleCheckout)
	unrealBot.Bot.Handle(tele.OnPayment, invoiceHandler.HandlePaymentSuccess)

}
