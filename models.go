package main

import (
	"time"
)

var item Item

// Item represents an item.
type Item struct {
	ID        int       `json:"id"`
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
}
