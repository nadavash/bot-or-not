package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
	"github.com/nadavash/bot-or-not/src/netutil"
)

var rooms = make([]*Room, 10)
var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	msg := message.WrapServerConnectionSuccessMessage(
		&message.ServerConnectionSuccessMessage{
			WelcomeMessage: "You're connected to the Bot or Not server!",
		},
	)
	if netutil.SendProtoMessage(ws, msg) != nil {
		return
	}
	AssignRoom(ws)
}

func AssignRoom(client *websocket.Conn) {
	for _, room := range rooms {
		if room.roomState == RoomStateWaiting {
			room.AddClient(client)
			break
		}
	}
}

func main() {
	fmt.Println("V5")
	for i := 0; i < cap(rooms); i++ {
		rooms[i] = NewRoom(AssignRoom)
	}

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)
	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
