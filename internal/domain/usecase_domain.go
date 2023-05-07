package domain

import "telegram_password_manager/internal/models"

type Usecase interface {
	GetStateByChatID(chatID int64) (models.State, error)
	CreateStateByChatID(chatID int64) error
	ReplaceStateByChatAndState(chatID int64, state models.StateType) (models.State, error)
	CheckExistingKey(chatID int64) error
	CreateSecretKey(state models.State) error
	ChechSecretKey(state models.State) error
	
	CreateServiceName(pass models.Password) error
	AddPassword(pass models.Password) error
	GetPassword(pass models.Password) (models.Password, error) 
	DeleteService(pass models.Password) error 
}
