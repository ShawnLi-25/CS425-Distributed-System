package node

import (
	msg "../Helper"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
	"strings"
)

// Sender is a type that implements the SendHearbeat() "method"
type Sender struct{}


func (s *Sender) NodeSend() {
	var membershipList []string
	var monitorList []string
	var localIdx int
	localHostName := msg.GetHostName()

	for {
		upQryChan<-UpdateQuery{0,""}
		membershipList <- memListChan

		monitorList = msg.GetMonitorList(membershipList, localHostName)
		
		
		
	}
}

func (s *Sender) SendHeartbeat(monitorAddress string, monitorID string, localID string) {
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

func (s *Sender) SendLeaveMsg(monitorAddress string, monitorID string, localID string) {
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

		log.Print("===LeaveMsg Sent to: " + string(monitorID) + "\n" + "===Msg is" + string(msg))
	}
}

func (s *Sender) SendJoinMsg(introducerAddress string, localID string) {
	joinMsg := msg.NewMessage(msg.JoinMsg, localID, []string{})
	joinPkg := msg.MsgToJSON(joinMsg)

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

		log.Print("===LeaveMsg Sent to: " + string(monitorID) + "\n" + "===Msg is" + string(msg))
}

func CreateID() string {
	hostName := msg.GetHostName()
	localTime := time.Now()
	// fmt.Println(localTime.Format(time.RFC3339))
	return hostName + ":" + msg.ConnPort + "+" + localTime.Format("20060102150405")
}
