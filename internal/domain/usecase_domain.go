package domain

import "telegram_password_manager/internal/models"

type Usecase interface {
	GetStateByChatID(chatID int64) (models.State, error)
	CreateStateByChatID(chatID int64) error
	ReplaceStateByChatAndState(chatID int64, state models.StateType) (models.State, error)
}
