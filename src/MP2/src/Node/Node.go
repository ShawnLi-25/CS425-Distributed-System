package node

import (
	"MP2/src/node"
	"log"
	"net"
)

type Node struct {
	node.Sender
	node.Listener
	node.Introducer
	node.Updater
}


//Called from main.go when the command is "JOIN\n"
//Create new node and run the node until LEAVE or crash
func RunNode(isIntroducer bool) {
	if(!isIntroducer){
	} else {
	}
}

//Called from main.go when the command is "LEAVE\n"
//Delete the Node
func StopNode() {
}
