package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
	"github.com/nadavash/bot-or-not/src/player"
)

var botRooms = make(map[uint32]*Room)
var humanRooms = make(map[uint32]*Room)
var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the initial GET request to a websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	email := r.FormValue("email")
	name := r.FormValue("name")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	humanPlayer := player.NewHumanPlayer(name, email, ws)
	msg := message.WrapServerConnectionSuccessMessage(
		&message.ServerConnectionSuccessMessage{
			WelcomeMessage: "You're connected to the Bot or Not server!",
		},
	)
	if humanPlayer.SendMessage(msg) != nil {
		return
	}
	AssignRoom(humanPlayer)
}

func AssignRoom(player player.Player) {
	goesToBotRoom := true
	if goesToBotRoom {
		botRoom := NewRoom(AssignRoom, true)
		botRooms[botRoom.GetRoomId()] = botRoom
		botRoom.AddPlayer(player)
	} else {
		// This may break if room is deleted while we are iterating through them
		for _, room := range humanRooms {
			if room.roomState == RoomStateWaiting {
				room.AddPlayer(player)
				return
			}
		}
		newRoom := NewRoom(AssignRoom, false)
		newRoom.AddPlayer(player)
		humanRooms[newRoom.GetRoomId()] = newRoom
	}
}

func main() {
	fmt.Println("V5")

	webDir := flag.String("web_dir", "./web", "Path to the web directory.")
	flag.Parse()

	// Start the web server on the default port.
	fs := http.FileServer(http.Dir(*webDir))
	go func() { http.ListenAndServe(":80", fs) }()

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
