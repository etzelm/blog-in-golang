package main

type article struct {
	ID       int    `json:"id"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Title    string `json:"title"`
	Blurb    string `json:"blurb"`
	Content  string `json:"content"`
}

func main() {
	//drop_tables()
	//about_me()
	how_its_made()
}
