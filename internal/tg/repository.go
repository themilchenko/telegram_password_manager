package tg

import (
	"errors"
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

func (db *DB) GetState(chatID int64) (models.State, error) {
	var res []models.State
	err := db.client.SelectTyped("users", "primary", 0, 1, tarantool.IterEq, tarantool.IntKey{int(chatID)}, &res)
	if err != nil {
		return models.State{}, nil
	}
	if len(res) == 0 {
		return models.State{}, errors.New("ErrNotFound")
	}
	if len(res) > 1 {
		return models.State{}, errors.New("ErrManyRows")
	}

	return models.State{
		ChatID: res[0].ChatID,
		ChatState: res[0].ChatState,
	}, nil
}

func (db *DB) CreateState(state models.State) (error) {
	_, err := db.client.Insert("users", []interface{}{
		state.ChatID, 
		state.ChatState,
		state.SecretKey,
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ReplaceState(state models.State) (models.State, error) {
	var res models.State
	_, err := db.client.Replace("users", []interface{}{
		state.ChatID, 
		state.ChatState,
		state.SecretKey,
	})
	if err != nil {
		return models.State{}, err
	}
	return res, nil
}

func (db *DB) CreatePassword(pass models.Password) error {
	_, err := db.client.Insert("passwords", []interface{}{
		pass.ChatID,
		pass.ServiceName,
		pass.Password,
		})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetPassword(pass models.Password) (models.Password, error) {
	var res models.Password
	err := db.client.SelectTyped("passwords", "primary", 0, 1, tarantool.IterEq, []interface{}{pass.ChatID, pass.ServiceName}, &res)
	if err != nil {
		return models.Password{}, nil
	}
	return res, nil
}

func (db *DB) ReplacePassword(pass models.Password) error {
	_, err := db.client.Replace("passwords", []interface{}{
		pass.ChatID, 
		pass.ServiceName,
		pass.Password,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeletePassword(pass models.Password) error {
	_, err := db.client.Delete("passwords", "primary", []interface{}{pass.ChatID, pass.ServiceName})
	if err != nil {
		return err
	}
	return nil
}

