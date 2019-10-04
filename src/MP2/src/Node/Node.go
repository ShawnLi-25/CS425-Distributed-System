package node

import (
	msg "../Helper"
)

var curNode Node = CreateNewNode()
var UpQryChan chan UpdateQuery = make(chan UpdateQuery)
var MemListChan chan []string = make(chan []string)
var LocalAddress string
var LocalID string

type Node struct {
	MemList []string
	InGroup bool
	Sender
	Listener
	Updater
	Introducer
}

func CreateNewNode() *Node {
	var newNode *Node = new(Node)
	newNode.MemList = []string{}
	newNode.InGroup = true
	// newSender := NewSender()
	// newListener := NewListener()
	// newIntroducer := NewIntroducer()
	// newUpdater := NewUpdater()
	// newNode := Node{
	// 	MemList:  newMemList,
	// 	Sender:   newSender,
	// 	Listener: newListener,
	// 	Updater:  newUpdater,
	// 	InGroup:  false,
	// }
	return newNode
}

//Called from main.go when the command is "JOIN\n"
//Create new node and run Listener,Sender and Updater
//in seperate goroutines
func RunNode(isIntroducer bool) {
	LocalID = msg.CreateID()
	LocalAddress = msg.GetHostName()

	go curNode.Updater.UpdateMembershipList(IntroducerAddress)

	//Firstly, send Join Msg to Introducer
	curNode.Sender.NodeSend(msg.JoinMsg)

	// curNode = CreateNewNode()
	if !isIntroducer {
		//false for non-intro, true for intro
		go curNode.Introducer.NodeHandleJoin()
		go curNode.Listener.NodeListen(false)
	} else {
		go curNode.Listener.NodeListen(true)
	}
}

//Called from main.go when the command is "LEAVE\n"
//Delete the Node
func StopNode() {
}
