package tg

import (
	"log"
	"telegram_password_manager/internal/domain"
	"telegram_password_manager/internal/models"

	bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Messages for responses
const (
	createPrivateKeyMsg = "Create your password to get access to the password manager. Remember that this password (secret key) will be a key for all your passwords that you will communicate with. Don't loose it."
	privateKeyMsg       = "Enter your secret key"
	incorrectSecretKey  = "Your secret key is incorrect"
	serviceNameMsgSet = "Enter the name of service you want to create password"
	serviceNameMsgDel = "Enter the name of service you want to delete password"
	serviceNameMsgGet = "Enter the name of service you want to get password"
)

type Handler struct {
	usecase domain.Usecase
}

func NewHandler(u domain.Usecase, b *bot.BotAPI) *Handler {
	return &Handler{
		usecase: u,
	}
}

// There is a router that will handle commads
func (h *Telegram) HandleCommand(update *bot.Update) {
	chatID := update.Message.Chat.ID

	state, err := h.usecase.GetStateByChatID(chatID)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			err = h.usecase.CreateStateByChatID(chatID)
		} else {
			log.Print(err.Error())
		}
	}

	var msg string
	var st models.StateType
	switch state.ChatState {
	case models.StateDeafault:
		// Handling command request and wait for secret key
		switch update.Message.Command() {
		case "set":
			err = h.usecase.CheckExistingKey(chatID)
			if err != nil {
				switch err.Error() {
				case "ErrSecKeyCreate":
					st = models.StateCreateSecretKeySet
					msg = createPrivateKeyMsg
				default:
					log.Print(err.Error())
				}
			} else {
				st = models.StateWaitSetKey
				msg = privateKeyMsg
			}
		case "get":
			err = h.usecase.CheckExistingKey(chatID)
			if err != nil {
				switch err.Error() {
				case "ErrSecKeyCreate":
					st = models.StateCreateSecretKeyGet
					msg = createPrivateKeyMsg
				default:
					log.Print(err.Error())
				}
			} else {
				st = models.StateWaitSetKey
				msg = privateKeyMsg
			}
		case "del":
			err = h.usecase.CheckExistingKey(chatID)
			if err != nil {
				switch err.Error() {
				case "ErrSecKeyCreate":
					st = models.StateCreateSecretKeyDel
					msg = createPrivateKeyMsg
				default:
					log.Print(err.Error())
				}
			} else {
				st = models.StateWaitSetKey
				msg = privateKeyMsg
			}
		default:
			msg = "I don't know that command"
		}
	case models.StateWaitGetKey:
		err = h.usecase.ChechSecretKey(state)
		if err != nil {
			switch err.Error() {
			case "ErrSecKeyReqiuered":
				st = models.StateCreateSecretKeyGet
				msg = privateKeyMsg
			case "ErrIncorrectSecretKey":
				st = models.StateDeafault
				msg = incorrectSecretKey
			default:
				log.Print(err.Error())
			}
		}
	case models.StateCreateSecretKeyDel:
		state.SecretKey = update.Message.Text
		err = h.usecase.CreateSecretKey(state)
		if err != nil {
			log.Println(err.Error())
		} else {
			st = models.StateEnterServiceNameDel
			msg = serviceNameMsgDel
		}
	case models.StateCreateSecretKeySet:
		state.SecretKey = update.Message.Text
		err = h.usecase.CreateSecretKey(state)
		if err != nil {
			log.Println(err.Error())
		} else {
			st = models.StateEnterServiceNameSet
			msg = serviceNameMsgSet
		}
	case models.StateCreateSecretKeyGet:
		state.SecretKey = update.Message.Text
		err = h.usecase.CreateSecretKey(state)
		if err != nil {
			log.Println(err.Error())
		} else {
			st = models.StateEnterServiceNameGet
			msg = serviceNameMsgGet
		}
	case models.StateEnterServiceNameGet:
		pass := models.Password{
			ChatID: chatID,
			ServiceName: update.Message.Text,
			Password: "",
		}
		err = h.usecase.CreateServiceName(pass)
		if err != nil {
			log.Print(err.Error())
		}
		st = models.StateEnterPassword
	case models.StateEnterServiceNameSet:
		err := h.db.CreatePassword(models.Password{
			ChatID: chatID,
			ServiceName: update.Message.Text,
		})
		if err != nil {
			log.Print(err.Error())
		}
		st = models.StateEnterPassword	
	case models.StateEnterPassword:
		p, err := h.db.GetPassword(models.Password{ChatID: chatID})
		if err != nil {
			log.Print(err.Error())
		}
		err = h.db.ReplacePassword(models.Password{
			ChatID: chatID,
			ServiceName: p.ServiceName,
			Password: update.Message.Text,
		})
		if err != nil {
			log.Print(err.Error())
		}
		st = models.StateDeafault
	// case models.StateWaitDelKey:
	// case models.StateWaitSetKey:
	default:
	}

	// Edit user's state after each request
	_, err = h.usecase.ReplaceStateByChatAndState(state.ChatID, st)
	if err != nil {
		log.Print(err.Error())
	}
	// return msg
	if err := h.SendMsg(update, msg); err != nil {
		log.Println(err.Error())
	}
}
	
func (h *Telegram) SendMsg(update *bot.Update, text string) error {
	msg := bot.NewMessage(update.Message.Chat.ID, "")
	msg.Text = text
	if _, err := h.bot.Send(msg); err != nil {
		return err
	}
	return nil
}
