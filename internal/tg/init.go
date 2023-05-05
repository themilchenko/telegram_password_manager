package tg

import (
	"errors"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
) 

type Telegram struct {
	bot *tgbot.BotAPI
	updCfg *tgbot.UpdateConfig
}

func NewBot(tokenAPI string) (*Telegram, error) {
	bot, err := tgbot.NewBotAPI(tokenAPI)
	if err != nil {
		return nil, errors.New("api key is incorrect")
	}

	bot.Debug = true
	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30

	return &Telegram{
		bot: bot,
		updCfg: &updateConfig,
	}, nil
}

func (b *Telegram) Run() error {
	updates := tgbot.GetUpdatesChan(b.updCfg)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "set":
			msg.Text = "ok set"
		case "get":
			msg.Text = "ok get"
		case "del":
			msg.Text = "ok del"
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			return err
		}
	}
}
