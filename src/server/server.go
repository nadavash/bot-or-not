package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var rooms = make([]*Room, 10)
var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Register our new client
	clients[ws] = true

	for _, room := range rooms {
		if room.roomState == RoomStateWaiting {
			room.addClient(ws)
			break
		}
	}
}

func main() {
	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)
	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
