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

var MayFailMap map[string]time.Time = make(map[string]time.Time)

func buildUDPServer(ConnPort string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+ConnPort)
	fmt.Println("Build UDP!!!")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Build UDP Success!!!")

	return ln
}

func (l *Listener) RunMSGListener() {
	ln := buildUDPServer(msg.ConnPort)
	fmt.Printf("===Listener:MSGListener listen on port %s\n", msg.ConnPort)

	for {
		select {
		case <-KillMsgListener:
			ln.Close()
			fmt.Println("===Listener: MSGListener Leave!!")
			// KillRoutine <- struct{}{}
			return
		default:
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
	log.Printf("Listener: Recieve %s message from Node: %s, Addrs: %s", receivedMsg.MessageType, receivedMsg.NodeID, msgAddr)

	switch receivedMsg.MessageType {
	case msg.FailMsg:
		if receivedMsg.Content[0] == LocalID {
			log.Println("Fail Msg: I'm gonna Delete myself !!")
			fmt.Println("Fail Msg: I'm gonna Delete myself !!")
			StopNode()
		} else {
			// fmt.Println("Fail Msg: Delete Node!!")
			// newList := DeleteNode(receivedMsg.Content[0])
			// UpdateMemHBMap()
			// UpQryChan <- UpdateQuery{2, receivedMsg.Content[0]}
			// retMemList := <-MemListChan
			if len(newList) != 0 {
				//I have a update on MemList, so this is the first time I receive the msg
				//and I will send to other nodes this new msg!!!!!
				log.Printf("Listener: NodeID %s is recognized as failed...\n", receivedMsg.Content[0])
				SendFailMsg(conn, receivedMsg.NodeID, receivedMsg.Content[0])
			}
		}
	case msg.LeaveMsg:
		fmt.Println("Leave Msg: Delete Node!!")
		newList := MembershipList
		if LocalID != receivedMsg.Content[0] {
			newList = DeleteNode(receivedMsg.Content[0])
			UpdateMemHBMap()
		}
		// UpQryChan <- UpdateQuery{2, receivedMsg.Content[0]}
		// retMemList := <-MemListChan
		// fmt.Print(retMemList)
		if len(newList) != 0 {
			log.Printf("Listener: NodeID %s is recognized as leave...\n", receivedMsg.Content[0])
			SendLeaveMsg(conn, receivedMsg.NodeID, receivedMsg.Content[0])
		}
	case msg.IntroduceMsg:
		retMemList := AddNewNode(receivedMsg.Content[0])
		UpdateMemHBMap()
		// UpQryChan <- UpdateQuery{1, receivedMsg.Content[0]}
		// retMemList := <-MemListChan
		if len(retMemList) != 0 {
			log.Printf("Listener: NodeID %s join the group, welcome!\n", receivedMsg.Content[0])
			SendIntroduceMsg(conn, receivedMsg.NodeID, receivedMsg.Content[0])
		}
	default:
		fmt.Println("===Listener:Can't recognize the msg")
	}
	log.Println("Listener: Return from HandleListenMsg ")
}

//Counting the timeout
func HBTimer(ln *net.UDPConn) {
	for {
		select {
		case <-KillHBTimer:
			// ln.Close()
			fmt.Println("===Listener: Timer Leave!!")
			// KillRoutine <- struct{}{}
			return
		default:
			time.Sleep(time.Second)
			curTime := time.Now()
			for NodeID, lastTime := range MemHBMap {
				timeDiff := curTime.Sub(lastTime)
				fmt.Printf("===HBTimer: For %d duration not received message from %s!!===\n", int64(timeDiff), NodeID)
				log.Printf("===HBTimer: For %d duration not received message from %s!!===\n", int64(timeDiff), NodeID)
				_, ok := MayFailMap[NodeID]
				if ok {

					//Oops! This guy may fail!! Let me check
					if int64(timeDiff)-msg.TimeOut*int64(time.Millisecond) > 0 {
						//Ahaaaaaa! You fail!!!
						fmt.Printf("HBTimer: %s timeout!!\n", NodeID)
						newList := DeleteNode(NodeID)
						if len(newList) != 0 {
							//I have a update on MemList, so this is the first time I receive the msg
							//and I will send to other nodes this new msg!!!!!
							log.Printf("HBTimer: %s timeout!!\n", NodeID)
							UpdateMemHBMap()
							SendFailMsg(ln, "", NodeID)
						}

					} else {
						//Sorry, you are still good~
						delete(MayFailMap, NodeID)
					}
				} else {
					if int64(timeDiff)-msg.TimeOut*int64(time.Millisecond) > 0 {
						//You may fail...
						//BE CAREFUL! time.Now() is used to
						//record WHEN I add the NodeID to MayFailMap
						//No other use, just a record!!
						MayFailMap[NodeID] = time.Now()
					} else {
						//You are good guy~
					}
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

	go HBTimer(ln)
	//For-loop only update the value of MemHBMap(NodeID, Time)
	for {
		// select {
		// case <-KillHBListener:
		// 	KillHBTimer <- struct{}{}
		// 	ln.Close()
		// 	fmt.Println("===Listener: HBListener Leave!!")
		// 	return
		// default:
		n, _, err := ln.ReadFromUDP(hbBuf)
		if err != nil {
			log.Println(err)
		}

		receivedMsg := msg.JSONToMsg([]byte(string(hbBuf[:n])))

		log.Printf("Received Message Type: %s...\n", receivedMsg.MessageType)

		if receivedMsg.MessageType == msg.HeartbeatMsg {
			if _, ok := MemHBMap[receivedMsg.NodeID]; ok {
				MemHBMap[receivedMsg.NodeID] = time.Now()
			} else {
				// log.Println("Listener: MemHBMap doesn't contain the NodeID" + receivedMsg.NodeID)
			}
			// log.Println("Listener: HBlistener doesn't receive a HeartbeatMsg")
			continue
		} else if receivedMsg.MessageType == msg.LeaveMsg && receivedMsg.NodeID == LocalID {
			KillHBTimer <- struct{}{}
			ln.Close()
			fmt.Println("===Listener: HBListener Leave!!")
			return
			// log.Println("Listener:Recieve Heartbeat from NodeID:", receivedMsg.NodeID)
		}

		// }
	}
	return
}
