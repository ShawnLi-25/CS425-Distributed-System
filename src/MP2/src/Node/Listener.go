package node

import (
	"fmt"
	"log"
	"net"
	"time"

	msg "../Helper"
)

type Listener struct {
}

func buildUDPServer(ConnPort string) *net.UDPConn {
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
	ln := buildUDPServer(msg.ConnPort)
	fmt.Printf("===Listener:MSGListener listen on port %s\n", msg.ConnPort)

	for {
		select {
		case <-KillRoutine:
			ln.Close()
			fmt.Println("===Listener: MSGListener Leave!!")
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
			fmt.Println("===Listener: Receive failMsg")
			UpQryChan <- UpdateQuery{2, receivedMsg.Content[0]}
			retMemList := <-MemListChan
			if len(retMemList) != 0 {
				SendFailMsg(conn, receivedMsg.Content[0])
			}
		}
	case msg.LeaveMsg:
		fmt.Println("===Listener: Receive leaveMsg")
		UpQryChan <- UpdateQuery{2, receivedMsg.Content[0]}
		retMemList := <-MemListChan
		if len(retMemList) != 0 {
			SendLeaveMsg(conn, receivedMsg.Content[0])
		}
		fmt.Println("===Listener: Current Membershiplist is empty !!!!!!")
	case msg.IntroduceMsg:
		fmt.Println("===Listener: receive IntroduceMsg")
		UpQryChan <- UpdateQuery{1, receivedMsg.Content[0]}
		retMemList := <-MemListChan
		if len(retMemList) != 0 {
			SendIntroduceMsg(conn, receivedMsg.Content[0])
		}
	default:
		fmt.Println("===Listener:Can't recognize the msg")
	}
	fmt.Println("Listener: Return from HandleListenMsg ")
}

//Counting the timeout
func HBTimer(ln *net.UDPConn) {
	for {
		select {
		case <-KillRoutine:
			ln.Close()
			fmt.Println("===Listener: Timer Leave!!")
			KillRoutine <- struct{}{}
			return
		default:
			time.Sleep(time.Second)
			curTime := time.Now()
			for NodeID, lastTime := range MemHBMap {
				timeDiff := curTime.Sub(lastTime)
				if timeDiff-2*msg.TimeOut*time.Second > 0 {
					SendFailMsg(ln, NodeID)
				}
			}
		}
	}
}

//Listen to Heartbeat and Check timeout
func (l *Listener) RunHBListener() {

	ln := buildUDPServer(msg.HeartbeatPort)
	fmt.Printf("===HBListener:Listen Heartbeat on port %s\n", msg.HeartbeatPort)

	hbBuf := make([]byte, 2048)

	//fmt.Printf("\nListener:::RunHBListener:::MemHBMap has %d elements.\n\n",len(MemHBMap))

	go HBTimer(ln)
	//For-loop only update the value of MemHBMap(NodeID, Time)
	for {
		select {
		case <-KillRoutine:
			ln.Close()
			fmt.Println("===Listener: HBListener Leave!!")
			KillRoutine <- struct{}{}
			return
		default:
			n, _, err := ln.ReadFromUDP(hbBuf)
			if err != nil {
				log.Println(err)
			}

			receivedMsg := msg.JSONToMsg([]byte(string(hbBuf[:n])))

			if receivedMsg.MessageType != msg.HeartbeatMsg {
				fmt.Println("Listener: HBlistener doesn't receive a HeartbeatMsg")
				continue
			} else {
				fmt.Println("Listener:Recieve Heartbeat from NodeID:", receivedMsg.NodeID)
			}

			//fmt.Printf("\nListener:::For-loop:::MemHBMap has %d elements.\n\n",len(MemHBMap))

			if _, ok := MemHBMap[receivedMsg.NodeID]; ok {
				MemHBMap[receivedMsg.NodeID] = time.Now()
			} else {
				fmt.Println("Listener: MemHBMap doesn't contain the NodeID" + receivedMsg.NodeID)
			}
		}
	}
}
