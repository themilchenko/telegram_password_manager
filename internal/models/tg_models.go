package models

// There are states of user that helps to route by commands
type StateType int
const (
	StateDeafault StateType = iota

	StateCreateSecretKeySet
	StateCreateSecretKeyGet
	StateCreateSecretKeyDel

	StateWaitSetKey
	StateWaitGetKey
	StateWaitDelKey

	StateIncorrect

	StateEnterServiceNameGet
	StateEnterServiceNameDel

	StateEnterServiceNameSet
	StateEnterPassword
)

type State struct {
	ChatID int64 `msgpack:"chat_id"`
	ChatState StateType `msgpack:"state"`
	SecretKey string `msgpack:"secret_key"`
}

type Password struct {
	ChatID int64 `msgpack:"chat_id"`
	ServiceName string `msgpack:"service_name"`
	Password string `msgpack:"password"`
}
