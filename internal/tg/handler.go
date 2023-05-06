package tg

import (
	"log"
	"telegram_password_manager/internal/domain"
	"telegram_password_manager/internal/models"

	// "telegram_password_manager/internal/domain"

	bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Messages for responses
const (
	privateKeyMsg = "Create your password to get access to the password manager. Remember that this password will be a key for all your passwords that you will communicate with. Dont loose it."
)


type Handler struct {
	bot *bot.BotAPI
	usecase domain.Usecase
}

func NewHandler(u domain.Usecase, b *bot.BotAPI) *Handler {
	return &Handler{
		bot: b,
		usecase: u,
	}
}

// There is a router that will handle commads
func (h *Handler) HandleCommand(update *bot.Update, states map[int64]int) {
	chatID := update.Message.Chat.ID
	if _, ok := states[chatID]; !ok {
		states[chatID] = int(models.StateDeafault)
	} 

	switch (states[chatID]) {
	case int(models.StateDeafault):
		// Handling command request and wait for secret key
		var msg string
		switch update.Message.Command() {
		case "set":
			states[chatID] = int(models.StateWaitSetKey)
			msg = privateKeyMsg
		case "get":
			states[chatID] = int(models.StateWaitGetKey)
			msg = privateKeyMsg
		case "del":
			states[chatID] = int(models.StateWaitDelKey)
			msg = privateKeyMsg
		default:
			msg = "I don't know that command"
		}
		if err := h.SendMsg(update, msg); err != nil {
			log.Println(err.Error())
		}
	case int(models.StateWaitGetKey):
	case int(models.StateRightAdd):
	case int(models.StateRightGet):
	case int(models.StateRightDel):
	default:
	}
}

func (h *Handler) SendMsg(update *bot.Update, text string) error {
	msg := bot.NewMessage(update.Message.Chat.ID, "")
	msg.Text = privateKeyMsg
	if _, err := h.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Set(msg *bot.MessageConfig, chatID int64) error {
	msg.Text = privateKeyMsg

	if _, err := h.bot.Send(msg); err != nil {
		log.Println(err.Error())
	}

	passwordMsg := bot.NewMessage(chatID, "")

	if _, err := h.bot.Send(passwordMsg); err != nil {
		log.Println(err.Error())
	}
	
	return nil
}

func (h *Handler) Get(msg *bot.MessageConfig, chatID int64) error {
	return nil
}

func (h *Handler) Remove(msg *bot.MessageConfig, chatID int64) error {
	return nil
}
