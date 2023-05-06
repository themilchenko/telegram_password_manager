package domain

import "telegram_password_manager/internal/models"

type Usecase interface {
	GetStateByChatID(chatID uint64) (uint64, error)
	CreateStateByChatID(chatID int64) (models.State, error)
}
