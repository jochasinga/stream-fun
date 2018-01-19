package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

var (
	currentItemID int
	currentUserID int
)

func init() {
	items := startScrape()
	for _, item := range items {
		CreateItem(item)
	}
}

func redisConnect() redis.Conn {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	return conn
}

// CreateItem creates a new Item in the database.
func CreateItem(item Item) {
	item.ID = currentItemID
	currentItemID++

	conn := redisConnect()
	defer conn.Close()

	b, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	// Save JSON blob to Redis
	reply, err := conn.Do("SET", "item:"+strconv.Itoa(item.ID), b)
	if err != nil {
		panic(err)
	}

	log.Println("GET ", reply)

}

// FindAllItem finds all available items. Right now it's stubbed.
func FindAllItem() []Item {
	conn := redisConnect()
	defer conn.Close()

	keys, err := conn.Do("KEYS", "item:*")
	if err != nil {
		panic(err)
	}

	var items []Item
	if keys, ok := keys.([]interface{}); ok {
		for _, k := range keys {
			var item Item
			reply, err := conn.Do("GET", k.([]byte))
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(reply.([]byte), &item); err != nil {
				panic(err)
			}
			items = append(items, item)
		}
		return items
	}
	return nil
}

// FindItemByID finds a single Item based on a given ID.
func FindItemByID(id int) Item {
	var item Item

	conn := redisConnect()
	defer conn.Close()

	reply, err := conn.Do("GET", "item:"+strconv.Itoa(id))
	if err != nil {
		panic(err)
	}

	log.Println("GET OK")
	if err = json.Unmarshal(reply.([]byte), &item); err != nil {
		panic(err)
	}
	return item
}
