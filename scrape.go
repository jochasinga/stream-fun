package main

import (
	"log"
	"os"
	"path/filepath"
)

var srcDir = "./assets"

func startScrape(callbackFn func()) []Item {
	filenames, err := filepath.Glob(filepath.Join(srcDir, "*.mp4"))
	if err != nil {
		panic(err)
	}

	for _, filename := range filenames {
		log.Println(filename)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Printf("file doesn't exist")
		} else {
			log.Printf("file exist")
		}
	}
	callbackFn()
	return []Item{}
}
