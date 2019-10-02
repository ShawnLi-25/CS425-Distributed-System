package node

import (
	msg "MP2/src/helper"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func SendHeartbeat(monitorAddress string, monitorID string, localID string) {
	heartBeatMsg := msg.NewMessage("HeartBeat", localID, []string{})
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

func CreateID() string {
	cmd := exec.Command("hostname")
	hName, _ := cmd.Output()
	hostName := string(hName)
	localTime := time.Now()
	// fmt.Println(localTime.Format(time.RFC3339))
	return hostName + ":" + msg.ConnPort + "+" + localTime.Format("20060102150405")
}
