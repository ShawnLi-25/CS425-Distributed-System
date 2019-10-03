package helper

import (
	//"bufio"
	"fmt"
	//"io"
	//"log"
	//"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	IntroducerAddress = "fa19-cs425-g73-01.cs.illinois.edu"
)


func GetIPAddressFromID(ID string) string {
	return strings.Split(ID, "+")[0]
}

func GetHostName() string {
	// Get client info(host name, ID, log file name)
	cmd := exec.Command("hostname")
	hName, _ := cmd.Output()
	hostName := string(hName)
	return hostName
}

func GetVMNumber() int {
	// Get client info(host name, ID, log file name)
	cmd := exec.Command("hostname")
	hName, _ := cmd.Output()
	hostName := string(hName)
	var machineNO int
	// var machineName string
	if hostName[15] == '0' {
		machineNO, _ = strconv.Atoi(hostName[16:17])
	} else {
		machineNO, _ = strconv.Atoi(hostName[15:17])
	}
	return machineNO
}

func IsIntroducer() bool {
	hostName := GetHostName()
	return hostName == IntroducerAddress
}
