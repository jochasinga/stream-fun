package main

import (
	"log"
	"net/http"
)

func main() {
	_ = startScrape(func() { log.Println("scrape done") })
	router := newRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
