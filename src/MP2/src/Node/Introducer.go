package node

import (
	"fmt"
	"log"
	"net"
	"os"
	msg "../Helper"
)

type Introducer struct{}

//Called from Node.go when the node type is Introducer
func (i *Introducer) NodeHandleJoin() {
	//Add Introducer itself to MemList
	UpQryChan <- UpdateQuery{1, LocalID}
	<-MemListChan

	//Set up UDP connection for upcoming JoinMsg
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+msg.IntroducePort)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println("Introducer: Start Listening for New-Join Node...")
	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Println(err.Error())
		// os.Exit(1)
	}
	fmt.Println("Introducer: Listening on port " + msg.IntroducePort)
	defer ln.Close()

	//Handle JoinMsg
	for {
		handleJoinMsg(ln)
	}
}

func handleJoinMsg(ln *net.UDPConn) {
	joinBuf := make([]byte, 1024)
	n, joinAddr, err := ln.ReadFromUDP(joinBuf)
	if err != nil {
		log.Println(err.Error())
	}

	joinMsg := msg.JSONToMsg([]byte(string(joinBuf[:n])))

	if joinMsg.MessageType == msg.JoinMsg {
		log.Println("Introducer: JoinMsg Received from Node:" + joinMsg.NodeID)

		//Send Introduce Message to Other node
		SendIntroduceMsg(ln,joinMsg.NodeID)

		//Add new node to introducer's merbership list
		UpQryChan <- UpdateQuery{1, joinMsg.NodeID}
		newMembershipList := <-MemListChan

		/*For debugging: print full MemList*/
		// for _, str := range newMembershipList {
		// 	fmt.Println("Introducer: Current Membership List is: " + str)
		// }

		//Send full membershiplist to new join node
		joinAckMsg := msg.NewMessage(msg.JoinAckMsg, LocalID, newMembershipList)
		joinAckPkg := msg.MsgToJSON(joinAckMsg)

		_, err := ln.WriteToUDP(joinAckPkg, joinAddr)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("Introducer: JoinAck Sent to Node:" + joinMsg.NodeID)
	} else {
		log.Printf("Introducer: Port %s only accepts JoinMsg\n", msg.IntroducePort)
	}
}
