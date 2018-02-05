package main

type Rating int

const (
	Eyesore Rating = iota + 1
	Bad
	Watchable
	Good
	Great
)

func (rating Rating) String() string {
	ratings := [...]string{
		"",
		"Eyesore",
		"Bad",
		"Watchable",
		"Good",
		"Great",
	}
	return ratings[rating]
}

// User represents a user of service.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Item represents an item.
type Item struct {
	ID                        int    `json:"id"`
	Filename                  string `json:"filename"`
	Title                     string `json:"title"`
	Description               string `json:"description"`
	ScreenshotURL             string `json:"screenshotUrl"`
	ScreenshotAsEncodedString string `json:"screenshotAsEncodedString"`
	ItemURL                   string `json:"itemUrl"`
	Ratings                   Rating `json:"ratings"`
}
