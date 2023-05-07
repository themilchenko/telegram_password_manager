package domain

import "telegram_password_manager/internal/models"

type Repository interface {
	GetState(chatID int64) (models.State, error)
	CreateState(state models.State) error
	ReplaceState(state models.State) (models.State, error)

	CreatePassword(pass models.Password) error 
	GetPassword(pass models.Password) (models.Password, error) 
	ReplacePassword(pass models.Password) error 
	DeletePassword(pass models.Password) error 
	Close() error
}
