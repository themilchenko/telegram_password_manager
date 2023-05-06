package tg

import (
	"telegram_password_manager/internal/models"

	"github.com/tarantool/go-tarantool"
)

type DB struct {
	client *tarantool.Connection
}

func NewDB(addr string, opts *tarantool.Opts) (*DB, error) {
	client, err := tarantool.Connect(addr, *opts)
	if err != nil {
		return nil, err
	}
	return &DB{client: client}, nil
}

func (db *DB) Close() error {
	err := db.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetState(chatID int64) (models.State, error) 

func (db *DB) CreateState(chatID int64) (models.State, error)

