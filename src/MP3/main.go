package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	Mem "./Membership"
	Config "./Config"
)

func main() {
	vmNumber := Config.GetVMNumber()
	logFile, err := os.OpenFile("MP2_"+strconv.Itoa(vmNumber)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	isIntroducer := Config.IsIntroducer()

	reader := bufio.NewReader(os.Stdin)

	for {
		var cmd string
		fmt.Println("Please type your command:")
		cmd, _ = reader.ReadString('\n')

		switch cmd {
		case "Join\n":
			log.Println("Main: Join the group")
			go Mem.RunNode(isIntroducer)
		case "Leave\n":
			log.Println("Main: Leave the group")
			go Mem.StopNode()
		case "List\n":
			log.Println("Main: Show the current Membership List")
			go Mem.ShowList()
		case "ID\n":
			log.Println("Main: Show the current Node ID")
			go Mem.ShowID()
		default:
			log.Println("Main: Don't support this command")
		}
	}
}
