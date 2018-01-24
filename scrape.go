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

func startScrape() []Item {
	items := []Item{}
	filenames, err := filepath.Glob(filepath.Join(srcDir, "movies/*.mp4"))
	if err != nil {
		panic(err)
	}

	for _, filename := range filenames {
		item := Item{}
		fileinfo, err := os.Stat(filename)
		if os.IsNotExist(err) {
			log.Panic(err)
		}

		fileparts := strings.Split(fileinfo.Name(), ".")
		hyphenatedName := fileparts[:len(fileparts)-1][0]
		screenshotFullPath := filepath.Join(srcDir, "screenshots", hyphenatedName+".jpg")
		if _, err := os.Stat(screenshotFullPath); !os.IsNotExist(err) {
			item.ScreenshotURL = screenshotFullPath
		}
		almostTitle := strings.Title(strings.Join(strings.Split(hyphenatedName, "-"), " "))
		splitTitle := strings.Split(almostTitle, " ")

		for i, word := range splitTitle {
			for _, joiner := range joiners {
				if word == strings.Title(joiner) {
					splitTitle[i] = strings.ToLower(word)
				}
			}
			finalTitle := strings.Join(splitTitle, " ")
			item.Title = finalTitle
		}
		items = append(items, item)
	}
	return items
}
