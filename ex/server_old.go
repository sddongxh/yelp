package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
	fmt.Println("this is handler")
	fmt.Println(r.URL, r.Method)
	//var s []byte
	//r.Body.Read(s)
	//fmt.Println(s)
	reqBody, _ := ioutil.ReadAll(r.Body)
	//fmt.Fprintf(w, "%+v", string(reqBody))
	fmt.Println(string(reqBody))
}

var g countHandler

func serveFunc(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.n++
	fmt.Fprintf(w, "g count is %v\n", g.n)
	fmt.Println("g count == ", g.n)
	fmt.Println(r.URL)
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func main() {
	//http.Handle("/count", new(countHandler))
	//http.HandleFunc("/g", serveFunc)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	//http.Handle()
	server := &http.Server{
		Addr:    ":8080",
		Handler: &g}
	log.Fatal(server.ListenAndServe())
	defer server.Close()

}
