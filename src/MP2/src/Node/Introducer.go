package node

import (
	"fmt"
	"log"
	"net"
	"os"

	msg "../Helper"
)

// Introducer is a type that implements the SendFullListToNewNode(), SendIntroduceMsg() "method"
type Introducer struct{}

func handleUDPConnection(conn *net.UDPConn){
	joinBuf := make([]byte, 1024)
	n, joinAddr, err := ln.ReadFromUDP(joinBuf)
	if err != nil {
		log.Println(err.Error())
	}
<<<<<<< HEAD
	fmt.Println("Start Listening for New-Join Node...")
	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	for {
		joinBuf := make([]byte, 128)
		n, joinAddr, err := ln.ReadFromUDP(joinBuf)
		if err != nil {
			log.Println(err.Error())
		}

		joinMsg := msg.JSONToMsg([]byte(string(joinBuf[:n])))

		if joinMsg.MessageType == msg.JoinMsg {
			fmt.Println("JoinMsg Received from Node... Address: " + joinAddr.IP.String())

			//Send Introduce Message to Other node
			SendIntroduceMsg(ln, joinMsg.NodeID)
=======

	joinMsg := msg.JSONToMsg([]byte(string(joinBuf[:n])))

	if joinMsg.MessageType == msg.JoinMsg {
		fmt.Println("JoinMsg Received from Node... Address: " + joinAddr.IP.String())

		//Send Introduce Message to Other node
		// SendIntroduceMsg()

		introduceMsg := msg.NewMessage(msg.JoinAckMsg, LocalID, []string{joinMsg.NodeID})
		introducePkg := msg.MsgToJSON(introduceMsg)

		for _, member := range MembershipList {
			if member == LocalID {
				continue
			}
			memberAddress := msg.GetIPAddressFromID(member)
>>>>>>> 4c882ff832057081423a7d87c214e79412c85e17

			udpAddr, err := net.ResolveUDPAddr(msg.ConnType, memberAddress+":"+msg.ConnPort)
			if err != nil {
				log.Println(err.Error())
			}
			_, wErr := ln.WriteToUDP(introducePkg, udpAddr)
			if wErr != nil {
				log.Println(wErr.Error())
			}
		}

		//Add new node to introducer's merbership list
		UpQryChan <- UpdateQuery{1, joinMsg.NodeID}
		newMembershipList := <-MemListChan
		//Send full membershiplist to new join node
		// SendJoinAckMsg(addr)

		joinAckMsg := msg.NewMessage(msg.JoinAckMsg, LocalID, newMembershipList)
		joinAckPkg := msg.MsgToJSON(joinAckMsg)

		msg, err := ln.WriteToUDP(joinAckPkg, joinAddr)
		if err != nil {
			log.Println(err.Error())
		}

		log.Print("JoinAck Sent to New Node:" + joinMsg.NodeID + "\nMsg is" + string(msg))
	}
}


func (i *Introducer) NodeHandleJoin() {
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+msg.IntroducePort)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println("Start Listening for New-Join Node...")
	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on port %s", msg.IntroducePort)
	defer ln.Close()

	for {
		handleUDPConnection(ln)
	}
}

// func SendIntroduceMsg() {

// }

// func SendJoinAckMsg(addr *net.UDPAddr) {

// }
