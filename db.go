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
	conn := redisConnect()
	defer conn.Close()

	reply, err := conn.Do("FLUSHDB")
	if err != nil {
		log.Printf("Could not FLUSHDB Redis with error %v\n", err)
	}
	log.Printf("Successfully FLUSHDB with reply %v\n", reply)

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
	key := "item:" + strconv.Itoa(item.ID)
	reply, err := conn.Do("SET", key, b)
	if err != nil {
		panic(err)
	}

	log.Printf("SET %s %s\n", reply, key)

	key = "screenshot:item:" + strconv.Itoa(item.ID)
	reply, err = conn.Do("SET", key, item.ScreenshotAsEncodedString)
	if err != nil {
		panic(err)
	}

	log.Printf("SET %s\n", reply)
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

// FindEncodedScreenshotByID look up the screenshot encoded string by id.
func FindEncodedScreenshotByID(id int) []byte {
	conn := redisConnect()
	defer conn.Close()

	reply, err := conn.Do("GET", "screenshot:item:"+strconv.Itoa(id))
	if err != nil {
		log.Fatalf("Failed with error: %v", err)
	}

	log.Println("GET OK")
	if reply == nil {
		return nil
	}
	if val, ok := reply.([]byte); ok {
		return val
	}
	return nil
}
