package Sender

import (
	msg "Helper/Message"
	"net"
)

func SendHeartbeat(monitor *net.UDPAddr, monitorID string, localID string) {
	heartBeatMsg := msg.NewMessage("HeartBeat", localID, []string{""})

}
