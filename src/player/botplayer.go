package player

import (
	"github.com/nadavash/bot-or-not/src/message"
)

type BotPlayer struct {
}

func (p *BotPlayer) SendMessage(msg *message.WrapperMessage) error {
	return nil
}

func (p *BotPlayer) ReceiveMessage() (*message.WrapperMessage, error) {
	return nil, nil
}
