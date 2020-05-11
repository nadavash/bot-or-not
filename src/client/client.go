package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahmetb/go-cursor"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
	"github.com/nadavash/bot-or-not/src/netutil"
)

type GameState int

const (
	GameStatePlaying  = 0
	GameStateDeciding = 1
	GameStateWaiting  = 2
	GameStateGameOver = 3
)

var state GameState = GameStateWaiting
var arePlayersBotsAnswer bool = false

func handleIncomingMessages(conn *websocket.Conn) {
	for {
		wrapperMsg := &message.WrapperMessage{}
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
		}

		if err := proto.Unmarshal(bytes, wrapperMsg); err != nil {
			log.Printf("Unrmarshalling error: %v", err)
		}

		switch wrapperMsg.MessageType {
		case message.MessageType_SERVER_CONNECTION_SUCCESS:
			fmt.Println(wrapperMsg.GetServerSuccess().GetWelcomeMessage())
		case message.MessageType_ROOM_CONNECTION_SUCCESS:
			state = GameStatePlaying
			fmt.Printf(
				"Successfully connected to room %d\n",
				wrapperMsg.GetRoomSuccess().GetRoomId(),
			)
		case message.MessageType_CHAT:
			chatMsg := wrapperMsg.GetChat()
			// Return to beginning.
			fmt.Print(cursor.ClearEntireLine())
			fmt.Printf("\r%s: %s\n> ", chatMsg.GetUsername(), chatMsg.GetMessage())
		case message.MessageType_GAME_OVER:
			state = GameStateDeciding
			fmt.Println("Bot Or Not?.")
		case message.MessageType_ANSWER_CORRECT:
			state = GameStateGameOver
			if wrapperMsg.GetAnswerCorrect().GetIsCorrectAnswer() {
				fmt.Println("You answered correctly!")
			} else {
				fmt.Println("You answered incorrectly!")
			}
		}
	}
}

func handleOutgoingMessages(scanner *bufio.Scanner, name string, conn *websocket.Conn) {
	for state != GameStateGameOver {
		fmt.Print("> ")
		if !scanner.Scan() {
			log.Printf("Scanner.Scan() returned false!")
		}
		switch state {
		case GameStatePlaying:
			netutil.SendProtoMessage(
				conn,
				message.WrapChatMessage(
					&message.ChatMessage{
						Email:    "example@test.com",
						Username: name,
						Message:  scanner.Text(),
					},
				),
			)
		case GameStateDeciding:
			arePlayersBotsAnswer = scanner.Text()[0] == 'b'
			netutil.SendProtoMessage(
				conn,
				message.WrapBotOrNotMessage(
					&message.BotOrNotMessage{
						ArePlayersBots: arePlayersBotsAnswer,
					},
				),
			)

			fmt.Println("Do you want to play again? (y/n)")
			if !scanner.Scan() {
				log.Printf("Scanner.Scan() returned false!")
			}

			netutil.SendProtoMessage(
				conn,
				message.WrapPlayAgainMessage(
					&message.PlayAgainMessage{
						PlayAgain: scanner.Text()[0] == 'y',
					},
				),
			)

			if scanner.Text()[0] == 'y' {
				state = GameStateWaiting
			} else {
				state = GameStateGameOver
			}
		}
	}
}

func handleGameLoop(scanner *bufio.Scanner, name string, conn *websocket.Conn) {
	go handleIncomingMessages(conn)
	handleOutgoingMessages(scanner, name, conn)
	conn.Close()
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
	handleGameLoop(scanner, name, conn)
}
