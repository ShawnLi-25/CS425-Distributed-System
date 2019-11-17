package main

import {
	"fmt"
	"io"
	"os"
	"bufio"
}

var MapperResult map[string]int = make(map[string]int)

func doMap(fileName string) [string]string {
	
	data, fileErr := os.Open("SDFS/" + fileName)
	if fileErr != nil {
		log.Println(fileErr)
		panic(fileErr)
	}
	defer data.Close()


	dataScanner := bufio.NewScanner(data)
	// buffer, _, err := ReadLine(dataScanner)

	for dataScanner.Scan() {
		fmt.Println(dataScanner.Text())
		
	}
	

	



}

func ()

func Parse(cmd string) []string {
	cmd = cmd[:len(cmd)-1]
	cmd = strings.Join(strings.Fields(cmd), " ")
	return strings.Split(cmd, " ")
}

func main() {
	
	reader := bufio.NewReader(os.Stdin)

	arg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Reading Error: ", err.Error())
		return
	}
	
	//arg should be file name
	doMap(arg)

	// parsedCmd = Parse(cmd)

	
	

}