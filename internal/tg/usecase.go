package tg

import (
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

func (u *Usecase) CreateStateByChatID(chatID int64) (models.State, error) {
	s, err := u.repository.CreateState(chatID)
	if err != nil {
		return models.State{}, err
	}
	return s, nil
}
