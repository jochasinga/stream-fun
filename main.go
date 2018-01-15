package main

import (
	"log"
	"net/http"
)

func main() {
	router := newRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
