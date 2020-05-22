package player

import (
	"fmt"
	"time"

	"github.com/mattshiel/eliza-go/eliza"
	"github.com/nadavash/bot-or-not/src/message"
)

type BotPlayer struct {
	messageChannel chan string
}

func NewBotPlayer() *BotPlayer {
	return &BotPlayer{messageChannel: make(chan string)}
}

func (p *BotPlayer) GetName() string {
	return "eliza"
}

func (p *BotPlayer) GetEmail() string {
	return ""
}

func (p *BotPlayer) SendMessage(msg *message.WrapperMessage) error {
	if msg.MessageType == message.MessageType_CHAT && msg.GetChat().Username != "Room" {
		response := eliza.ReplyTo(msg.GetChat().Message)
		fmt.Printf("Got message from eliza %s\n", response)
		p.messageChannel <- response
	}
	return nil
}

func (p *BotPlayer) ReceiveMessage() (*message.WrapperMessage, error) {
	m := <-p.messageChannel
	time.Sleep(time.Second * 2)
	return message.WrapChatMessage(&message.ChatMessage{
		Username: "eliza",
		Message:  m,
	}), nil
}
