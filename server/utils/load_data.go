package utils

import (
	"bufio"
	"encoding/json"
	"os"
)

//BusiInfo infomation of a business
type BusiInfo map[string]interface{}

//LoadJSONFile load json file to a Map array
func LoadJSONFile(filename string) ([]BusiInfo, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	reader := bufio.NewReader(jsonFile)

	var businesses []BusiInfo
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var info BusiInfo
		json.Unmarshal([]byte(line), &info)
		//fmt.Println(info)
		businesses = append(businesses, info)
	}
	return businesses, nil
}
