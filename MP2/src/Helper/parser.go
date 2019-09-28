package parser

import (
	"bufio"
	"fmt"
	"strconv"
	"log"
	"io"
	"os"
	"os/exec"
	"strings"
)

func getClientInfo() int {
	// Get client info(host name, ID, log file name)
	cmd := exec.Command("hostname")
	hName,_ := cmd.Output()
	hostName := string(hName)
	var machineNO int
	// var machineName string
	if hostName[15] == '0'{
		machineNO, _  = strconv.Atoi(hostName[16:17])
	} else {
		machineNO, _ = strconv.Atoi(hostName[15:17])
	}
	return machineNO
}