package node

import (
	msg "../Helper"
	"fmt"
	"net"
	"os"
	"log"
)

// Introducer is a type that implements the SendFullListToNewNode(), SendIntroduceMsg() "method"
type Introducer struct{}

func (i *Introducer) NodeHandleJoin() {
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+ msg.IntroducePort)
	if err != nil {
		log.Println(err.Error())
	}	
	fmt.Println("Start Listening for New-Join Node...")
	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	for {
		joinBuf := make([]byte, 128)
		n, addr, err := ln.ReadFromUDP(joinBuf)
		if err != nil {
			log.Println(err.Error())
		}

		joinMsg := msg.JSONToMsg([]byte(string(joinBuf[:n])))

		if joinMsg.MessageType == msg.JoinMsg {
			fmt.Println("JoinMsg Received from Node... Address: "+ addr.IP.String())
				
		}
	}
}

func (i *Introducer) SendIntroduceMsg() {

}

func SendJoinAckMsg(newJoinAddress string, newJoinID string) {

	msg, err := conn.Write(msg.JoinAckMsg)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	log.Print("===JoinAck Sent to " + "\n" + "===Msg is" + string(msg))
}
