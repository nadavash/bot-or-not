package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nadavash/bot-or-not/src/message"
	"github.com/nadavash/bot-or-not/src/player"
)

type RoomState int

const (
	RoomStateWaiting    RoomState = 0
	RoomStateInProgress RoomState = 1
	RoomStateFinished   RoomState = 2
	gameTimeSeconds               = 20
	roomLimit                     = 2
)

type Room struct {
	roomId           uint32
	roomState        RoomState
	players          []player.Player
	broadcastChannel chan playerMessagePair
	finishedCallback func(player.Player)
	isBotRoom        bool
}

type playerMessagePair struct {
	message *message.WrapperMessage
	player  player.Player
}

func NewRoom(onPlayerFinished func(player.Player), botRoom bool) *Room {
	r := new(Room)
	r.roomId = rand.Uint32() % 10000
	r.players = make([]player.Player, 0, roomLimit)
	r.roomState = RoomStateWaiting
	r.broadcastChannel = make(chan playerMessagePair)
	r.finishedCallback = onPlayerFinished
	r.isBotRoom = botRoom
	if botRoom {
		go r.generateBotPlayers()
	}
	return r
}

func (r *Room) GetRoomId() uint32 {
	return r.roomId
}

func (r *Room) AddPlayer(player player.Player) error {
	fmt.Println("addingPlayer")
	fmt.Println(len(r.players))
	if r.roomState != RoomStateWaiting {
		return errors.New("Cannot add players to a Room that's in progress or finished")
	}
	r.players = append(r.players, player)
	if len(r.players) == roomLimit {
		fmt.Println("4 players added")
		go r.test()
		go r.broadcastMessages()
		r.roomState = RoomStateInProgress
		r.sendRoomMessage("Room full! Game starting now.")
		go r.handleGameLogic()
		for _, player := range r.players {
			fmt.Println("accepting from player")
			go r.acceptIncomingChatMessages(player)
		}
	}

	msg := message.WrapRoomConnectionSuccesMessage(
		&message.RoomConnectionSuccessMessage{
			RoomId: r.roomId,
		},
	)
	return player.SendMessage(msg)
}

func (r *Room) test() {
	fmt.Println("in test")
}

func (r *Room) acceptIncomingChatMessages(player player.Player) {
	for r.roomState == RoomStateInProgress {
		// Read in a new message as JSON and map it to a Message object
		wrapperMsg, err := player.ReceiveMessage()
		if err != nil {
			log.Printf("Error occurred while reading proto message: %v", err)
			r.removePlayer(player)
			// TODO: instead of closing the connection, boot all of the players out
			// and force them to search for a new game.
			// player.Close()
			return
		}

		if r.roomState == RoomStateInProgress {
			chatMsg := wrapperMsg.GetChat()
			r.broadcastChannel <- playerMessagePair{
				player:  player,
				message: message.WrapChatMessage(chatMsg),
			}
		} else if r.roomState == RoomStateFinished {
			decisionMsg := wrapperMsg.GetBotOrNot()

			fmt.Println("Message:", wrapperMsg)
			fmt.Println("accepting decisions from players")
			fmt.Println("Decision received:", decisionMsg)
			log.Printf("decision: %v", decisionMsg.ArePlayersBots)
			isAnswerCorrect := decisionMsg.ArePlayersBots == r.isBotRoom
			player.SendMessage(
				message.WrapAnswerCorrectMessage(
					&message.AnswerCorrectMessage{
						IsCorrectAnswer: isAnswerCorrect,
					},
				),
			)

			r.roomState = RoomStateFinished
			r.acceptPlayAgain(player)
		}
	}
}

func (r *Room) acceptPlayAgain(player player.Player) {
	fmt.Println("about to accept play again")
	_, err := player.ReceiveMessage()
	fmt.Println("reading play again")
	if err != nil {
		log.Printf("error: %v", err)
		r.removePlayer(player)
		//player.Close()
	}
	r.finishedCallback(player)
}

func (r *Room) broadcastMessages() {
	fmt.Println("starting broadcast messages")
	for r.roomState == RoomStateInProgress || len(r.broadcastChannel) > 0 {
		fmt.Println("in loop")
		playerMessage := <-r.broadcastChannel
		fmt.Println(playerMessage.message)
		for _, player := range r.players {
			if player == playerMessage.player {
				continue
			}
			err := player.SendMessage(playerMessage.message)
			if err != nil {
				log.Printf("error: %v", err)
				// player.Close()
			}
		}
	}
}

func (r *Room) handleGameLogic() {
	fmt.Println("4 player")
	for secondsLeft := gameTimeSeconds; secondsLeft > 0; secondsLeft -= 60 {
		r.sendRoomMessage(
			fmt.Sprintf("%d seconds left in the game", secondsLeft))

		if secondsLeft < 60 {
			time.Sleep(time.Duration(secondsLeft) * time.Second)
		} else {
			time.Sleep(time.Minute)
		}
	}

	r.sendRoomMessage("Times up, game is over!")
	r.roomState = RoomStateFinished
	r.broadcastChannel <- playerMessagePair{
		message.WrapGameOverMessage(),
		nil,
	}
}

func (r *Room) sendRoomMessage(s string) {
	go func() {
		r.broadcastChannel <- playerMessagePair{
			message.WrapChatMessage(
				&message.ChatMessage{
					Username: "Room",
					Message:  s,
				},
			),
			nil,
		}
	}()
}

func (r *Room) removePlayer(player player.Player) {
	for i, playerPointer := range r.players {
		if player == playerPointer {
			r.players = append(r.players[i:], r.players[i+1:]...)
			break
		}
	}
}

func (r *Room) generateBotPlayers() {
	botPlayers := 0
	instantBotPlayers := rand.Intn(roomLimit + 1)
	botPlayers += instantBotPlayers
	for i := 0; i < instantBotPlayers; i++ {
		botPlayer := player.NewBotPlayer()
		r.AddPlayer(botPlayer)
	}
	for ; botPlayers < roomLimit-1; botPlayers++ {
		time.Sleep(time.Second * time.Duration(rand.Intn(30)))
		botPlayer := player.NewBotPlayer()
		r.AddPlayer(botPlayer)
	}
}
