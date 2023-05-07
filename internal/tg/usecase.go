package tg

import (
	"errors"
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

func (u *Usecase) CheckExistingKey(chatID int64) error {
	st, err := u.repository.GetState(chatID)
	if err != nil {
		return err
	}

	if len(st.SecretKey) == 0 {
		return errors.New("ErrSecKeyCreate")
	}
	return nil
}

func (u *Usecase) ChechSecretKey(state models.State) error {
	st, err := u.repository.GetState(state.ChatID)
	if err != nil {
		return err
	}

	if !crypto.CheckHashPassword(state.SecretKey, st.SecretKey) {
		return errors.New("ErrIncorrectSecretKey")
	}
	return nil
}

func (u *Usecase) CreateServiceName(pass models.Password) error {
	err := u.repository.CreatePassword(pass)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) AddPassword(pass models.Password) error {
	err := u.repository.ReplacePassword(pass)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) GetPassword(pass models.Password) (models.Password, error) {
	res, err := u.repository.GetPassword(pass)
	if err != nil {
		return models.Password{}, err
	}
	return res, nil
}

func (u *Usecase) DeleteService(pass models.Password) error {
	err := u.repository.DeletePassword(pass)
	if err != nil {
		return err
	}
	return nil
}

