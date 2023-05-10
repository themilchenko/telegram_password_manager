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
		ChatID:    chatID,
		ChatState: models.StateDeafault,
		SecretKey: "",
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) ReplaceStateByChatAndState(chatID int64, state models.StateType) error {
	err := u.repository.ReplaceState(models.State{
		ChatID:         chatID,
		ChatState:      state,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) AddRequestServiceName(state models.State) error {
	err := u.repository.ReplaceState(state)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) CreateSecretKey(state models.State) error {
	hashed, err := crypto.HashPassword(state.SecretKey)
	if err != nil {
		return err
	}
	err = u.repository.ReplaceState(models.State{
		ChatID:    state.ChatID,
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

func (u *Usecase) CheckSecretKey(chatID int64, key string) error {
	st, err := u.repository.GetState(chatID)
	if err != nil {
		return err
	}
	if crypto.CheckHashPassword(key, st.SecretKey) != true {
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

func (u *Usecase) GetSimplePassword(pass models.Password) (models.Password, error) {
	res, err := u.repository.GetPassword(pass)
	if err != nil {
		return models.Password{}, err
	}
	return res, nil
}

func (u *Usecase) AddPassword(pass models.Password, key string) error {
	encrypted, err := crypto.Encrypt([]byte(key), pass.Password)
	if err != nil {
		return err
	}
	pass.Password = encrypted

	err = u.repository.ReplacePassword(pass)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) GetPassword(pass models.Password, key string) (models.Password, error) {
	res, err := u.repository.GetPassword(pass)
	if err != nil {
		return models.Password{}, err
	}
	decrypted, err := crypto.Decrypt([]byte(key), res.Password)
	if err != nil {
		return models.Password{}, err
	}
	res.Password = decrypted

	return res, nil
}

func (u *Usecase) DeleteService(pass models.Password) error {
	err := u.repository.DeletePassword(pass)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) GetServicesByChatID(chatID int64) ([]models.Password, error) {
	res, err := u.repository.GetServices(chatID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Usecase) CheckServiceExists(chatID int64, serviceName string) error {
	_, err := u.repository.GetPassword(models.Password{
		ChatID: chatID,
		ServiceName: serviceName,
	})
	if err != nil {
		return err
	}
	return nil
}
