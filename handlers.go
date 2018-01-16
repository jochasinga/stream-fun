package main

import (
	"encoding/json"
	"net/http"

	mux "github.com/julienschmidt/httprouter"
)

// BrowseHandler handles the browse page.
func BrowseHandler(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	items := FindAll()
	if err := json.NewEncoder(w).Encode(items); err != nil {
		panic(err)
	}
}

// IndexHandler handles the index route. At the moment, it is always serving
// a local movie file "test.mp4"
func IndexHandler(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	http.ServeFile(w, r, "./test.mp4")
}
