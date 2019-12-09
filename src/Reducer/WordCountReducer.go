package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)


//filename_prefix_K
func wordCount(content string) int {

	var cnt = 0

	//Find value (Todo: make sure [ will not appear at end of line)
	for idx, c := range content {
		if c == '[' {
			startIdx = idx + 1
		} else if c == ']' {
			val = content[startIdx, idx]
			fmt.Println(val)
			cnt += strconv.Atoi(val)
		}
	}

	return cnt
}


func postProcess(value string) {


}

func main {

	fileDir := os.Args[1]

	//How do we know the key?

	data, fileErr := os.Open(fileDir)
	if fileErr != nil {
		fmt.Printf("os.Open() error: Can't open file %s\n", fileDir)
		return
	}
	defer data.Close()

	var totalCnt = 0

	var scanner = bufio.NewScanner(data)

	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		totalCnt += gatherLink(scanner.Text())
	}

	res := PostProcess(MapperResult)

	fmt.Fprint(os.Stdout, res)

}