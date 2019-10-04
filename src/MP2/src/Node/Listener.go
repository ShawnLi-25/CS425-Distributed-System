package node

import (
	"fmt"
	"log"
	"net"

	msg "../Helper"
)


// Listener is a type that implements the ListenMsg(), ListenJoinMsg() "method"
type Listener struct {
}

func (l *Listener) NodeListen(port string) {
	fmt.Println("Initialize new listener...")
	con, err := net.ListenPacket(msg.ConnType, ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listen on port %s\n", port)
	defer con.Close()

}

//ListenMsg: Listen to Heartbeat or Leave Msg
func (l *Listener) ListenMsg() {

}

//ListenJoinMsg: Listen to Join Msg (Introducer-only)
func (l *Listener) ListenJoinMsg() {

}
