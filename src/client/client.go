package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"../message"
)

func handleIncomingMessages(conn *websocket.Conn) {
	for {
		var message message.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("error: %v", err)
		}
		// Return to beginning.
		fmt.Printf("\r")
		fmt.Printf("%s: %s\n", message.Username, message.Message)
	}
}

func handleOutgoingMessages(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("What's your name?\nname: ")
	scanner.Scan()
	name := scanner.Text()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			log.Printf("Scanner.Scan() returned false!")
		}
		err := conn.WriteJSON(
			&message.Message{
				Email:    "example@test.com",
				Username: name,
				Message:  scanner.Text(),
			})
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}

func main() {
	requestHeader := http.Header{}
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://localhost:8000/ws", requestHeader)
	if err != nil {
		log.Fatal("Error occurred during Dialer.Dial(): ", err)
	}
	go handleIncomingMessages(conn)
	handleOutgoingMessages(conn)
}
