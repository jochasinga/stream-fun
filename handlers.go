package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	mux "github.com/julienschmidt/httprouter"
)

type response map[string]interface{}

func (res response) wrap(any interface{}) {
	res["data"] = (interface{})(any)
}

// BrowseHandler handles the browse page.
func BrowseHandler(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	items := FindAllItem()
	payload := make(response)
	payload.wrap(items)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		panic(err)
	}
}

// ItemHandler handles finding an Item by ID.
func ItemHandler(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	id := getItemID(ps)
	item := FindItemByID(id)
	payload := make(response)
	payload.wrap(item)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		panic(err)
	}
}

// ScreenshotHandler serves a screenshot encoded string
func ScreenshotHandler(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	id := getItemID(ps)
	imageEncodedStr := FindEncodedScreenshotByID(id)
	if _, err := w.Write(imageEncodedStr); err != nil {
		log.Fatalf("Failed with error: %v", err)
	}
}

// WatchHandler handles serving an individual item.
func WatchHandler(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	id := getItemID(ps)
	itemURL := FindItemByID(id).ItemURL
	http.ServeFile(w, r, itemURL)
}

func getItemID(ps mux.Params) int {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		panic(err)
	}
	return id
}
