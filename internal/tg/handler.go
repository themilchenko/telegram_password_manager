package tg

import (
	"log"
	"telegram_password_manager/internal/domain"
	"telegram_password_manager/internal/models"

	bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Messages for responses
const (
	privateKeyMsg = "Create your password to get access to the password manager. Remember that this password will be a key for all your passwords that you will communicate with. Dont loose it."
)

type Handler struct {
	bot     *bot.BotAPI
	usecase domain.Usecase
}

func NewHandler(u domain.Usecase, b *bot.BotAPI) *Handler {
	return &Handler{
		bot:     b,
		usecase: u,
	}
}

// There is a router that will handle commads
func (h *Telegram) HandleCommand(update *bot.Update, states map[int64]int) {
	chatID := update.Message.Chat.ID

	state, err := h.usecase.GetStateByChatID(chatID)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			err = h.usecase.CreateStateByChatID(chatID)
		} else {
			log.Print(err.Error())
		}
	}

	switch state.ChatState {
	case models.StateDeafault:
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
	case models.StateWaitGetKey:
	case models.StateRightAdd:
	case models.StateRightGet:
	case models.StateRightDel:
	default:
	}
}

func (h *Telegram) SendMsg(update *bot.Update, text string) error {
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
