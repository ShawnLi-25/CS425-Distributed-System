package main

import (
	"log"
	"os/exec"
)

func main() {
	temp := "./webTest"
	cmd := exec.Command("./WebMapper", temp)
	res, _ := cmd.Output()
	s := string(res)
	log.Println(s)
	//fmt.Println(s)
}
