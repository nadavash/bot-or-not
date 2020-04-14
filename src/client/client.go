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

type GameState int

const (
	GameStatePlaying  = 0
	GameStateDeciding = 1
	GameStateWaiting  = 2
	GameStateGameOver = 3
)

var state GameState = GameStateWaiting

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
			state = GameStatePlaying
			fmt.Printf("Successfully connected to room %g\n", messageBody["roomId"])
		case message.MessageTypeChat:
			// Return to beginning.
			fmt.Print(cursor.ClearEntireLine())
			fmt.Printf("\r%s: %s\n> ", messageBody["username"], messageBody["message"])
		case message.MessageTypeGameOver:
			state = GameStateDeciding
			fmt.Println("Game over. Disconnecting from server.")
			fmt.Println("Bot Or Not?.")
			return
		}
	}
}

func handleOutgoingMessages(scanner *bufio.Scanner, name string, conn *websocket.Conn) {
	for state != GameStateGameOver {
		fmt.Print("> ")
		if !scanner.Scan() {
			log.Printf("Scanner.Scan() returned false!")
		}
		var err error = nil
		switch state {
		case GameStatePlaying:
			err = conn.WriteJSON(
				&message.ChatMessage{
					Email:    "example@test.com",
					Username: name,
					Message:  scanner.Text(),
				})
		case GameStateDeciding:
			err = conn.WriteJSON(
				&message.BotOrNotAnswerMessage{
					ArePlayersBotsAnswer: scanner.Text()[0] == 'b',
				})
			if err != nil {
				log.Printf("error: %v", err)
			}
			fmt.Println("Do you want to play again? (y/n)")
			if !scanner.Scan() {
				log.Printf("Scanner.Scan() returned false!")
			}
			err = conn.WriteJSON(
				&message.PlayAgainMessage{
					PlayAgain: scanner.Text()[0] == 'y',
				})
			if err != nil {
				log.Printf("error: %v", err)
			}
			if scanner.Text()[0] == 'y'{
				state = GameStateWaiting
			} else {
				state = GameStateGameOver
			}
		}
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
