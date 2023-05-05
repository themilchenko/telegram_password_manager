package main

import (
	"log"
	"os"

	tg "telegram_password_manager/internal/tg"
	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_API_TOKEN"))
	// if err != nil {
	// 	log.Fatal("Token is incorrect")
	// }
	//
	// bot.Debug = true
	//
	// updateConfig := tgbotapi.NewUpdate(0)
	// updateConfig.Timeout = 30

	bot, err := tg.NewBot(os.Getenv("TG_API_TOKEN"))
	if err != nil {
		log.Fatalf("%s", err.Error)
	}

	if err = bot.Run(); err != nil {
		log.Fatalf("%s", err.Error)
	}
	// updates := bot.GetUpdatesChan(updateConfig)

	// for update := range updates {
	// 	// We handle only messages and commands at this moment
	// 	if update.Message == nil || !update.Message.IsCommand() {
	// 		continue
	// 	}
	//
	// 	// Now that we know we've gotten a new message, we can construct a
	// 	// reply! We'll take the Chat ID and Text from the incoming message
	// 	// and use it to create a new message.
	// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	//
	// 	switch update.Message.Command() {
	// 	case "set":
	// 		msg.Text = "ok set"
	// 	case "get":
	// 		msg.Text = "ok get"
	// 	case "del":
	// 		msg.Text = "ok del"
	// 	default:
	// 		msg.Text = "I don't know that command"
	// 	}
	//
	// 	// Okay, we're sending our message off! We don't care about the message
	// 	// we just sent, so we'll discard it.
	// 	if _, err := bot.Send(msg); err != nil {
	// 		// Note that panics are a bad way to handle errors. Telegram can
	// 		// have service outages or network errors, you should retry sending
	// 		// messages or more gracefully handle failures.
	// 		panic(err)
	// 	}
	// }
}
