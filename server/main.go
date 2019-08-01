package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

type yelpHandler struct {
	mu sync.Mutex
	id int // request id
}

func serveGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")
	fmt.Fprintf(w, "12,3,12134412")
	fmt.Println(r.RequestURI)
	u, _ := url.Parse(r.RequestURI)
	query := u.Query()
	fmt.Printf("find_desc = %v\n", query["find_desc"])
	fmt.Printf("l = %v\n", query["l"])
	fmt.Printf("r = %v\n", query["r"])
	fmt.Printf("ranking = %v\n", query["ranking"])
}

func servePost(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	req := string(reqBody)
	fmt.Println(req)
	fmt.Fprintf(w, req)
}

func (h *yelpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	h.mu.Lock()
	defer h.mu.Unlock()
	h.id++
	//	fmt.Fprintf(w, "Your request (%d) has been received\n", h.id)
	switch r.Method {
	case "POST":
		fmt.Println("POST received")
		servePost(w, r)
	case "GET":
		fmt.Println("GET received")
		serveGet(w, r)
	default:
		fmt.Printf("Other type of request received: %s\n", r.Method)
	}

}

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: new(yelpHandler)}
	log.Fatal(server.ListenAndServe())
	defer server.Close()
}
