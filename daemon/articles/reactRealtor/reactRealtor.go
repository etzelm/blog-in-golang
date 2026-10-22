package reactRealtor

import (
	"os"

	"github.com/etzelm/blog-in-golang/daemon/ddb"
	"github.com/etzelm/blog-in-golang/src/models"
)

// ReactRealtor function for addition/modification of ReactRealtor article
func ReactRealtor() {
	data, _ := os.ReadFile("articles/reactRealtor/articlePicture.html")
	ap := string(data)

	data, _ = os.ReadFile("articles/reactRealtor/panelPicture.html")
	pp := string(data)

	data, _ = os.ReadFile("articles/reactRealtor/reactRealtor.html")
	hh := string(data)

	item := models.Item{
		ArticlePicture: ap,
		Author:         "<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>",
		Categories:     "Disciplines,Frontend Development,My Projects",
		CreatedDate:    "May 17th, 2020",
		Excerpt: "I recently had the opportunity to explore the combined " +
			"capabilities of the Go, Gin, and React libraries for an interview " +
			"assessment. This post is about that journey.",
		HTMLHold:     hh,
		ModifiedDate: "May 20th, 2020",
		PanelPicture: pp,
		PostID:       3,
		PostTitle:    "Go & React: A 1, 2 Punch Combo",
		ShortTitle:   "Go & React",
		PostType:     "standard",
	}

	ddb.PutArticle("reactRealtor", item)
}
