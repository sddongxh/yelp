package main

import (
	"fmt"
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
}

var g countHandler

func serveFunc(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.n++
	fmt.Fprintf(w, "g count is %v\n", g.n)
	fmt.Println("g count == ", g.n)

}

func main() {
	//http.Handle("/count", new(countHandler))
	//http.HandleFunc("/g", serveFunc)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	server := &http.Server{
		Addr:    ":8080",
		Handler: &g}
	log.Fatal(server.ListenAndServe())
	defer server.Close()

}
