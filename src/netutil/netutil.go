package netutil

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/nadavash/bot-or-not/src/message"
	"google.golang.org/protobuf/proto"
)

func ReadProtoMessage(conn *websocket.Conn) (*message.WrapperMessage, error) {
	_, bytes, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	wrapperMsg := &message.WrapperMessage{}
	if err := proto.Unmarshal(bytes, wrapperMsg); err != nil {
		return nil, err
	}

	return wrapperMsg, nil
}

func SendProtoMessage(conn *websocket.Conn, msg *message.WrapperMessage) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling WrapperMessage proto:", err)
		return err
	}

	if err = conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
		log.Println("Error while writing proto message: ", err)
		return err
	}

	return nil
}
