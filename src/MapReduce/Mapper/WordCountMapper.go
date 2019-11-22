package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	//"encoding/json"
	//"io/ioutil"

	"./helper"
)


//TODO Exclude "," 
func Parse(cmd string) []string {
	cmd = strings.Join(strings.Fields(cmd), " ")
	return strings.Split(cmd, " ")
}

//Count each word in a word list
func countFromWordList(wordList []string, wordMap map[string]int) {
	//Iterate word list
	for _, word := range(wordList) {
		if _, ok := wordMap[word]; ok {
			//If the word exists in word map
			wordMap[word]++
		} else {
			//Not exists in word map
			wordMap[word] = 1
		}
	}
}

func main() {
	var wordMap map[string]int
	wordMap = make(map[string]int)

	//Read from arguments
	filepath := os.Args[1]
	prefix := os.Args[2]

	//Open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("os.Open() can't open file %s\n", filepath)
		return
	}
	defer file.Close()

	//Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		//Parse each line
		wordList := Parse(scanner.Text())
		
		//Count each word in the line
		countFromWordList(wordList, wordMap)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error")
	}

	fmt.Println(wordMap)
	helper.WriteWordMapToJsonFile(wordMap, prefix)
}
