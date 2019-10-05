package node

import (
	"fmt"
	"log"
	"net"
	"time"

	msg "../Helper"
)

// Listener is a type that implements the ListenMsg(), ListenJoinMsg() "method"
type Listener struct {
}

func (l *Listener) RunMSGListener() {

	fmt.Println("Listener:Run message listener...")
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+msg.ConnPort)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listener:MSGListener listen on port %s\n", msg.ConnPort)

	for {

		select {
			case <-KillRoutine:
				ln.Close()
				fmt.Println("Listener: Leave!!")
				return
			default:
				fmt.Println("Listener: Works!!")
				HandleListenMsg(ln)
		}
	}

}

func HandleListenMsg(conn *net.UDPConn) {
	msgBuf := make([]byte, 1024)

	n, msgAddr, err := conn.ReadFromUDP(msgBuf)
	if err != nil {
		log.Fatal(err)
	}
	receivedMsg := msg.JSONToMsg([]byte(string(msgBuf[:n])))
	log.Printf("Listender: Recieve %s message from Node: %s, Addrs: %s", receivedMsg.MessageType, receivedMsg.NodeID, msgAddr)

	switch receivedMsg.MessageType {
	case msg.FailMsg:
		UpQryChan <- UpdateQuery{2, receivedMsg.NodeID}
		retMemList := <-MemListChan
		if len(retMemList) != 0 {
			SendFailMsg(conn, receivedMsg.NodeID)
		}
	case msg.LeaveMsg:
		UpQryChan <- UpdateQuery{2, receivedMsg.Content[0]}
		retMemList := <-MemListChan
		if len(retMemList) != 0 {
			SendLeaveMsg(conn, receivedMsg.Content[0])
		}
	case msg.IntroduceMsg:
		fmt.Println("Listener: receive IntroduceMsg")
		UpQryChan <- UpdateQuery{1, receivedMsg.Content[0]}
		retMemList := <-MemListChan
		if len(retMemList) != 0 {
			SendIntroduceMsg(conn, receivedMsg.Content[0])
		}
	default:
		fmt.Println("Listener:Can't recognize the msg")
	}
	fmt.Println("Listener: Return from HandleListenMsg ")
	return
}

//Listen to Heartbeat and Check timeout
func (l *Listener) RunHBListener() {
	fmt.Println("HBListener:Initialize heartbeat listener...")
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+msg.HeartbeatPort)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HBListener:Listen Heartbeat on port %s\n", msg.HeartbeatPort)
	defer ln.Close()
	hbBuf := make([]byte, 1024)
	ln.SetReadDeadline(time.Now().Add(2*msg.TimeOut*time.Second))
	for {
		n, msgAddr, err := ln.ReadFromUDP(hbBuf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Listener:Recieve Heartbeat from UDP client: %s", msgAddr)
		if n > 0 {
			//No delay, refresh deadline
			ln.SetReadDeadline(time.Now().Add(2*msg.TimeOut*time.Second))
			receivedMsg := msg.JSONToMsg([]byte(string(hbBuf[:n])))
			msg.PrintMsg(receivedMsg)
		}

		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			//Timeout error
			fmt.Printf("HBListener: Client %s Timeout!\n", msgAddr)
			//TODO Send timeout msg
		}
	}
}
