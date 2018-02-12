package main

import (
	mux "github.com/julienschmidt/httprouter"
)

// Route represents a URL route data.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handle  mux.Handle
}

// Routes represents many routes
type Routes []Route

func newRouter() *mux.Router {
	router := mux.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Pattern, route.Handle)
	}
	return router
}

var routes = Routes{
	// Route{"Index", "GET", "/", IndexHandler},
	Route{"Browse", "GET", "/browse", BrowseHandler},
	Route{"Item", "GET", "/item/:id", ItemHandler},
	Route{"Watch", "GET", "/watch/:id", WatchHandler},
	Route{"Screenshot", "GET", "/screenshot/:id", ScreenshotHandler},
	Route{"Countdown", "GET", "/countdown/:id", CountdownHandler},
}
