package player

import (
	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
	"github.com/nadavash/bot-or-not/src/netutil"
)

type HumanPlayer struct {
	name   string
	email  string
	client *websocket.Conn
}

func NewHumanPlayer(name string, email string, c *websocket.Conn) *HumanPlayer {
	return &HumanPlayer{name: name, email: email, client: c}
}

func (p *HumanPlayer) GetName() string {
	return p.name
}

func (p *HumanPlayer) GetEmail() string {
	return p.email
}

func (p *HumanPlayer) SendMessage(msg *message.WrapperMessage) error {
	return netutil.SendProtoMessage(p.client, msg)
}

func (p *HumanPlayer) ReceiveMessage() (*message.WrapperMessage, error) {
	return netutil.ReadProtoMessage(p.client)
}
