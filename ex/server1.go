package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var articles []Article

type countHandler struct {
	mu sync.Mutex // guards n
	n  int
	j  []byte
}

var database []string

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
	fmt.Println("this is handler")
	//json.NewEncoder(w).Encode(articles)
	//bz := ""
	//fmt.Fprintf(w, bz)
	data := database[h.n%len(database)]
	fmt.Fprintf(w, string(data))
}

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
	for i := 0; i < 1000; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		//database = append(database, line)
		//fmt.Println(line)
		//zip := line[]
		json.Unmarshal([]byte(line), &business)
		//zip := business["postal_code"]
		//fmt.Println(zip)
		data, _ := json.MarshalIndent(business, "", "   ")
		database = append(database, string(data))
	}
	json, err := json.Marshal(business)
	//data = json
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))

	// articles = []Article{
	// 	Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
	// 	Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	// }
	http.Handle("/", new(countHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
