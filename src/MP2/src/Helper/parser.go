package helper

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	Introducer := "cs425-fa19-
)


func getIPAddressFromID(ID string) string {
	return strings.Split(ID, "+")[0]
}

func getHostName() string {
	// Get client info(host name, ID, log file name)
	cmd := exec.Command("hostname")
	hName, _ := cmd.Output()
	hostName := string(hName)
	return hostName
}

func getVMNumber() int {
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
