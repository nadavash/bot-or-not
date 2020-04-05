package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
)

type clientMessagePair struct {
	message message.Message
	client  *websocket.Conn
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan clientMessagePair)

var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	// Register our new client
	clients[ws] = true

	for {
		var msg message.Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- clientMessagePair{client: ws, message: msg}
		fmt.Println("Message received:", msg)
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		clientMessage := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			if client == clientMessage.client {
				continue
			}
			err := client.WriteJSON(clientMessage.message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)
	// Start listening for incoming chat messages
	go handleMessages()
	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
