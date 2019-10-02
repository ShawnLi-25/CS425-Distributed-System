package helper

import (
	"encoding/json"
	"fmt"
)

const (
	ConnHostName  = "fa19-cs425-g73-%02d.cs.illinois.edu"
	ConnType      = "udp"
	ConnPort      = "8888"
	ConnlocalHost = "localhost"
	TimeOut       = 1
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
	var m Message
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
	}
	return m
}
