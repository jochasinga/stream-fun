package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
)

// const (
// 	day   = time.Duration(24) * time.Hour
// 	week  = time.Duration(24*7) * time.Hour
// 	month = time.Duration(24*7*30) * time.Hour
// )

var (
	currentItemID      = 1
	currentUserID      = 1
	sampleReleaseDates = [...]time.Time{
		time.Date(2018, time.April, 13, 23, 0, 0, 0, time.UTC),
		time.Date(2018, time.February, 14, 12, 0, 0, 0, time.UTC),
		time.Date(2018, time.November, 16, 15, 0, 0, 0, time.UTC),
		time.Date(2018, time.December, 6, 22, 0, 0, 0, time.UTC),
	}
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
	for i, item := range items {
		item.Watchers = rand.Intn(10000000) + 10000
		item.GrossTotal = rand.Intn(100000000) + 1000000
		item.Ratings = func() Rating {
			rating := Rating(rand.Intn(6) + 1)
			if rating > 5 {
				return 5
			}
			return rating
		}()
		// item.ReleaseStatus = ReleaseStatus(rand.Intn(3))

		// if item.ReleaseStatus > Showing {
		// 	if item.ReleaseStatus == ThisWeek {
		// 		item.Countdown = time.Duration(24*(1+rand.Intn(3))) * time.Hour
		// 	}
		// 	if item.ReleaseStatus == Upcoming {
		// 		item.Countdown = time.Duration(24*7*(1+rand.Intn(2))) * time.Hour
		// 	}
		// }
		// fmt.Println(item.Countdown)
		if i < len(sampleReleaseDates) {
			item.ReleaseDate = sampleReleaseDates[i]
		}
		CreateItem(item)
	}
}

func injectMockData(target Item, itemFunc func(Item) Item) Item {
	return itemFunc(target)
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

func redisSubscribe(ctx context.Context, channel string, conn *websocket.Conn) error {
	redisconn := redisConnect()
	defer conn.Close()

	const healthCheckPeriod = time.Minute

	psc := redis.PubSubConn{Conn: redisconn}
	if err := psc.Subscribe(channel); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				log.Printf("writing data %v\n", n.Data)
				_ = conn.WriteMessage(websocket.TextMessage, n.Data)
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()

Loop:
	for {
		select {
		case <-ticker.C:
			if err := psc.Ping(""); err != nil {
				break Loop
			}
		case <-ctx.Done():
			break Loop
		case err := <-done:
			return err
		}
	}

	return <-done
}

func redisPublish(channel, data string) {
	conn := redisConnect()
	defer conn.Close()
	conn.Do("PUBLISH", channel, data)
}
