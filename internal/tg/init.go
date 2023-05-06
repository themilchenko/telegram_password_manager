package tg

import (
	"errors"

	domain "telegram_password_manager/internal/domain"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tarantool/go-tarantool"
) 

type Telegram struct {
	userStates map[int64]int

	bot *tgbot.BotAPI
	updCfg tgbot.UpdateConfig

	db *DB

	Router *Handler
	Usecase Usecase
	// Usecase
	UsecaseDomain domain.Usecase
	// RepoDomain domain.Repository
}

func NewBot(tokenAPI string) (*Telegram, error) {
	userStates := make(map[int64]int, 0)

	// Create connection to tarantool
	db, err := NewDB("127.0.0.1:3031", &tarantool.Opts{
		User: "guest",
	})
	if err != nil {
		return nil, err
	}

	// Create conncetion to Telegram API
	bot, err := tgbot.NewBotAPI("6004504657:AAEk0sk69zTH8EP1WuleOhnOU4_qJ3Ig6p4")
	if err != nil {
		return nil, errors.New("api key is incorrect")
	}
	bot.Debug = true
	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30

	// Create instanse of server
	tg := &Telegram{
		userStates: userStates,
		bot: bot,
		updCfg: updateConfig,
		db: db,
	}

	tg.Router = NewHandler(tg.UsecaseDomain, bot)
	tg.Usecase = NewUsecase(db)
	return tg, nil
}

// Here we listen for updates and handle its
func (b *Telegram) Run() error {
	updates := b.bot.GetUpdatesChan(b.updCfg)

	for update := range updates {
		// We ignore empty messages
		if update.Message != nil {
			b.Router.HandleCommand(&update, b.userStates)
		}
	}
	
	return nil
}
