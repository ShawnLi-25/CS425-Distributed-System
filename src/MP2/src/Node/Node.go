package node

import (
	//"MP2/src/node"
	"log"
	"net"
)

type Node struct {
	MemList      []string
	InGroup      bool
	Sender       node.Sender
	Listener     node.Listener
	Updater      node.Updater
}

func CreateNewNode() Node{
	var newMemList []string
	newSender     := NewSender()
	newListener   := NewListener()
	newIntroducer := NewIntroducer()
	newUpdater    := NewUpdater()
	newNode := Node {
		MemList : newMemList
		Sender  : newSender
		Listener: newListener
		Updater : newUpdater
		InGroup : false
	}
	return newNode
}


//Called from main.go when the command is "JOIN\n"
//Create new node and run the node until LEAVE or crash
func RunNode(isIntroducer bool) {
	var node Node
	var	upQryChan := make(chan UpdateQuery)
	var memListChan = make(chan []string)

	node := CreateNewNode()
	if(!isIntroducer){
		//false for non-intro, true for intro
		go NodeListen(&node,false) 
	} else {
		go NodeListen(&node,true)
	}
	go NodeSend(&node)
	go NodeUpdate(&node)
}

//Called from main.go when the command is "LEAVE\n"
//Delete the Node
func StopNode() {
}
