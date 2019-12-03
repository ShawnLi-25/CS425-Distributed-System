package main

import (
	"fmt"
	"os/exec"
)

func main() {
	temp := "./webTest"
	cmd := exec.Command("./WebMapper", temp)
	res, _ := cmd.Output()
	s := string(res)
	fmt.Println(s)
}
