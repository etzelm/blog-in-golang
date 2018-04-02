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
	switch args := os.Args[1:]; args[0] {
	case "0":
		drop_tables()
	case "1":
		about_me()
	case "2":
		how_its_made()
	default:
		fmt.Printf("No Input Given")
	}
}
