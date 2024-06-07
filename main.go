package main

import (
	"unrealbot/cmd/bot"
	"unrealbot/config"
	"unrealbot/internal/chat"
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

	// /-- Только для подписчиков
	memberOnly := unrealBot.Bot.Group()

	memberOnly.Use(CheckMembership)

	memberOnly.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет, "+c.Sender().FirstName+"! 👋", emptyMenu)
	})

	memberOnly.Handle(tele.OnContact, unrealBot.ContactHandler)
	memberOnly.Handle(tele.OnText, messageHandler.HandleMessage)
	// --/

	unrealBot.Bot.Handle(tele.OnCheckout, checkoutHandler.HandleCheckout)

	unrealBot.Bot.Handle(&btnPay, invoiceHandler.HandleInvoice)

	unrealBot.Bot.Handle(&btnSubscribe, func(c tele.Context) error {
		channel := &tele.Chat{ID: channelID, Type: "privatechannel"}
		link, err := c.Bot().InviteLink(channel)
		if err != nil {
			return c.Send("Произошла ошибка при формировании пригласительной ссылки.")
		}
		return c.Send(link)
	})

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
	// --/

}

// --- Мидлвейры --- /

// CheckMembership - мидлвейр, проверяет подписку на канал
func CheckMembership(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user := c.Recipient()
		channel := &tele.Chat{ID: channelID}

		chatMember, err := c.Bot().ChatMemberOf(channel, user)

		if err != nil {
			return c.Send("Ошибка при проверке подписки: " + err.Error())
		}

		if isMember(chatMember) {
			return next(c)
		}

		// return c.Send("У вас нет доступа к этому боту.", guestMenu)
		return c.Send("У вас нет доступа к этому боту. Свяжитесь с администратором: @frntbck", emptyMenu)
	}
}

func isMember(chatMember *tele.ChatMember) bool {
	return chatMember.Role == tele.Member || chatMember.Role == tele.Administrator || chatMember.Role == tele.Creator
}

// --/
