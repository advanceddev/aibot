package main

import (
	"log"
	"unrealbot/cmd/bot"
	"unrealbot/internal/config"
	"unrealbot/internal/handlers/chat"
	"unrealbot/internal/middlewares"

	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func main() {

	cfg := config.MustLoad()

	unrealBot, err := bot.InitBot(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	defer unrealBot.Bot.Stop()

	// registerHandlers(unrealBot)
	if err := registerHandlers(unrealBot); err != nil {
		log.Fatalf("Не удалось зарегистрировать обработчики: %v", err)
	}
	unrealBot.Bot.Start()

}

func registerHandlers(unrealBot *bot.UnrealBot) error {

	messageHandler := chat.NewMessageHandler(unrealBot)
	commandHandler := chat.NewCommandHandler(unrealBot)

	// Создаем группу для админов и закрываем ее мидлвэйром
	adminOnly := unrealBot.Bot.Group()
	adminOnly.Use(middleware.Whitelist(unrealBot.AdminUserID))

	// Создаем группу для пользователей и закрываем ее мидлвэйром
	memberOnly := unrealBot.Bot.Group()
	memberOnly.Use(middlewares.CheckMembership(unrealBot)) // Проверить подписку и запросить доступ

	// Хэндлеры группы membersOnly
	memberOnly.Handle("/start", commandHandler.StartHandler)     // Обработчик команды /start
	memberOnly.Handle(tele.OnText, messageHandler.HandleMessage) // Обработчик текстового сообщения

	return nil
}
