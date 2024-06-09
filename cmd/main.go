package main

import (
	"unrealbot/cmd/bot"
	"unrealbot/internal/config"
	"unrealbot/internal/handlers/chat"
	"unrealbot/internal/handlers/payments"
	"unrealbot/internal/middlewares"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
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

	// Создаем группу для админов и закрываем ее мидлвэйром
	adminOnly := unrealBot.Bot.Group()
	adminOnly.Use(middleware.Whitelist(unrealBot.AdminUserID))
	adminOnly.Handle("/balance", commandHandler.BalanceHandler) // Запросить текущий баланс командой /balance

	// Создаем группу для пользователей и закрываем ее мидлвэйром
	memberOnly := unrealBot.Bot.Group()
	memberOnly.Use(middlewares.CheckMembership(unrealBot)) // Проверить подписку и запросить доступ

	// Хэндлеры группы membersOnly
	memberOnly.Handle("/start", commandHandler.StartHandler)         // Обработчик команды /start
	memberOnly.Handle(tele.OnContact, commandHandler.ContactHandler) // Обработчик на отправленный Контакт
	memberOnly.Handle(tele.OnText, messageHandler.HandleMessage)     // Обработчик текстового сообщения

	// Публичные хэндлеры
	unrealBot.Bot.Handle(tele.OnCheckout, checkoutHandler.HandleCheckout)     // Обработчик созданного платежа
	unrealBot.Bot.Handle(tele.OnPayment, invoiceHandler.HandlePaymentSuccess) // Обработчик успешного платежа

}
