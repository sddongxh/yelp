package main

import (
	"fmt"
	"net/http"
)

func main() {
	var client *http.Client = &http.Client{}

	resp, err := client.Get("http://172.30.144.57:8080/g")
	//req, err := http.NewRequest("GET")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)

	req, err := http.NewRequest("GET", "http://172.30.144.57:8080/g", nil)
	client.Do(req)

}
