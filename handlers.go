package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
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
	payload := make(response)
	payload.wrap(imageEncodedStr)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Panic(err)
	}
}

// WatchHandler handles serving an individual item.
func WatchHandler(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	id := getItemID(ps)
	itemURL := FindItemByID(id).ItemURL
	http.ServeFile(w, r, itemURL)
}

// CountdownHandler handles serving a countdown timer websocket.
func CountdownHandler(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	itemID := getItemID(ps)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Origin")
	}
	serveWs(w, r, itemID)
}

func getItemID(ps mux.Params) int {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		panic(err)
	}
	return id
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(w http.ResponseWriter, r *http.Request, id int) {
	log.Println("serving Ws")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	item := FindItemByID(id)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		redisconn := redisConnect()
		defer redisconn.Close()

		for _ = range ticker.C {
			untilRelease := time.Until(item.ReleaseDate)
			redisconn.Do("PUBLISH", "countdown:"+strconv.Itoa(id), strconv.Itoa(int(untilRelease)))
			// redisconn.Do("PUBLISH", "countdown:"+strconv.Itoa(id), strconv.Itoa(int(item.Countdown)))
			// item.Countdown--

		}
	}()

	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer wsconn.Close()

	go func() {
		countdownChannel := "countdown:" + strconv.Itoa(id)
		redisSubscribe(ctx, countdownChannel, wsconn)
	}()

	for {
		messageType, p, err := wsconn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("recv: %s", p)
		if err := wsconn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	log.Printf("Received: %v\n", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	log.Println("body:", string(body))
}
