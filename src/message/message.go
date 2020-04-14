package message

type MessageType int

const (
	MessageTypeChat                    = 0
	MessageTypeServerConnectionSuccess = 1
	MessageTypeRoomConnectionSuccess   = 2
	MessageTypeGameOver                = 3
	MessageTypeBotOrNotAnswer          = 4
	MessagePlayAgain                   = 4
)

type MessageBase struct {
	MessageType MessageType `json:"messageType"`
	MessageBody interface{} `json:"messageBody"`
}

// Message defines a protocol for the client and server to send chat messages
// to each other.
type ChatMessage struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type ServerConnectionSuccessMessage struct {
	WelcomeMessage string `json:welcomeMessage`
}

type RoomConnectionSuccessMessage struct {
	RoomId uint32 `json:"roomId"`
}

type GameOverMessage struct{}

type PlayAgainMessage struct{
	PlayAgain bool
}

type BotOrNotAnswerMessage struct {
	ArePlayersBotsAnswer bool
}
