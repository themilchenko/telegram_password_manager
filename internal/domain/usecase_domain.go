package domain

import "telegram_password_manager/internal/models"

type Usecase interface {
	GetStateByChatID(chatID int64) (models.State, error)
	CreateStateByChatID(chatID int64) error
	ReplaceStateByChatAndState(chatID int64, state models.StateType) error
	CheckExistingKey(chatID int64) error
	CreateSecretKey(state models.State) error
	CheckSecretKey(chatID int64, key string) error
	GetSimplePassword(pass models.Password) (models.Password, error)
	CreateServiceName(pass models.Password) error
	AddRequestServiceName(state models.State) error
	AddPassword(pass models.Password, key string) error
	GetPassword(pass models.Password, key string) (models.Password, error)
	DeleteService(pass models.Password) error
	GetServicesByChatID(chatID int64) ([]models.Password, error)
	CheckServiceExists(chatID int64, serviceName string) error
}
