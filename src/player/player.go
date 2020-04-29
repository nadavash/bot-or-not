package player

import "github.com/nadavash/bot-or-not/src/message"

type Player interface {
	SendMessage(msg *message.WrapperMessage) error
	// conn.ReadMessage
	ReceiveMessage() (*message.WrapperMessage, error)
}
