package player

import (
	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
	"github.com/nadavash/bot-or-not/src/netutil"
)

type HumanPlayer struct {
	client *websocket.Conn
}

func NewHumanPlayer(c *websocket.Conn) *HumanPlayer {
	return &HumanPlayer{client: c}
}

func (p *HumanPlayer) SendMessage(msg *message.WrapperMessage) error {
	return netutil.SendProtoMessage(p.client, msg)
}

func (p *HumanPlayer) ReceiveMessage() (*message.WrapperMessage, error) {
	return netutil.ReadProtoMessage(p.client)
}
