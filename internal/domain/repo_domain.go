package domain

import "telegram_password_manager/internal/models"

type Repository interface {
	GetState(chatID int64) (models.State, error)
	CreateState(chatID int64) (models.State, error)
	Close() error
}
