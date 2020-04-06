package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	cursor "github.com/ahmetb/go-cursor"
	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
)

func handleIncomingMessages(conn *websocket.Conn) {
	for {
		var message message.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("error: %v", err)
		}
		// Return to beginning.
		fmt.Print(cursor.ClearEntireLine())
		fmt.Printf("\r%s: %s\n> ", message.Username, message.Message)
	}
}

func handleOutgoingMessages(scanner *bufio.Scanner, name string, conn *websocket.Conn) {
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
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("What's your name?\nname: ")
	scanner.Scan()
	name := scanner.Text()

	requestHeader := http.Header{}
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://0b5eb1db.ngrok.io/ws", requestHeader)
	if err != nil {
		log.Fatal("Error occurred during Dialer.Dial(): ", err)
	}
	go handleIncomingMessages(conn)
	handleOutgoingMessages(scanner, name, conn)
}
