package main

import (
	"unrealbot/cmd/bot"
	"unrealbot/internal/config"
	"unrealbot/internal/handlers/chat"
	"unrealbot/internal/middlewares"

	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func main() {

	cfg := config.MustLoad()

	unrealBot := bot.UnrealBot{
		APIToken:    cfg.APIToken,
		APIUrl:      cfg.APIUrl,
		Bot:         bot.InitBot(cfg.BotToken),
		BotID:       cfg.BotID,
		ChannelID:   cfg.ChannelID,
		AdminUserID: cfg.AdminUserID,
		AiModelIdentifier: cfg.AiModelIdentifier,
	}

	defer unrealBot.Bot.Stop()

	registerHandlers(unrealBot)
	unrealBot.Bot.Start()

}

func registerHandlers(unrealBot bot.UnrealBot) {

	messageHandler := chat.NewMessageHandler(&unrealBot)
	commandHandler := chat.NewCommandHandler(&unrealBot)

	// Создаем группу для админов и закрываем ее мидлвэйром
	adminOnly := unrealBot.Bot.Group()
	adminOnly.Use(middleware.Whitelist(unrealBot.AdminUserID))

	// Создаем группу для пользователей и закрываем ее мидлвэйром
	memberOnly := unrealBot.Bot.Group()
	memberOnly.Use(middlewares.CheckMembership(unrealBot)) // Проверить подписку и запросить доступ

	// Хэндлеры группы membersOnly
	memberOnly.Handle("/start", commandHandler.StartHandler)     // Обработчик команды /start
	memberOnly.Handle(tele.OnText, messageHandler.HandleMessage) // Обработчик текстового сообщения

}
