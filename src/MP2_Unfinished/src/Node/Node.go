package node

import (
	"fmt"
	"time"

	//"log"

	msg "../Helper"
)

var curNode *Node = CreateNewNode()

// var UpQryChan chan UpdateQuery = make(chan UpdateQuery)
// var MemListChan chan []string = make(chan []string)
var KillHBListener chan struct{} = make(chan struct{})
var KillHBSender chan struct{} = make(chan struct{})
var KillMsgListener chan struct{} = make(chan struct{})

// var KillUpdater chan struct{} = make(chan struct{})
var KillIntroducer chan struct{} = make(chan struct{})
var KillHBTimer chan struct{} = make(chan struct{})

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

	// go curNode.Updater.UpdateMembershipList()
	go curNode.Listener.RunMSGListener()
	if !isIntroducer {
		//Non-intro send JoinMsg to Introducer
		curNode.Sender.SendJoin()
	} else {
		//Introducer receive JoinMsg from non-intro
		go curNode.Introducer.NodeHandleJoin()
	}

	go curNode.Sender.SendHeartbeat()
	time.Sleep(time.Second)
	go curNode.Listener.RunHBListener()

}

//Called from main.go when the command is "LEAVE\n"
//Delete the Node
func StopNode() {
	// if byLocal {
	// 	curNode.Sender.SendLeave(true)
	// } else {
	// 	curNode.Sender.SendLeave(false)
	// }

	fmt.Println("1")
	curNode.Sender.SendLeave()
	KillMsgListener <- struct{}{}

	Status = false
	fmt.Println("1")

	KillHBListener <- struct{}{}

	KillHBSender <- struct{}{}

	//When Leave, Clear all elements
	MembershipList = make([]string)
	MemHBMap = make(map[string]time.Time)

	if msg.IsIntroducer() {
		KillIntroducer <- struct{}{}
	}
	fmt.Println("Node: Stop Node...")
	// time.Sleep(3 * time.Second)
	// <-KillHBListener
}

//Called from main.go when the command is "LIST\n"
//Show the List
func ShowList() {
	if Status {
		// MembershipList
		// UpQryChan <- UpdateQuery{0, ""}
		// curList := <-MemListChan
		fmt.Println("The current membership list is:")
		fmt.Print(MembershipList, "\n")
		// fmt.Println()
	} else {
		fmt.Println("This server doesn't belong to a group")
	}
	return
}

//Called from main.go when the command is "ID\n"
//Show Local ID
func ShowID() {
	if Status {
		fmt.Println("The current node ID is:")
		fmt.Print(LocalID, "\n")
	} else {
		fmt.Println("This server doesn't belong to a group")
	}
	return
}

// func ShowMonitoringList() {
// 	if Status {
// 		fmt.Println("The current MonitorList is:")
// 		fmt.Print(LocalID, "\n")
// 	} else {
// 		fmt.Println("This server doesn't belong to a group")
// 	}
// 	return
// }
