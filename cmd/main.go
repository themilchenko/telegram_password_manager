package main

import (
	"fmt"
	"log"

	tg "telegram_password_manager/internal/tg"
)

func main() {
	// Create new instance of bot using API key
	fmt.Println("Hello there")
	bot, err := tg.NewBot("6004504657:AAEk0sk69zTH8EP1WuleOhnOU4_qJ3Ig6p4")
	// bot, err := tg.NewBot(os.Getenv("TG_API_TOKEN"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Run our bot to communicate with
	if err = bot.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
