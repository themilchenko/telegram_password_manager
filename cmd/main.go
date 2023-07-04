package main

import (
	"fmt"
	"log"
	"os"

	tg "telegram_password_manager/internal/tg"
)

func main() {
	// Create new instance of bot using API key
	fmt.Println("Hello there")

	bot, err := tg.NewBot(os.Getenv("TG_API_TOKEN"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Run our bot to communicate with
	if err = bot.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
