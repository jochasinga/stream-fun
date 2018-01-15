package main

import (
	"net/http"

	mux "github.com/julienschmidt/httprouter"
)

// IndexHandler handles the index route. At the moment, it is always serving
// a local movie file "test.mp4"
func IndexHandler(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	http.ServeFile(w, r, "./test.mp4")
}
