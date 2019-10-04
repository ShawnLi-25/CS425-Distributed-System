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

func handleListenMsg(conn *net.UDPConn){
	msgBuf := make([]byte, 124)

	n, msgAddr, err := conn.ReadFromUDP(msgBuf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Recieve Msg from UDP client: %s", msgAddr)

	receivedMsg := msg.JSONToMsg([]byte(string(msgBuf[:n])))
	msg.PrintMsg(receivedMsg)
	switch receivedMsg.MessageType {
		case msg.HeartbeatMsg:
			fmt.Println("===Receive Heartbeat===")
		case msg.JoinMsg:
			fmt.Println("===Receive JoinMsg===")
		case msg.FailMsg:
			fmt.Println("===Receive FailMsg===")
		case msg.LeaveMsg:
			fmt.Println("===Receive LeaveMsg===")
		case msg.IntroduceMsg:
			fmt.Println("===Receive IntroduceMsg===")
		default:
			fmt.Println("Can't recognize the msg")
	}
}

func (l *Listener) NodeListen() {
	fmt.Println("Initialize new listener...")
	udpAddr,err := net.ListenPacket(msg.ConnType, ":"+msg.ConnPort)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listen on port %s\n", port)
	defer ln.Close()
	
	for {
		handleListenMsg(ln)
	}
}

//ListenMsg: Listen to Heartbeat or Leave Msg
func (l *Listener) ListenMsg() {

}

//ListenJoinMsg: Listen to Join Msg (Introducer-only)
func (l *Listener) ListenJoinMsg() {

}
