package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
)

var rooms = make([]*Room, 10)
var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	ws.WriteJSON(
		message.MessageBase{
			MessageType: message.MessageTypeServerConnectionSuccess,
			MessageBody: message.ServerConnectionSuccessMessage{
				WelcomeMessage: "You're connected to the Bot or Not server!",
			},
		},
	)

	for _, room := range rooms {
		if room.roomState == RoomStateWaiting {
			room.AddClient(ws)
			break
		}
	}
}

func main() {
	fmt.Println("V5")
	for i := 0; i < cap(rooms); i++ {
		rooms[i] = NewRoom()
	}

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)
	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
