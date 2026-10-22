package googleSRE

import (
	"os"

	"github.com/etzelm/blog-in-golang/daemon/ddb"
	"github.com/etzelm/blog-in-golang/src/models"
)

// GoogleSRE function for addition/modification of GoogleSRE article
func GoogleSRE() {
	data, _ := os.ReadFile("articles/googleSRE/articlePicture.html")
	ap := string(data)

	data, _ = os.ReadFile("articles/googleSRE/panelPicture.html")
	pp := string(data)

	data, _ = os.ReadFile("articles/googleSRE/googleSRE.html")
	hh := string(data)

	item := models.Item{
		ArticlePicture: ap,
		Author:         "<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>",
		Categories:     "Disciplines,Distributed Systems",
		CreatedDate:    "April 11th, 2018",
		Excerpt: "If you've ever smelled bad SCRUM you're going to want to learn how to internalize " +
			"these important lessons I've identified in the O'Reilly book, <i><a style=\"color:#9C6708;\" " +
			"href=\"https://landing.google.com/sre/book.html\" target=\"_blank\">Site Reliability Engineering: " +
			"How Google Runs Production Systems</a></i>.",
		HTMLHold:     hh,
		ModifiedDate: "March 18th, 2022",
		PanelPicture: pp,
		PostID:       2,
		PostTitle:    "How To Internalize Site Reliability Engineering's Top 5 Golden Lessons",
		ShortTitle:   "SRE Internalization",
		PostType:     "standard",
	}

	ddb.PutArticle("googleSRE", item)
}
