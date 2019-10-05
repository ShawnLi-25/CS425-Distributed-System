package node

import (
	"fmt"
	msg "../Helper"
)

var curNode *Node = CreateNewNode()
var UpQryChan chan UpdateQuery = make(chan UpdateQuery)
var MemListChan chan []string = make(chan []string)
var LocalAddress string
var LocalID string
var KILL chan int

type Node struct {
	Sender
	Listener
	Updater
	Introducer
}

func CreateNewNode() *Node {
	var newNode *Node = new(Node)
	return newNode
}

//Called from main.go when the command is "JOIN\n"
//Create new node and run Listener,Sender and Updater
//in seperate goroutines
func RunNode(isIntroducer bool) {
	LocalID = msg.CreateID()
	fmt.Println("Node: Local ID is: " + LocalID)
	LocalAddress = msg.GetHostName()
	fmt.Println("Node: Local Address is: " + LocalAddress)

	go curNode.Updater.UpdateMembershipList()
	if !isIntroducer {
		//Non-intro send JoinMsg to Introducer
		curNode.Sender.NodeSend(msg.JoinMsg)
	} else {
		//Introducer receive JoinMsg from non-intro
		go curNode.Introducer.NodeHandleJoin()
	}

	//go curNode.Listener.RunHBListener()
	go curNode.Listener.RunMSGListener()
	//go curNode.Sender.NodeSend(msg.HeartbeatMsg)
	k := <- KILL
	if k == 1 {
		return
	}
}

//Called from main.go when the command is "LEAVE\n"
//Delete the Node
func StopNode() {
	curNode.Sender.NodeSend(msg.LeaveMsg)
	KILL <- 1
}
