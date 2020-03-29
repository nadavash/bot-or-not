package message

// Message defines a protocol for the client and server to send chat messages
// to each other.
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
