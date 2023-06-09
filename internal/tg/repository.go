package tg

import (
	"errors"

	"github.com/tarantool/go-tarantool"

	"telegram_password_manager/internal/models"
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
	err := db.client.SelectTyped(
		"users",
		"primary",
		0,
		1,
		tarantool.IterEq,
		tarantool.IntKey{int(chatID)},
		&res,
	)
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
		ChatID:         res[0].ChatID,
		ChatState:      res[0].ChatState,
		SecretKey:      res[0].SecretKey,
		RequestService: res[0].RequestService,
	}, nil
}

func (db *DB) CreateState(state models.State) error {
	_, err := db.client.Insert("users", []interface{}{
		state.ChatID,
		state.ChatState,
		state.SecretKey,
		state.RequestService,
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ReplaceState(state models.State) error {
	st, err := db.GetState(state.ChatID)
	if err != nil {
		return err
	}

	if len(state.RequestService) != 0 {
		_, err = db.client.Replace("users", []interface{}{
			state.ChatID,
			state.ChatState,
			st.SecretKey,
			state.RequestService,
		})
	} else if len(state.SecretKey) == 0 {
		_, err = db.client.Replace("users", []interface{}{
			state.ChatID,
			state.ChatState,
			st.SecretKey,
			st.RequestService,
		})
	} else {
		_, err = db.client.Replace("users", []interface{}{
			state.ChatID,
			state.ChatState,
			state.SecretKey,
			state.RequestService,
		})
	}
	if err != nil {
		return err
	}
	return nil
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
	var res []models.Password
	err := db.client.SelectTyped(
		"passwords",
		"primary",
		0,
		1,
		tarantool.IterEq,
		[]interface{}{pass.ChatID, pass.ServiceName},
		&res,
	)
	if err != nil {
		return models.Password{}, nil
	}
	if len(res) == 0 {
		return models.Password{}, errors.New("ErrNotFound")
	}
	return res[0], nil
}

func (db *DB) ReplacePassword(pass models.Password) error {
	var err error
	p, err := db.GetPassword(pass)
	if len(pass.Password) == 0 {
		_, err = db.client.Replace("passwords", []interface{}{
			p.ChatID,
			p.ServiceName,
			p.Password,
		})
	} else {
		_, err = db.client.Replace("passwords", []interface{}{
			pass.ChatID,
			pass.ServiceName,
			pass.Password,
		})
	}
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

func (db *DB) GetServices(chatID int64) ([]models.Password, error) {
	var res []models.Password
	err := db.client.SelectTyped(
		"passwords",
		"primary",
		0,
		50,
		tarantool.IterEq,
		tarantool.IntKey{int(chatID)},
		&res,
	)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("ErrNotFound")
	}
	return res, nil
}
