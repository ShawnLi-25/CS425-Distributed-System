package helper

import(
	"encoding/json"
	"os"
)

func WriteWordMapToJsonFile(mymap map[string]int, prefix string) error {
	for word, count := range(mymap) {
		filebyte, _ := json.MarshalIndent(map[string]int{word:count}, "", " ")

		file, err := os.OpenFile(prefix + "_" + word, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()

		if err != nil {
			return err
		}

		if _, err := file.Write(filebyte); err != nil{
			return err
		}
	}

	return nil
}

