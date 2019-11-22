package main

import {
	"fmt"
	"strings"
	"io/ioutil"
	"os"
	"bufio"
	"encoding/json"
}

var MapperResult map[string][]string  = make(map[string][]string)

const JsonFileName = "webMap.json"

func parsePair(pair string) {
	res := strings.Split(pair, " ")
	if len(res) > 2 {
		fmt.Println("Data Error!")
	}
	src := res[0]
	tgt := res[1]
	
	MapperResult[tgt] = append(MapperResult[tgt], src)
}

func main(fileDir string) {
	
	data, fileErr := os.Open(fileDir)
	if fileErr != nil {
		log.Println(fileErr)
		panic(fileErr)
	}
	defer data.Close()

	var scanner = bufio.NewScanner(data)
	
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		parsePair(scanner.Text())
	}

	b, err := json.Marshal(MapperResult)
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile(webMap, b, 0644)

}