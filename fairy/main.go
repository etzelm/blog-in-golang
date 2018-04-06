package main

import (
	"fmt"
	"os"
)

type ItemInfo struct {
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Blurb    string `json:"blurb"`
	Content  string `json:"content"`
}

type Item struct {
	ID    int      `json:"id"`
	Title string   `json:"title"`
	Info  ItemInfo `json:"info"`
}

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
	case "42":
		create_table()
	//Execute Order 66 Meme
	case "66":
		drop_tables()
	default:
		fmt.Printf("No Input Given")
	}
}