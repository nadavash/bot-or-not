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
		var m message.MessageBase
		err := conn.ReadJSON(&m)
		if err != nil {
			log.Printf("error: %v", err)
		}

		messageBody := m.MessageBody.(map[string]interface{})
		switch m.MessageType {
		case message.MessageTypeServerConnectionSuccess:
			fmt.Println(messageBody["WelcomeMessage"])
		case message.MessageTypeRoomConnectionSuccess:
			fmt.Printf("Successfully connected to room %g\n", messageBody["roomId"])
		case message.MessageTypeChat:
			// Return to beginning.
			fmt.Print(cursor.ClearEntireLine())
			fmt.Printf("\r%s: %s\n> ", messageBody["username"], messageBody["message"])
		case message.MessageTypeGameOver:
			fmt.Println("Game over. Disconnecting from server.")
			conn.Close()
			return
		}
	}
}

func handleOutgoingMessages(scanner *bufio.Scanner, name string, conn *websocket.Conn) {
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			log.Printf("Scanner.Scan() returned false!")
		}
		err := conn.WriteJSON(
			&message.ChatMessage{
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
	conn, _, err := dialer.Dial("ws://localhost:8000/ws", requestHeader)
	if err != nil {
		log.Fatal("Error occurred during Dialer.Dial(): ", err)
	}
	go handleIncomingMessages(conn)
	handleOutgoingMessages(scanner, name, conn)
}
