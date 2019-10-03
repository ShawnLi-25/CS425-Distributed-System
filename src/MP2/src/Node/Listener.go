package node

import (
	"log"
	"net"
	"fmt"
)

const (
	connType = "udp"
)


// Listener is a type that implements the ListenMsg(), ListenJoinMsg() "method"
type Listener struct{
	Connection PacketConn
}

func NewListener(port string) Listener{
	fmt.Println("Initialize new listener...")
	con, err := net.ListenPacket(connType,":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listen on port %s\n", port)
	defer con.Close()

	newListener := Listener{
		Connection : con
	}
	return newListener
}

//ListenMsg: Listen to Heartbeat or Leave Msg
func (l *Listener) ListenMsg() {
	
}

//ListenJoinMsg: Listen to Join Msg (Introducer-only)
func (l *Listener) ListenJoinMsg() {

}
