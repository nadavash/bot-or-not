package player

import "github.com/nadavash/bot-or-not/src/message"

type Player interface {
	GetName() string
	GetEmail() string
	SendMessage(msg *message.WrapperMessage) error
	ReceiveMessage() (*message.WrapperMessage, error)
}
