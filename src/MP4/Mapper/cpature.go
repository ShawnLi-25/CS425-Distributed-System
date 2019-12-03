package main

import (
	"log"
	"os"
	"os/exec"
)

func parseMapRes(res []byte, prefix string) error {
	s := string(res)

	isKey := true

	var key, val []byte

	for i := 0; i < len(s); i++ {
		if isKey {
			if s[i] == ':' {
				isKey = false
			} else {
				key = append(key, s[i])
			}
		} else {
			if s[i] == '\n' {
				WriteFile(key, val, prefix)
			} else {
				val = append(val, s[i])
			}

		}

	}
	// var reader = strings.NewReader(s)

	// for reader.Read() {
	// 	if
	// }

	return nil
}

func WriteFile(key []byte, val []byte, prefix string) err {
	fileName := prefix + string(key)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)

	n, err := file.Write(val)
	if err != nil || n <= 0 {
		return err
	}

	return nil

}

func main() {
	temp := "./webTest"
	cmd := exec.Command("./WebMapper", temp)
	res, _ := cmd.Output()
	s := string(res)
	log.Println(s)
	//fmt.Println(s)
}
