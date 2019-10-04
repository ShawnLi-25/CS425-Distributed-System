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
		membershipList <- MemListChan

		monitorList = msg.GetMonitorList(membershipList, LocalAddress)

		for _, v := range monitorList {
			monitorAdd := msg.GetIPAddressFromID(v)
			SendHearbeat(v, monitorAdd, LocalID)
		}
	}


	

}

func SendHeartbeat(monitorAddress string, monitorID string, localID string) {
	heartBeatMsg := msg.NewMessage(msg.HeartbeatMsg, localID, []string{})
	heartBeatPkg := msg.MsgToJSON(heartBeatMsg)

	for {
		udpAddr, err := net.ResolveUDPAddr(msg.ConnType, monitorAddress)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		conn, err := net.DialUDP(msg.ConnType, nil, udpAddr)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		msg, err := conn.Write(heartBeatPkg)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		log.Print("===HeartBeat Sent to: " + string(monitorID) + "\n" + "===Msg is" + string(msg))
		time.Sleep(time.Second) //send heartbeat 1 second
	}
}

func SendLeaveMsg(monitorAddress string, monitorID string, localID string) {
	leaveMsg := msg.NewMessage(msg.LeaveMsg, localID, []string{})
	leavePkg := msg.MsgToJSON(leaveMsg)

	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, monitorAddress)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	conn, err := net.DialUDP(msg.ConnType, nil, udpAddr)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	msg, err := conn.Write(leavePkg)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	log.Print("===LeaveMsg Sent to Introducer" + "\n" + "===Msg is" + string(msg))
}

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

	msg, err := conn.Write(joinMsg)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	fmt.Print("JoinMsg Sent to Introducer..."+ "\n" + "Msg is" + string(msg))

	//Set 3s Deadline for Ack
	conn.SetReadDeadline(time.Now().Add(time.Duration(3) * time.Second))

	//Read from Introducer 
	joinAck := make([]byte, 128)
	n, err := conn.Read(joinAck)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	
	joinAckMsg := msg.JSONToMsg([]byte(string(joinAck[:n])))
	fmt.Println("JoinAckMsg Received from Introducer...")

	if joinAckMsg.MessageType == msg.JoinAckMsg {
		curMembershipList := joinAckMsg.Content
		
		//Copy the current membership list locally 
		for _, member := range curMembershipList {
			UpQryChan <- UpdateQuery{1, member}
			<- MemListChan
		}
		return true
	} else {
		log.Println("Received Wrong Ack...")
		return false
	}
	return true
}


// func CreateID() string {
// 	hostName := msg.GetHostName()
// 	localTime := time.Now()
// 	// fmt.Println(localTime.Format(time.RFC3339))
// 	return hostName + ":" + msg.ConnPort + "+" + localTime.Format("20060102150405")
// }
