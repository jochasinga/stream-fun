package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	srcDir  = "./assets"
	joiners = [...]string{"of", "and", "with", "or", "if", "is", "am", "are"}
)

func startScrape(callbackFn func()) []Item {
	filenames, err := filepath.Glob(filepath.Join(srcDir, "*.mp4"))
	if err != nil {
		panic(err)
	}

	for _, filename := range filenames {
		fileinfo, err := os.Stat(filename)
		if os.IsNotExist(err) {
			log.Panic(err)
		}

		fileparts := strings.Split(fileinfo.Name(), ".")
		hyphenatedName := fileparts[:len(fileparts)-1][0]
		almostTitle := strings.Title(strings.Join(strings.Split(hyphenatedName, "-"), " "))
		splitAgain := strings.Split(almostTitle, " ")

		for i, word := range splitAgain {
			for _, joiner := range joiners {
				if word == strings.Title(joiner) {
					splitAgain[i] = strings.ToLower(word)
				}
			}
		}

		finalTitle := strings.Join(splitAgain, " ")
		log.Println(finalTitle)
	}
	callbackFn()
	return []Item{}
}
