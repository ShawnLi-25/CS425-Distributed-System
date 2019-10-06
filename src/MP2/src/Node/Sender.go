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

var MonitorList []string

// Sender is a type that implements the SendHearbeat() "method"
type Sender struct{}

func (s *Sender) NodeSend(msgType string) {
	// var membershipList []string
	// var monitorList []string
	// localHostName := msg.GetHostName()

	//Join the group
	if msgType == msg.JoinMsg {
		joinSucceed := SendJoinMsg(msg.IntroducerAddress)

		if !joinSucceed {
			fmt.Println("Introducer is down!!")
			return
		}
	} else if msgType == msg.LeaveMsg {
		// UpQryChan <- UpdateQuery{0, ""}
		// membershipList =<- MemListChan
		isIntroducer := msg.IsIntroducer()
		if isIntroducer {
			fmt.Println("Close Introducer Port")
			msg.CloseIntroducePort(LocalID)
		}
		fmt.Println("Not Introducer??????")
		msg.CloseConnPort(LocalID)
	}
	return

}

func (s *Sender) SendHeartbeat() {
	heartBeatMsg := msg.NewMessage(msg.HeartbeatMsg, LocalID, []string{})
	heartBeatPkg := msg.MsgToJSON(heartBeatMsg)

	for {
		select {
		case <-KillRoutine:
			// ln.Close()
			fmt.Println("====Heartbeat Sender: Leave!!")
			KillRoutine <- struct{}{}
			return

		default:
			UpQryChan <- UpdateQuery{0, ""}
			membershipList := <-MemListChan

			MonitorList = msg.GetMonitorList(membershipList, LocalAddress)

			for _, monitorID := range MonitorList {
				monitorAddress := msg.GetIPAddressFromID(monitorID)
				udpAddr, err := net.ResolveUDPAddr(msg.ConnType, monitorAddress+":"+msg.HeartbeatPort)
				if err != nil {
					log.Println(err.Error())
					// os.Exit(1)
				}
				conn, err := net.DialUDP(msg.ConnType, nil, udpAddr)
				if err != nil {
					log.Println(err.Error())
					// os.Exit(1)
				}

				_, err = conn.Write(heartBeatPkg)
				if err != nil {
					log.Println(err.Error())
				}

				fmt.Printf("Sender: HeartBeat Sent to: %s...\n", monitorID)
				conn.Close()
			}
			time.Sleep(time.Second) //send heartbeat 1 second
		}
	}

}

func SendJoinMsg(introducerAddress string) bool {
	joinMsg := msg.NewMessage(msg.JoinMsg, LocalID, []string{})
	joinPkg := msg.MsgToJSON(joinMsg)

	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, introducerAddress+":"+msg.IntroducePort)
	if err != nil {
		log.Println(err.Error())
		// os.Exit(1)
	}
	conn, err := net.DialUDP(msg.ConnType, nil, udpAddr)
	if err != nil {
		log.Println(err.Error())
		// os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write(joinPkg)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	log.Println("Sender: JoinMsg Sent to Introducer...")

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

	log.Printf("Sender: JoinAckMsg Received from Introducer, the message type is: %s...: ", joinAckMsg.MessageType)

	if joinAckMsg.MessageType == msg.JoinAckMsg {
		curMembershipList := joinAckMsg.Content

		//Copy the current membership list locally
		for _, member := range curMembershipList {
			UpQryChan <- UpdateQuery{1, member}
			<-MemListChan
		}
		return true
	} else {
		log.Println("Sender: Received Wrong Ack...")
		return false
	}
	return true
}

func SendLeaveMsg(ln *net.UDPConn, leaveNodeID string) {

	leaveMsg := msg.NewMessage(msg.LeaveMsg, LocalID, []string{leaveNodeID})
	leavePkg := msg.MsgToJSON(leaveMsg)
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

		conn, err := net.DialUDP(msg.ConnType, nil, udpAddr)
		if err != nil {
			log.Println(err.Error())
			// os.Exit(1)
		}
		defer conn.Close()

		_, wErr := ln.WriteToUDP(leavePkg, udpAddr)
		if wErr != nil {
			log.Println(wErr.Error())
		}
		fmt.Printf("Sender:LeaveMsg Sent to Monitor: %s...\n", member)
	}

}

func SendIntroduceMsg(ln *net.UDPConn, newNodeID string) {
	introduceMsg := msg.NewMessage(msg.IntroduceMsg, LocalID, []string{newNodeID})
	introducePkg := msg.MsgToJSON(introduceMsg)
	monitorList := msg.GetMonitorList(MembershipList, LocalAddress)

	for _, member := range monitorList {
		// fmt.Println(i,member)
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
		log.Println("Sender:IntroduceMsg Sent to: " + member)
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
		fmt.Printf("Sender: FailMsg Sent to Monitor: %s...\n ", memberAddress)
	}

}
