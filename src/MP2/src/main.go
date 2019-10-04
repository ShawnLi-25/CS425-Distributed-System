package main

import (
	"fmt"
	"./Helper"
	"./Node"
	// "flag"
	//"log"
	"os"
	"bufio"
)


func main() {
	isIntroducer := helper.IsIntroducer()

	reader := bufio.NewReader(os.Stdin)
	
	nodeProcess:for {
		var cmd string
		fmt.Println("Please type your command:")
		cmd, _ = reader.ReadString('\n')

		switch cmd{
			case "JOIN\n"://TODO if node is already in group??
				fmt.Println("Join the group")
				go node.RunNode(isIntroducer)
				continue
			case "LEAVE\n"://TODO if node is not in group??
				fmt.Println("Leave the group")
				go node.StopNode()
				break nodeProcess
			default:
				fmt.Println("Don't support this command")
				continue
		}
	}
	/**
	logFile, fileErr := os.OpenFile("MP2.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Println(fileErr)
	}
	defer logFile.Close()
	var updateChan = make(chan node.UpdateQuery)
	go node.UpdateMembershipList(updateChan)
	go node.ListenMsg()
	**/
}

