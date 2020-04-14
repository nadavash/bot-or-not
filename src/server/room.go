package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
	"github.com/mitchellh/mapstructure"
)

type RoomState int

const (
	RoomStateWaiting    = 0
	RoomStateInProgress = 1
	RoomStateFinished   = 2
	gameTimeSeconds     = 5
	roomLimit           = 1
)

type Room struct {
	roomId           uint32
	roomState        RoomState
	clients          []*websocket.Conn
	broadcastChannel chan clientMessagePair
	finishedCallback func(*websocket.Conn)
}

type clientMessagePair struct {
	message message.MessageBase
	client  *websocket.Conn
}

func NewRoom(onClientFinished func(*websocket.Conn)) *Room {
	r := new(Room)
	r.roomId = rand.Uint32() % 10000
	r.clients = make([]*websocket.Conn, 0, roomLimit)
	r.roomState = RoomStateWaiting
	r.broadcastChannel = make(chan clientMessagePair)
	r.finishedCallback = onClientFinished
	return r
}

func (r *Room) GetRoomId() uint32 {
	return r.roomId
}

func (r *Room) AddClient(client *websocket.Conn) error {
	fmt.Println("addingClient")
	fmt.Println(len(r.clients))
	if r.roomState != RoomStateWaiting {
		return errors.New("Cannot add clients to a Room that's in progress or finished")
	}
	r.clients = append(r.clients, client)
	if len(r.clients) == roomLimit {
		fmt.Println("4 clients added")
		go r.test()
		go r.broadcastMessages()
		r.roomState = RoomStateInProgress
		r.sendRoomMessage("Room full! Game starting now.")
		go r.handleGameLogic()
		for _, client := range r.clients {
			fmt.Println("accepting from client")
			go r.acceptIncomingChatMessages(client)
		}
	}
	client.WriteJSON(
		message.MessageBase{
			MessageType: message.MessageTypeRoomConnectionSuccess,
			MessageBody: message.RoomConnectionSuccessMessage{
				RoomId: r.roomId,
			},
		},
	)
	return nil
}

func (r *Room) test() {
	fmt.Println("in test")
}

func (r *Room) acceptIncomingChatMessages(client *websocket.Conn) {
	for r.roomState == RoomStateInProgress {
		var msg interface{}
		// Read in a new message as JSON and map it to a Message object
		err := client.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			r.removeClient(client)
			client.Close()
			break
		}
		if r.roomState == RoomStateInProgress {
			var chatMessage message.ChatMessage
			mapstructure.Decode(msg, &chatMessage)
			r.broadcastChannel <- clientMessagePair{
				client: client,
				message: message.MessageBase{
					MessageType: message.MessageTypeChat,
					MessageBody: chatMessage,
				},
			}
		} else if r.roomState == RoomStateFinished {
			var decisionMessage message.BotOrNotAnswerMessage
			mapstructure.Decode(msg, &decisionMessage)

			fmt.Println("accepting decisions from clients")
			fmt.Println("Decision received:", decisionMessage)
			log.Printf("decision: %v", decisionMessage.ArePlayersBotsAnswer)
			r.roomState = RoomStateFinished
			r.acceptPlayAgain(client)
		}
	}
}

func (r *Room) acceptPlayAgain(client *websocket.Conn) {
	fmt.Println("about to accept play again")
	var msg message.PlayAgainMessage
	err := client.ReadJSON(&msg)
	fmt.Println("reading play again")
	if err != nil {
		log.Printf("error: %v", err)
		r.removeClient(client)
		client.Close()
	}
	r.finishedCallback(client)
}

func (r *Room) broadcastMessages() {
	fmt.Println("starting broadcast messages")
	for r.roomState == RoomStateInProgress || len(r.broadcastChannel) > 0 {
		fmt.Println("in loop")
		clientMessage := <-r.broadcastChannel
		fmt.Println(clientMessage.message)
		for _, client := range r.clients {
			if client == clientMessage.client {
				continue
			}
			err := client.WriteJSON(clientMessage.message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
			}
		}
	}

	//for _, client := range r.clients {
	//	client.Close()
	//}
}

func (r *Room) handleGameLogic() {
	fmt.Println("4 client")
	for secondsLeft := gameTimeSeconds; secondsLeft > 0; secondsLeft-- {
		r.sendRoomMessage(
			fmt.Sprintf("%d seconds left in the game", secondsLeft))
		time.Sleep(time.Second)
	}

	r.sendRoomMessage("Times up, game is over!")
	r.roomState = RoomStateFinished
	r.broadcastChannel <- clientMessagePair{
		message.MessageBase{
			MessageType: message.MessageTypeGameOver,
			MessageBody: message.GameOverMessage{},
		},
		nil,
	}
}

func (r *Room) sendRoomMessage(s string) {
	go func() {
		r.broadcastChannel <- clientMessagePair{
			message.MessageBase{
				MessageType: message.MessageTypeChat,
				MessageBody: message.ChatMessage{
					Username: "Room",
					Message:  s,
				},
			},
			nil,
		}
	}()
}

func (r *Room) removeClient(client *websocket.Conn) {
	for i, clientPointer := range r.clients {
		if client == clientPointer {
			r.clients = append(r.clients[i:], r.clients[i+1:]...)
			break
		}
	}
}
