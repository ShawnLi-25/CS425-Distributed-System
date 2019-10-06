package node

import (
	"fmt"
	"log"
	"net"
	"time"
	msg "../Helper"
)

var MemHBMap map[string]time.Time = make(map[string]time.Time)



type Listener struct {
}

func buildUDPServer(ConnPort string) *net.UDPConn{
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+ConnPort)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	return ln
}


func (l *Listener) RunMSGListener() {
	fmt.Println("here")

	ln := buildUDPServer(msg.ConnPort)
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
	var newMemHBMap map[string]time.Time = make(map[string]time.Time)
	MemHBList := msg.GetMonitoringList(MembershipList, LocalAddress)
	fmt.Println("Listener: Correct Monitoring List is:...")
	fmt.Print(MemHBList)
	fmt.Println("..")
	if len(oldMemHBMap) == 0 {//New MemHBMap
		for _, c := range MemHBList {
			newMemHBMap[c] = time.Now()
			fmt.Println("===Map insert new element===")
			fmt.Print(newMemHBMap)
		
		}
	} else {                   //old MemHBMap has values
		for _, c := range MemHBList {
			if LastTime, ok := oldMemHBMap[c]; ok {
				newMemHBMap[c] = LastTime
			fmt.Println("===Map insert new element===")
			fmt.Print(newMemHBMap)

			} else {
				newMemHBMap[c] = time.Now()
			fmt.Println("===Map insert new element===")
			fmt.Print(newMemHBMap)

			}
		}
	}
	//fmt.Printf("\nListener:::getMem:::MemHBMap has %d elements.\n\n",len(MemHBMap))
	fmt.Println("Listener: Current Monitoring List is:...")
	fmt.Print(newMemHBMap)
	fmt.Println("..")
	return newMemHBMap

}

//Counting the timeout
func HBTimer(ln *net.UDPConn) {
	for{
		time.Sleep(time.Second)
		curTime := time.Now()
		for NodeID, lastTime := range MemHBMap {
			timeDiff := curTime.Sub(lastTime)
			if timeDiff - 2*msg.TimeOut*time.Second > 0{
				SendFailMsg(ln, NodeID)
			}
		}
		MemHBMap = getMemHBMap(MemHBMap)

		//fmt.Printf("\nListener:::HBTimer:::MemHBMap has %d elements.\n\n",len(MemHBMap))

	}
}


//Listen to Heartbeat and Check timeout
func (l *Listener) RunHBListener() {

	ln := buildUDPServer(msg.HeartbeatPort)
	fmt.Printf("HBListener:Listen Heartbeat on port %s\n", msg.HeartbeatPort)

	hbBuf := make([]byte, 1024)
	
	MemHBMap = getMemHBMap(MemHBMap)
	//fmt.Printf("\nListener:::RunHBListener:::MemHBMap has %d elements.\n\n",len(MemHBMap))
	
	go HBTimer(ln)
	//For-loop only update the value of MemHBMap(NodeID, Time)
	for {
		n, _, err := ln.ReadFromUDP(hbBuf)
		if err != nil {
			log.Println(err)
		}

		receivedMsg := msg.JSONToMsg([]byte(string(hbBuf[:n])))
		fmt.Println("Listener:Recieve Heartbeat from NodeID:", receivedMsg.NodeID)
		
		if receivedMsg.MessageType != msg.HeartbeatMsg {
			fmt.Println("Listener: HBlistener doesn't receive a HeartbeatMsg")
			continue
		}
		
		//fmt.Printf("\nListener:::For-loop:::MemHBMap has %d elements.\n\n",len(MemHBMap))

		if _, ok := MemHBMap[receivedMsg.NodeID]; ok {
			MemHBMap[receivedMsg.NodeID] = time.Now()
		} else {
			fmt.Println("Listener: MemHBMap doesn't contain the NodeID"+receivedMsg.NodeID)
		}
	}
}
