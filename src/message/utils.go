package message

func WrapServerConnectionSuccessMessage(msg *ServerConnectionSuccessMessage) *WrapperMessage {
	return &WrapperMessage{
		MessageType: MessageType_SERVER_CONNECTION_SUCCESS,
		Message:     &WrapperMessage_ServerSuccess{msg},
	}
}

func WrapRoomConnectionSuccesMessage(msg *RoomConnectionSuccessMessage) *WrapperMessage {
	return &WrapperMessage{
		MessageType: MessageType_ROOM_CONNECTION_SUCCESS,
		Message:     &WrapperMessage_RoomSuccess{msg},
	}
}

func WrapChatMessage(msg *ChatMessage) *WrapperMessage {
	return &WrapperMessage{
		MessageType: MessageType_CHAT,
		Message:     &WrapperMessage_Chat{msg},
	}
}

func WrapGameOverMessage() *WrapperMessage {
	return &WrapperMessage{
		MessageType: MessageType_GAME_OVER,
		Message:     &WrapperMessage_GameOver{&GameOverMessage{}},
	}
}

func WrapPlayAgainMessage(msg *PlayAgainMessage) *WrapperMessage {
	return &WrapperMessage{
		MessageType: MessageType_PLAY_AGAIN,
		Message:     &WrapperMessage_PlayAgain{msg},
	}
}

func WrapBotOrNotMessage(msg *BotOrNotMessage) *WrapperMessage {
	return &WrapperMessage{
		MessageType: MessageType_BOT_OR_NOT,
		Message:     &WrapperMessage_BotOrNot{msg},
	}
}

func WrapAnswerCorrectMessage(msg *AnswerCorrectMessage) *WrapperMessage {
	return &WrapperMessage{
		MessageType: MessageType_ANSWER_CORRECT,
		Message:     &WrapperMessage_AnswerCorrect{msg},
	}
}
