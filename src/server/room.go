package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
)

type RoomState int

const (
	RoomStateWaiting    = 0
	RoomStateInProgress = 1
	RoomStateFinished   = 2
	gameTime            = 2
	roomLimit           = 2
)

type Room struct {
	roomState        RoomState
	clients          []*websocket.Conn
	broadcastChannel chan clientMessagePair
}

type clientMessagePair struct {
	message message.Message
	client  *websocket.Conn
}

func NewRoom() *Room {
	r := new(Room)
	r.clients = make([]*websocket.Conn, 0, roomLimit)
	r.roomState = RoomStateWaiting
	r.broadcastChannel = make(chan clientMessagePair)
	return r
}

func (r *Room) addClient(client *websocket.Conn) error {
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
			go r.acceptIncomingMessages(client)
		}
	}
	return nil
}

func (r *Room) test() {
	fmt.Println("in test")
}

func (r *Room) acceptIncomingMessages(client *websocket.Conn) {
	for r.roomState == RoomStateInProgress {
		var msg message.Message
		// Read in a new message as JSON and map it to a Message object
		err := client.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			r.removeClient(client)
			client.Close()
			break
		}
		// Send the newly received message to the broadcast channel
		r.broadcastChannel <- clientMessagePair{client: client, message: msg}
		fmt.Println("Message received:", msg)
	}
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

	for _, client := range r.clients {
		client.Close()
	}
}

func (r *Room) handleGameLogic() {
	fmt.Println("4 client")
	for minutesLeft := gameTime; minutesLeft > 0; minutesLeft-- {
		r.sendRoomMessage(
			fmt.Sprintf("%d minutes left in the game", minutesLeft))
		time.Sleep(time.Minute)
	}

	r.sendRoomMessage("Times up, game is over!")
	r.roomState = RoomStateFinished
}

func (r *Room) sendRoomMessage(s string) {
	go func() {
		r.broadcastChannel <- clientMessagePair{
			message.Message{
				Username: "Room",
				Message:  s,
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
