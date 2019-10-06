package node

import (
	"fmt"
	//"log"
	"time"

	msg "../Helper"
)

var curNode *Node = CreateNewNode()
var UpQryChan chan UpdateQuery = make(chan UpdateQuery)
var MemListChan chan []string = make(chan []string)
var KillRoutine chan struct{} = make(chan struct{})

// var Kill
var LocalAddress string
var LocalID string
var Status bool

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
	Status = true

	go curNode.Updater.UpdateMembershipList()
	go curNode.Listener.RunMSGListener()
	if !isIntroducer {
		//Non-intro send JoinMsg to Introducer
		curNode.Sender.NodeSend(msg.JoinMsg)
	} else {
		//Introducer receive JoinMsg from non-intro
		go curNode.Introducer.NodeHandleJoin()
	}

	go curNode.Listener.RunHBListener()

	go curNode.Sender.SendHeartbeat()
}

//Called from main.go when the command is "LEAVE\n"
//Delete the Node
func StopNode(byLocal bool) {
	if byLocal {
		curNode.Sender.NodeSend(msg.LeaveMsg)
	}
	Status = false
	KillRoutine <- struct{}{}
	fmt.Println("Node: Stop Node...")
	time.Sleep(2 * time.Second)
	<-KillRoutine
}

//Called from main.go when the command is "LIST\n"
//Show the List
func ShowList() {
	if Status {
		UpQryChan <- UpdateQuery{0, ""}
		curList := <-MemListChan
		fmt.Println("The current membership list is:")
		fmt.Print(curList)
	} else {
		fmt.Println("This server doesn't belong to a group")
	}
}

//Called from main.go when the command is "ID\n"
//Show Local ID
func ShowID() {
	if Status {
		fmt.Println("The current node ID is:")
		fmt.Print(LocalID)
	} else {
		fmt.Println("This server doesn't belong to a group")
	}
}
