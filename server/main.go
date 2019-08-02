package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"yelp/server/utils"
)

//BusiInfo alias name
type BusiInfo = utils.BusiInfo

//QuadTreeNode alias name
type QuadTreeNode = utils.QuadTreeNode

//UnitDB alias name
type UnitDB = utils.UnitDB

//GeoHashDB alias name
type GeoHashDB = utils.GeoHashDB

//Cell range of a cell
type Cell = utils.Cell

//global variables
var dbTree *QuadTreeNode
var db GeoHashDB

func initServer(filename string) {
	data, err := utils.LoadJSONFile(filename)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Data size:", len(data))
	db = *utils.NewGeoHashDB()
	dbTree = utils.BuildQuadTreeFromData(data, 100, db)
	fmt.Println("Area infomation:", *dbTree)
	fmt.Println("GEO hash information:", len(db), "units")
	fmt.Println("Server has been successfully initiated")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

type yelpHandler struct {
	mu sync.Mutex
	id int // request id
}

func formatSearchResult(info []BusiInfo, dists []float64) string {
	var res string
	res += fmt.Sprintf("Stars Name Distance Address")
	for i, b := range info {
		res += fmt.Sprintf("<p>%v | %v | %v | %v</p>", b["stars"], b["name"], int(dists[i]), b["address"])
	}
	return res
}

func serveGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	u, _ := url.Parse(r.RequestURI)
	query := u.Query()
	lon, _ := strconv.ParseFloat(query["lon"][0], 64)
	lat, _ := strconv.ParseFloat(query["lat"][0], 64)
	dist, _ := strconv.ParseFloat(query["r"][0], 64)
	searchTerms := query["find_desc"][0]
	fmt.Println(lon, lat, dist, searchTerms)
	dbIndices := utils.FindAllDBInRange(dbTree, lat, lon, dist)
	businesses, distances := utils.FindALlBusiInRange(db, dbIndices, lat, lon, dist, searchTerms, "*")
	ret := formatSearchResult(businesses, distances)
	fmt.Fprintf(w, ret)
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
	filename := "/home/xihua/go/src/yelp/data/vegas.json"
	initServer(filename)
	server := &http.Server{
		Addr:    ":8080",
		Handler: new(yelpHandler)}
	log.Fatal(server.ListenAndServe())
	defer server.Close()
}
