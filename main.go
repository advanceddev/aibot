package main

import (
	"unrealbot/cmd/bot"
	"unrealbot/config"
	"unrealbot/internal/chat"
	"unrealbot/internal/middlewares"
	"unrealbot/internal/payments"

	tele "gopkg.in/telebot.v3"
)

type checkoutHandler struct {
	bot *tele.Bot
}

var (
	channelID    int64
	menu         = &tele.ReplyMarkup{ResizeKeyboard: true}
	guestMenu    = &tele.ReplyMarkup{ResizeKeyboard: true}
	emptyMenu    = &tele.ReplyMarkup{RemoveKeyboard: true}
	btnPromo     = menu.URL("💥 The Absolute Basstards", "@tabdnb")
	btnPay       = menu.Text("📢 Оплатить подписку")
	btnSubscribe = menu.Text("🎸 Подписаться")
)

func main() {

	cfg := config.MustLoad()

	channelID = cfg.ChannelID

	unrealBot := bot.UnrealBot{
		APIToken:              cfg.APIToken,
		APIUrl:                cfg.APIUrl,
		Bot:                   bot.InitBot(cfg.BotToken),
		BotID:                 cfg.BotID,
		PaymentProviderAPIKey: cfg.PaymentProviderAPIKey,
	}

	defer unrealBot.Bot.Stop()

	setupMenu()

	registerHandlers(unrealBot)
	unrealBot.Bot.Start()
}

func setupMenu() {
	menu.Reply(
		menu.Row(btnPromo),
	)
	guestMenu.Reply(
		// Кнопка для оплаты
		// guestMenu.Row(btnPay),

		// Кнопка для подписки
		guestMenu.Row(btnSubscribe),
	)
}

func registerHandlers(unrealBot bot.UnrealBot) {

	checkoutHandler := payments.NewCheckoutHandler(&unrealBot)
	invoiceHandler := payments.NewInvoiceHandler(&unrealBot)
	messageHandler := chat.NewMessageHandler(&unrealBot)
	commandHandler := chat.NewCommandHandler(&unrealBot)

	// Создаем группу хэндлеров
	memberOnly := unrealBot.Bot.Group()

	// Добавляем мидлвару к группе хэндлеров
	memberOnly.Use(middlewares.CheckMembership)

	// Хэндлеры группы membersOnly
	memberOnly.Handle("/start", commandHandler.StartHandler)
	memberOnly.Handle(tele.OnContact, commandHandler.ContactHandler)
	memberOnly.Handle(tele.OnText, messageHandler.HandleMessage)

	unrealBot.Bot.Handle(tele.OnCheckout, checkoutHandler.HandleCheckout)
	unrealBot.Bot.Handle(&btnPay, invoiceHandler.HandleInvoice)

	// TODO: вынести хэндлер в отдельный модуль
	unrealBot.Bot.Handle(&btnSubscribe, func(c tele.Context) error {
		channel := &tele.Chat{ID: channelID, Type: "privatechannel"}
		link, err := c.Bot().InviteLink(channel)
		if err != nil {
			return c.Send("Произошла ошибка при формировании пригласительной ссылки.")
		}
		return c.Send(link)
	})

	// TODO: вынести хэндлер в отдельный модуль
	// /-- В случае успешного платежа отправляем уникальную пригласительную ссылку на канал
	unrealBot.Bot.Handle(tele.OnPayment, func(c tele.Context) error {
		if c.Message().Payment != nil {
			channel := &tele.Chat{ID: channelID, Type: "privatechannel"}
			link, err := c.Bot().InviteLink(channel)
			if err != nil {
				return c.Send("Произошла ошибка при формировании пригласительной ссылки.")
			}
			return c.Send(link)
		}
		return nil
	})

}
