package main

import (
	"fmt"
	"os"
)

type article struct {
	ID       int    `json:"id"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Title    string `json:"title"`
	Blurb    string `json:"blurb"`
	Content  string `json:"content"`
}

func main() {
	//Switch Case Depends on Article ID Number
	switch args := os.Args[1:]; args[0] {
	case "0":
		about_me()
	case "1":
		how_its_made()
	//Execute Order 66 Meme
	case "66":
		drop_tables()
	default:
		fmt.Printf("No Input Given")
	}
}
