package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vanng822/go-solr/solr"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Article models the result for the /getAll/ route
type Article struct {
	ID               int
	Title            string
	Subtitle         string
	ShortDescription string
	BuyNowPrice      float32
	CurrentBidPrice  float32
	URL              string
	MainImageURL     string
}

var solrServer *solr.SolrInterface

// During init we populate Solr with dummy data
func init() {
	var err error
	solrServer, err = solr.NewSolrInterface("http://localhost:8983/solr/", "searchAws")
	if err != nil {
		log.Fatal(err)
	}
	_, errs := solrServer.DeleteAll()
	if errs != nil {
		log.Fatal(err)
	}
	docs := make([]solr.Document, 1)
	doc1 := make(solr.Document)
	doc1["ID"] = 1
	doc1["Title"] = "Test title"
	doc1["Subtitle"] = "Test subtitle"
	docs = append(docs, doc1)
	response, err := solrServer.Add(docs, 1, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err = solrServer.Commit()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Commit response: %s", strconv.FormatBool(response.Success))
	for key, value := range response.Result {
		fmt.Println(key, "=", value)
	}
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/getAll", GetAllHandler).Methods("GET")
	http.Handle("/", r)
	log.Printf("Server started and listening on port %d.", 3232)
	log.Fatal(http.ListenAndServe(":3232", nil))
}

// GetAllHandler handles requests
func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")

	query := solr.NewQuery()
	query.Q("Title:Test title")
	s := solrServer.Search(query)
	resp, _ := s.Result(nil)
	fmt.Println(resp.Results.Docs)
	response, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
	duration := time.Since(start)
	log.Printf("\t%s\t%s",
		r.RequestURI,
		duration)
	return
}
