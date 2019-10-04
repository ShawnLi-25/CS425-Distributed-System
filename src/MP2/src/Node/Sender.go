package node

import (
	"fmt"
	"log"
	"net"
	"os"

	msg "../Helper"

	//"os/exec"
	"time"
	//"strings"
)

// Sender is a type that implements the SendHearbeat() "method"
type Sender struct{}

func (s *Sender) NodeSend(msgType string) {
	var membershipList []string
	var monitorList []string
	// localHostName := msg.GetHostName()

	//Join the group
	if msgType == msg.JoinMsg {
		joinSucceed := SendJoinMsg(msg.IntroducerAddress)

		if !joinSucceed {
			fmt.Println("Introducer is down!!")
			return
		}
	} else if msgType == msg.HeartbeatMsg {
		UpQryChan <- UpdateQuery{0, ""}
		membershipList =<- MemListChan

		monitorList = msg.GetMonitorList(membershipList, LocalAddress)

		for _, v := range monitorList {
			monitorAdd := msg.GetIPAddressFromID(v)
			SendHeartbeat(v, monitorAdd, LocalID)
		}
	} else if msgType == msg.LeaveMsg {
		UpQryChan <- UpdateQuery{0, ""}
		membershipList =<- MemListChan

		monitorList = msg.GetMonitorList(membershipList, LocalAddress)

		for _, v := range monitorList {
			monitorAdd := msg.GetIPAddressFromID(v)
			SendLeaveMsg(v, monitorAdd, LocalID)
		}
	}


}

func SendHeartbeat(monitorAddress string, monitorID string, localID string) {
	heartBeatMsg := msg.NewMessage(msg.HeartbeatMsg, localID, []string{})
	heartBeatPkg := msg.MsgToJSON(heartBeatMsg)

	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, monitorAddress + ":" + msg.HeartbeatPort)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	conn, err := net.DialUDP(msg.ConnType, nil, udpAddr)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	for {
		msg, err := conn.Write(heartBeatPkg)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		fmt.Print("HeartBeat Sent to: " + string(monitorID) + "\nMsg is" + string(msg))
		time.Sleep(time.Second) //send heartbeat 1 second
	}
	conn.Close()
}

// func SendLeaveMsg(ln *net.UDPConn, failNodeID string) {
// 	leaveMsg := msg.NewMessage(msg.LeaveMsg, localID, []string{})
// 	leavePkg := msg.MsgToJSON(leaveMsg)

// 	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, monitorAddress + ":" + msg.ConnPort)
// 	if err != nil {
// 		log.Println(err.Error())
// 		os.Exit(1)
// 	}
// 	conn, err := net.DialUDP(msg.ConnType, nil, udpAddr)
// 	if err != nil {
// 		log.Println(err.Error())
// 		os.Exit(1)
// 	}
// 	defer conn.Close()

// 	msg, err := conn.Write(leavePkg)
// 	if err != nil {
// 		log.Println(err.Error())
// 		os.Exit(1)
// 	}
// 	fmt.Print("LeaveMsg Sent to: " + string(monitorID) + "\nMsg is" + string(msg))
// }

func SendJoinMsg(introducerAddress string) bool{
	joinMsg := msg.NewMessage(msg.JoinMsg, LocalID, []string{})
	joinPkg := msg.MsgToJSON(joinMsg)

	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, introducerAddress + ":" + msg.IntroducePort)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	conn, err := net.DialUDP(msg.ConnType, nil, udpAddr)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write(joinPkg)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Sender: JoinMsg Sent to Introducer...")

	//Set 3s Deadline for Ack
	conn.SetReadDeadline(time.Now().Add(time.Duration(3) * time.Second))

	//Read from Introducer 
	joinAck := make([]byte, 2048)
	n, err := conn.Read(joinAck)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	
	joinAckMsg := msg.JSONToMsg([]byte(string(joinAck[:n])))
	
	fmt.Printf("Sender: JoinAckMsg Received from Introducer, the message type is: %s...: ", joinAckMsg.MessageType)

	if joinAckMsg.MessageType == msg.JoinAckMsg {
		curMembershipList := joinAckMsg.Content
		
		//Copy the current membership list locally 
		for _, member := range curMembershipList {
			UpQryChan <- UpdateQuery{1, member}
			<- MemListChan
		}
		return true
	} else {
		log.Println("Sender: Received Wrong Ack...")
		return false
	}
	return true
}

func SendIntroduceMsg(ln *net.UDPConn, newNodeID string) {
	introduceMsg := msg.NewMessage(msg.IntroduceMsg, LocalID, []string{newNodeID})
	introducePkg := msg.MsgToJSON(introduceMsg)

	monitorList := msg.GetMonitorList(MembershipList, LocalAddress)

	for _, member := range monitorList {
		if member == LocalID {
			continue
		}

		memberAddress := msg.GetIPAddressFromID(member)

		udpAddr, err := net.ResolveUDPAddr(msg.ConnType, memberAddress+":"+msg.ConnPort)
		if err != nil {
			log.Println(err.Error())
		}

		_, wErr := ln.WriteToUDP(introducePkg, udpAddr)
		if wErr != nil {
			log.Println(wErr.Error())
		}
	fmt.Print("IntroduceMsg Sent to: "+ string(member))
	}
}

func SendFailMsg(ln *net.UDPConn, failNodeID string) {
	failMsg := msg.NewMessage(msg.FailMsg, LocalID, []string{failNodeID})
	failPkg := msg.MsgToJSON(failMsg)

	monitorList := msg.GetMonitorList(MembershipList, LocalAddress)

	for _, member := range monitorList {
		if member == LocalID {
			continue
		}

		memberAddress := msg.GetIPAddressFromID(member)

		udpAddr, err := net.ResolveUDPAddr(msg.ConnType, memberAddress+":"+msg.ConnPort)
		if err != nil {
			log.Println(err.Error())
		}

		_, wErr := ln.WriteToUDP(failPkg, udpAddr)
		if wErr != nil {
			log.Println(wErr.Error())
		}
	fmt.Print("FailMsg Sent to: "+ string(member))
	}

}

// func CreateID() string {
// 	hostName := msg.GetHostName()
// 	localTime := time.Now()
// 	// fmt.Println(localTime.Format(time.RFC3339))
// 	return hostName + ":" + msg.ConnPort + "+" + localTime.Format("20060102150405")
// }
