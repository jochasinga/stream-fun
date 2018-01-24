package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	mux "github.com/julienschmidt/httprouter"
)

// BrowseHandler handles the browse page.
func BrowseHandler(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	items := FindAllItem()
	if err := json.NewEncoder(w).Encode(items); err != nil {
		panic(err)
	}
}

// ItemHandler handles finding an Item by ID.
func ItemHandler(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	id := getItemID(ps)
	item := FindItemByID(id)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		panic(err)
	}
}

// WatchHandler handles serving an individual item.
func WatchHandler(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	id := getItemID(ps)
	itemURL := FindItemByID(id).ItemURL
	// http.ServeFile(w, r, ".assets/test.mp4")
	http.ServeFile(w, r, itemURL)
}

func getItemID(ps mux.Params) int {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		panic(err)
	}
	return id
}
