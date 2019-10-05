package node

import (
	"fmt"
	"log"
	"net"
	// "os"
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
	log.Println("Introducer: Listening on port " + msg.IntroducePort)

	//Handle JoinMsg

	for {
		select {
			case <-KillRoutine:
				ln.Close()
				log.Println("Introducer: Leave!!")
				return
			default:
				log.Println("Introducer: Works!!")
				HandleJoinMsg(ln)
		}
	}

}

func HandleJoinMsg(ln *net.UDPConn) {
	joinBuf := make([]byte, 1024)
	n, joinAddr, err := ln.ReadFromUDP(joinBuf)
	if err != nil {
		log.Println(err.Error())
	}

	joinMsg := msg.JSONToMsg([]byte(string(joinBuf[:n])))

	if joinMsg.MessageType == msg.JoinMsg {
		log.Printf("Introducer: JoinMsg Received from Node: %s...\n", joinMsg.NodeID)

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
		log.Printf("Introducer: JoinAck Sent to Node: %s...\n", joinMsg.NodeID)
	} else if joinMsg.MessageType == msg.LeaveMsg{
		log.Printf("Introducer: Introducer Leave... Close Port:%s...\n", msg.IntroducePort)
	} 
	return 
}
