package main

import (
	"bufio"
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	srcDir  = "./assets"
	joiners = [...]string{
		"of", "and", "with", "or", "if", "is",
		"am", "are", "in", "on", "the", "a", "an",
	}
)

func base64EncodeFileToString(filename string) (string, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return "", err
	}

	info, _ := file.Stat()
	size := info.Size()
	buf := make([]byte, size)
	reader := bufio.NewReader(file)
	reader.Read(buf)

	imgbase64Str := base64.StdEncoding.EncodeToString(buf)
	return imgbase64Str, nil
}

func getScreenshotFullPath(filename string) string {
	// Remove the file extension.
	filenameParts := strings.Split(filename, ".")
	withoutExt := filenameParts[:len(filenameParts)-1][0]
	screenshotFullPath := filepath.Join(srcDir, "screenshots", withoutExt+".jpg")
	if _, err := os.Stat(screenshotFullPath); !os.IsNotExist(err) {
		return screenshotFullPath
	}
	return ""
}

func getTitleFromFilename(filename string) string {
	// Remove the file extension.
	filenameParts := strings.Split(filename, ".")
	withoutExt := filenameParts[:len(filenameParts)-1][0]
	almostTitle := strings.Title(strings.Join(strings.Split(withoutExt, "-"), " "))
	splitTitle := strings.Split(almostTitle, " ")

	// Rudimentary way to not capitalize joiner words.
	for i, word := range splitTitle {
		for _, joiner := range joiners {
			if word == strings.Title(joiner) {
				splitTitle[i] = strings.ToLower(word)
			}
		}
	}
	return strings.Join(splitTitle, " ")
}

func startScrape() []Item {
	items := []Item{}
	filenames, err := filepath.Glob(filepath.Join(srcDir, "movies/*.mp4"))
	if err != nil {
		log.Panic(err)
	}

	for _, filename := range filenames {
		fileInfo, err := os.Stat(filename)
		if os.IsNotExist(err) {
			log.Panic(err)
		}

		screenshotFullPath := getScreenshotFullPath(fileInfo.Name())
		titleText := getTitleFromFilename(fileInfo.Name())
		encodedFileStr, _ := base64EncodeFileToString(screenshotFullPath)

		item := Item{
			ItemURL:       filename,
			ScreenshotURL: screenshotFullPath,
			Title:         titleText,
			ScreenshotAsEncodedString: encodedFileStr,
		}

		// log.Printf("base64string for %s is %s", filename, encodedFileStr)
		items = append(items, item)
	}
	return items
}
