package Message

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	MessageType string
	NodeID      string
	Content     []string
}

func NewMessage(Type string, ID string, Content []string) Message {
	newMessage := Message{
		MessageType: Type,
		NodeID:      ID,
		Content:     Content,
	}
	return newMessage
}

func MsgToJSON(message Message) []byte {
	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	return b
}

func JSONToMsg(b []byte) Message {
	var m Message
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
	}
	return m
}
