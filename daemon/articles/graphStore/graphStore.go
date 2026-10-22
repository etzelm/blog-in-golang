package graphStore

import (
	"os"

	"github.com/etzelm/blog-in-golang/daemon/ddb"
	"github.com/etzelm/blog-in-golang/src/models"
)

// GraphStore function for addition/modification of GraphStore article
func GraphStore() {
	data, _ := os.ReadFile("articles/graphStore/articlePicture.html")
	ap := string(data)

	data, _ = os.ReadFile("articles/graphStore/panelPicture.html")
	pp := string(data)

	data, _ = os.ReadFile("articles/graphStore/graphStore.html")
	hh := string(data)

	item := models.Item{
		ArticlePicture: ap,
		Author:         "<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>",
		Categories:     "Distributed Systems,My Projects",
		CreatedDate:    "April 10th, 2018",
		Excerpt: "The goal of this project is to provide a REST-accessible graph storage service that " +
			"is available as a resource named gs and would listen at: <br><a style=\"color:#9C6708;\" " +
			"href=\"http://server-hostname:3000/gs\" target=\"_blank\">http://server-hostname:3000/gs</a>.",
		HTMLHold:     hh,
		ModifiedDate: "August 10th, 2019",
		PanelPicture: pp,
		PostID:       0,
		PostTitle:    "Scalable, Fault Tolerant, & Strongly Consistent Graph Store API",
		ShortTitle:   "Fault Tolerant Graph Store API",
		PostType:     "standard",
	}

	ddb.PutArticle("graphStore", item)
}
