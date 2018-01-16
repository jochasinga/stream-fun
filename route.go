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
	// Only an index page for now.
	Route{"Index", "GET", "/", IndexHandler},
	Route{"Browse", "GET", "/browse", BrowseHandler},
}
