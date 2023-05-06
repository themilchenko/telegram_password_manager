package models

// There are states of user that helps to route by commands
type StateType int
const (
	StateDeafault StateType = iota
	StateWaitSetKey
	StateWaitGetKey
	StateWaitDelKey
	StateIncorrect
	StateRightGet
	StateRightAdd
	StateRightDel
)

type TgMessage struct {
	MessageID int
	UserID uint64 
	ChatID uint64
	Text string
}

type State struct {
	ChatID int64
	ChatState StateType
}
