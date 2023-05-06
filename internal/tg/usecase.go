package tg

import (
	"telegram_password_manager/internal/crypto"
	"telegram_password_manager/internal/domain"
	"telegram_password_manager/internal/models"
)

type Usecase struct {
	repository domain.Repository
}

func NewUsecase(r domain.Repository) Usecase {
	return Usecase{
		repository: r,
	}
}

func (u *Usecase) GetStateByChatID(chatID int64) (models.State, error) {
	s, err := u.repository.GetState(chatID)
	if err != nil {
		return models.State{}, err
	}
	return s, nil
}

func (u *Usecase) CreateStateByChatID(chatID int64) error {
	err := u.repository.CreateState(models.State{
		ChatID: chatID,
		ChatState: models.StateDeafault,
		SecretKey: "",
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) ReplaceStateByChatAndState(chatID int64, state models.StateType) (models.State, error) {
	s, err := u.repository.ReplaceState(models.State{
		ChatID:    chatID,
		ChatState: state,
	})
	if err != nil {
		return models.State{}, err
	}
	return s, nil
}

func (u *Usecase) CreateSecretKey(state models.State) error {
	hashed, err := crypto.HashPassword(state.SecretKey)
	if err != nil {
		return err
	}
	_, err = u.repository.ReplaceState(models.State{
		ChatID: state.ChatID,
		ChatState: state.ChatState,
		SecretKey: hashed,
	})
	if err != nil {
		return err
	}

	return nil
}


func (u *Usecase) ChechSecretKey(state models.State) error {
	st, err := u.repository.GetState(state.ChatID)
	if err != nil {
		return err
	}

	crypto.CheckHashPassword(state.SecretKey, st.SecretKey)
	return nil
}
