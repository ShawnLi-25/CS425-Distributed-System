package main

import (
	"exec"
	"fmt"
)

func main() {
	temp = "./webTest"
	cmd := exec.Command("./WebMapper", temp)
	res, _ := cmd.Output()
	s = string(res)
	fmt.Println(s)
}
