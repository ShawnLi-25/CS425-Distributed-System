package node

import (
	"fmt"
	"log"
	"net"
	"time"

	msg "../Helper"
)




// Listener is a type that implements the ListenMsg(), ListenJoinMsg() "method"
type Listener struct {
}

func (l *Listener) RunMSGListener() {

	fmt.Println("Listener:Run message listener...")
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+msg.ConnPort)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listener:MSGListener listen on port %s\n", msg.ConnPort)

	for {

		select {
			case <-KillRoutine:
				ln.Close()
				fmt.Println("Listener: Leave!!")
				KillRoutine <- struct{}{}
				return
			default:
				fmt.Println("Listener: Works!!")
				HandleListenMsg(ln)
		}
	}

}

func HandleListenMsg(conn *net.UDPConn) {
	msgBuf := make([]byte, 1024)

	n, msgAddr, err := conn.ReadFromUDP(msgBuf)
	if err != nil {
		log.Fatal(err)
	}
	receivedMsg := msg.JSONToMsg([]byte(string(msgBuf[:n])))
	log.Printf("Listender: Recieve %s message from Node: %s, Addrs: %s", receivedMsg.MessageType, receivedMsg.NodeID, msgAddr)

	switch receivedMsg.MessageType {
	case msg.FailMsg:
			//Triggered by False-Positive Situation
		if receivedMsg.NodeID == LocalID {
			log.Printf("Listener: NodeID %s is recognized as failed...\n", LocalID)
			StopNode(false)
			return
		} else {
			fmt.Println("Listener: receive failMsg")
			UpQryChan <- UpdateQuery{2, receivedMsg.Content[0]}
			retMemList := <-MemListChan
			if len(retMemList) != 0 {
				SendFailMsg(conn, receivedMsg.Content[0])
			}
		}
	case msg.LeaveMsg:
		fmt.Println("Listener: receive leaveMsg")
		UpQryChan <- UpdateQuery{2, receivedMsg.Content[0]}
		retMemList := <-MemListChan
		if len(retMemList) != 0 {
			SendLeaveMsg(conn, receivedMsg.Content[0])
		}
	case msg.IntroduceMsg:
		fmt.Println("Listener: receive IntroduceMsg")
		UpQryChan <- UpdateQuery{1, receivedMsg.Content[0]}
		retMemList := <-MemListChan
		if len(retMemList) != 0 {
			SendIntroduceMsg(conn, receivedMsg.Content[0])
		}
	default:
		fmt.Println("Listener:Can't recognize the msg")
	}
	fmt.Println("Listener: Return from HandleListenMsg ")
	return
}

//Use MembershipList to update the key in MemHBMap(NodeID, Time)
func getMemHBMap(oldMemHBMap map[string]time.Time) map[string]time.Time {
	var MemHBMap map[string]time.Time
	MemHBList := msg.GetMonitoringList(MembershipList, LocalAddress)
	if len(oldMemHBMap) == 0 {//New MemHBMap
		for _, c := range MemHBList {
			MemHBMap[c] = time.Now()
		}
	} else {                   //old MemHBMap has values
		for NodeID, _ := range oldMemHBMap{
			for i, c := range MemHBList {
				if NodeID == c { //find same NodeID
					MemHBList = append(MemHBList[:i], MemHBList[i+1:]...) //delete
					break
				}
			}
			//Not found the NodeID
			delete(oldMemHBMap, NodeID)
		}
		for NodeID, resTime := range oldMemHBMap{
			MemHBMap[NodeID] = resTime
		}
		for _, c := range MemHBList{
			MemHBMap[c] = time.Now()
		}
	}
	return MemHBMap
}

//Counting the timeout
func HBTimer(ln *net.UDPConn, MemHBMap map[string]time.Time) {
	for{
		time.Sleep(time.Second)

		for NodeID, lastTime := range MemHBMap {
			timeDiff := time.Now().Sub(lastTime)
			if timeDiff - 2*msg.TimeOut*time.Second > 0{
				SendFailMsg(ln, NodeID)
			}
		}
		MemHBMap = getMemHBMap(MemHBMap)
	}
}


//Listen to Heartbeat and Check timeout
func (l *Listener) RunHBListener() {
	fmt.Println("HBListener:Initialize heartbeat listener...")
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+msg.HeartbeatPort)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HBListener:Listen Heartbeat on port %s\n", msg.HeartbeatPort)
	defer ln.Close()
	hbBuf := make([]byte, 1024)
	
	//Initialize MemHBMap
	var MemHBMap map[string]time.Time = make(map[string]time.Time)
	MemHBMap = getMemHBMap(MemHBMap)
	
	go HBTimer(ln, MemHBMap)
	//For-loop only update the value of MemHBMap(NodeID, Time)
	for {
		n, msgAddr, err := ln.ReadFromUDP(hbBuf)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Listener:Recieve Heartbeat from UDP client: %s", msgAddr)
		receivedMsg := msg.JSONToMsg([]byte(string(hbBuf[:n])))
		
		if receivedMsg.MessageType != msg.HeartbeatMsg {
			fmt.Println("Listener: HBlistener doesn't receive a HeartbeatMsg")
			continue
		}
		
		if _, ok := MemHBMap[receivedMsg.NodeID]; ok {
			MemHBMap[receivedMsg.NodeID] = time.Now()
		} else {
			fmt.Println("Listener: MemHBMap doesn't contain the NodeID"+receivedMsg.NodeID)
		}
	}
}
