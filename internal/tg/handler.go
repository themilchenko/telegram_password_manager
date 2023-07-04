package tg

import (
	"fmt"
	"log"
	"time"

	bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"telegram_password_manager/internal/domain"
	"telegram_password_manager/internal/models"
)

// Messages for responses
const (
	createPrivateKeyMsg = "Create your password to get access to the password manager. Remember that this password (secret key) will be a key for all your passwords that you will communicate with. Don't loose it."
	privateKeyMsg       = "Enter your secret key"
	incorrectSecretKey  = "Your secret key is incorrect"
	serviceNameMsgSet   = "Enter the name of service you want to create password"
	serviceNameMsgDel   = "Enter the name of service you want to delete password"
	serviceNameMsgGet   = "Enter the name of service you want to get password"
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
	isPass := false

	switch state.ChatState {
	case models.StateDeafault:
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
				st = models.StateEnterServiceNameSet
				msg = "Enter the name of new service"
			}
		case "get":
			res, err := h.usecase.GetServicesByChatID(chatID)
			if err != nil {
				switch err.Error() {
				case "ErrNotFound":
					st = models.StateDeafault
					msg = "You haven't created service"
				default:
					log.Print(err.Error())
				}
			} else {
				st = models.StateEnterServiceNameGet
				msg = "Choose service. Enter its name\n"
				for num, item := range res {
					msg += fmt.Sprintf("%d) %s\n", num+1, item.ServiceName)
				}
			}
		case "del":
			res, err := h.usecase.GetServicesByChatID(chatID)
			if err != nil {
				switch err.Error() {
				case "ErrNotFound":
					st = models.StateDeafault
					msg = "You haven't created service"
				default:
					log.Print(err.Error())
				}
			} else {
				st = models.StateEnterServiceNameDel
				msg = "Choose service. Enter its name\n"
				for num, item := range res {
					msg += fmt.Sprintf("%d) %s\n", num+1, item.ServiceName)
				}
			}
		case "help":
			msg = "Enter /set to create password\nEnter /get to get password\nEnter /del to delete service"
		default:
			msg = "I don't know that command"
		}
	case models.StateWaitGetKey:
		err = h.usecase.CheckSecretKey(chatID, update.Message.Text)
		if err != nil {
			switch err.Error() {
			case "ErrIncorrectSecretKey":
				st = models.StateDeafault
				msg = incorrectSecretKey

			default:
				log.Print(err.Error())
			}
		} else {
			go h.deleteMessage(chatID, int(update.Message.MessageID), 2*time.Second)
			p, err := h.usecase.GetPassword(models.Password{
				ChatID:      chatID,
				ServiceName: state.RequestService,
			}, update.Message.Text)
			if err != nil {
				log.Print(err.Error())
			}
			st = models.StateDeafault
			msg = p.Password
			isPass = true
		}
	case models.StateWaitDelKey:
		err = h.usecase.CheckSecretKey(chatID, update.Message.Text)
		if err != nil {
			switch err.Error() {
			case "ErrIncorrectSecretKey":
				msg = "Incorrect secret key. Try again"
				st = models.StateDeafault
			}
		} else {
			go h.deleteMessage(chatID, int(update.Message.MessageID), 2*time.Second)
			err = h.usecase.DeleteService(models.Password{
				ChatID:      chatID,
				ServiceName: state.RequestService,
			})
			if err != nil {
				log.Print(err.Error())
			}
			msg = fmt.Sprintf("Service %s successfully deleted", state.RequestService)
			st = models.StateDeafault
		}
	case models.StateWaitSetKey:
		err = h.usecase.CheckSecretKey(chatID, update.Message.Text)
		if err != nil {
			switch err.Error() {
			case "ErrIncorrectSecretKey":
				msg = "Incorrect secret key. Try again"
				st = models.StateDeafault
				err = h.usecase.DeleteService(models.Password{
					ChatID:      chatID,
					ServiceName: state.RequestService,
				})
			}
		} else {
			go h.deleteMessage(chatID, int(update.Message.MessageID), 2*time.Second)
			p, _ := h.usecase.GetSimplePassword(models.Password{
				ChatID:      chatID,
				ServiceName: state.RequestService,
			})
			err = h.usecase.AddPassword(p, update.Message.Text)
			if err != nil {
				log.Print(err.Error())
			}
			msg = "Created"
			st = models.StateDeafault
		}
	case models.StateCreateSecretKeySet:
		state.SecretKey = update.Message.Text
		err = h.usecase.CreateSecretKey(state)
		if err != nil {
			log.Println(err.Error())
		} else {
			go h.deleteMessage(chatID, int(update.Message.MessageID), 2*time.Second)
			st = models.StateDeafault
			msg = "Secret key was created. Now you can `/set` password"
		}
	case models.StateEnterServiceNameGet:
		err := h.usecase.CheckServiceExists(chatID, update.Message.Text)
		if err != nil {
			if err.Error() == "ErrNotFound" {
				msg = "There is no service with its name"
				st = models.StateDeafault
			} else {
				log.Print(err.Error())
			}
		} else {
			h.usecase.AddRequestServiceName(models.State{
				ChatID:         chatID,
				ChatState:      state.ChatState,
				RequestService: update.Message.Text,
			})
			if err != nil {
				log.Print(err.Error())
			}
			msg = "Enter your secret key"
			st = models.StateWaitGetKey
		}
	case models.StateEnterServiceNameSet:
		h.usecase.CheckServiceExists(chatID, update.Message.Text)
		if err != nil {
			msg = "You already have service with this name"
			st = models.StateDeafault
		} else {
			err = h.usecase.AddRequestServiceName(models.State{
				ChatID:         chatID,
				ChatState:      state.ChatState,
				RequestService: update.Message.Text,
			})
			if err != nil {
				log.Print(err.Error())
			}
			st = models.StateEnterPassword
			msg = "Enter password 🔐"
		}
	case models.StateEnterServiceNameDel:
		err := h.usecase.CheckServiceExists(chatID, update.Message.Text)
		if err != nil {
			if err.Error() == "ErrNotFound" {
				msg = "There is no service with its name"
				st = models.StateDeafault
			} else {
				log.Print(err.Error())
			}
		} else {
			h.usecase.AddRequestServiceName(models.State{
				ChatID:         chatID,
				ChatState:      state.ChatState,
				RequestService: update.Message.Text,
			})
			if err != nil {
				log.Print(err.Error())
			}
			msg = "Enter your secret key 🔑"
			st = models.StateWaitDelKey
		}
	case models.StateEnterPassword:
		go h.deleteMessage(chatID, int(update.Message.MessageID), 2*time.Second)
		h.usecase.CreateServiceName(models.Password{
			ChatID:      state.ChatID,
			ServiceName: state.RequestService,
			Password:    update.Message.Text,
		})
		if err != nil {
			log.Print(err.Error())
		}
		st = models.StateWaitSetKey
		msg = "Enter secret key 🔑"
	default:
	}

	err = h.usecase.ReplaceStateByChatAndState(chatID, st)
	if err != nil {
		log.Print(err.Error())
	}

	if isPass {
		msg += "\nThis message will be deleted after 30 seconds!"
	}

	msgID, err := h.sendMsg(update, msg)
	if err != nil {
		log.Println(err.Error())
	}

	if isPass {
		go h.deleteMessage(chatID, int(msgID), 30*time.Second)
	}
}

func (h *Telegram) sendMsg(update *bot.Update, text string) (int, error) {
	msg := bot.NewMessage(update.Message.Chat.ID, "")
	msg.Text = text

	status, err := h.bot.Send(msg)
	if err != nil {
		return 0, err
	}

	return status.MessageID, nil
}

func (h *Telegram) deleteMessage(chatID int64, messageID int, delay time.Duration) {
	time.Sleep(delay)
	msg := bot.NewDeleteMessage(chatID, messageID)

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
