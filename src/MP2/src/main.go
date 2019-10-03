package main

import (
	//node "MP2/src/node"
	"fmt"
	"./Helper"
	"./Node"
	// "flag"
	//"log"
	"os"
	"bufio"
)

const (
	IntroducerAddress = "fa19-cs425-g73-01.cs.illinois.edu"
)
/*
type Node struct {
	node.Sender
	node.Listener
	node.Introducer
	node.Updater
}
*/
func main() {
	// isJoinPtr := flag.Bool("join", false, "join the group")
	// isLeavePtr := flag.Bool("leave", false, "voluntarily leave the group")
	// showListPtr := flag.Bool("membership", false, "show the current membership list")
	// showIDPtr := flag.Bool("ID", false, "show self's ID")

	// flag.Parse()
	
	isIntroducer := helper.IsIntroducer()

	reader := bufio.NewReader(os.Stdin)
	
	nodeProcess:for {
		var cmd string
		fmt.Println("Please type your command:")
		cmd, _ = reader.ReadString('\n')

		switch cmd{
			case "JOIN\n":
				fmt.Println("Join the group")
				node.RunNode(isIntroducer)
				continue
			case "LEAVE\n":
				fmt.Println("Leave the group")
				node.StopNode()
				break nodeProcess
			default:
				fmt.Println("Don't support this command")
				continue
		}
	}
	/*
	fmt.Println("Start running server...")

	logFile, fileErr := os.OpenFile("MP2.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Println(fileErr)
	}
	defer logFile.Close()

	var updateChan = make(chan node.UpdateQuery)
	go node.UpdateMembershipList(updateChan)
	go node.ListenMsg()
	// switch {
	// case *isJoinPtr:

	// }
	*/
}

func NodeBehavior() {
	/**
	Check if this server is introducer (hard core to VM01)
	Execute different logic
	**/
}
