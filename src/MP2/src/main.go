package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	helper "./Helper"
	node "./Node"
)

func main() {
	vmNumber := helper.GetVMNumber()
	logFile, err := os.OpenFile("MP2_"+strconv.Itoa(vmNumber)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	isIntroducer := helper.IsIntroducer()

	reader := bufio.NewReader(os.Stdin)

	for {
		var cmd string
		fmt.Println("Please type your command:")
		cmd, _ = reader.ReadString('\n')

		switch cmd {
		case "Join\n": //TODO if node is already in group??
			log.Println("Main: Join the group")
			go node.RunNode(isIntroducer)
		case "Leave\n": //TODO if node is not in group??
			log.Println("Main: Leave the group")
			go node.StopNode(true)
		case "List\n": //TODO if node hasn't joined a group??
			log.Println("Main: Show the current Membership List")
			go node.ShowList()
		case "ID\n":
			log.Println("Main: Show the current Node ID")
			go node.ShowID()
		default:
			log.Println("Main: Don't support this command")
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
