package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	jsonFile, err := os.Open("/home/xihua/code/datasets/yelp/dataset/business.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	reader := bufio.NewReader(jsonFile)
	//var row []byte
	//hash := make(map[int]int)
	var business map[string]interface{}
	for i := 0; i < 2; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(line)
		//zip := line[]
		json.Unmarshal([]byte(line), &business)
		zip := business["postal_code"]
		fmt.Println(zip)
	}

}
