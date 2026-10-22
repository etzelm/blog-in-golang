package awsEMR

import (
	"os"

	"github.com/etzelm/blog-in-golang/daemon/ddb"
	"github.com/etzelm/blog-in-golang/src/models"
)

// AmazonEMR function for addition/modification of AmazonEMR article
func AmazonEMR() {
	data, _ := os.ReadFile("articles/awsEMR/articlePicture.html")
	ap := string(data)

	data, _ = os.ReadFile("articles/awsEMR/panelPicture.html")
	pp := string(data)

	data, _ = os.ReadFile("articles/awsEMR/awsEMR.html")
	hh := string(data)

	item := models.Item{
		ArticlePicture: ap,
		Author:         "<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>",
		Categories:     "Cloud Services,Distributed Systems",
		CreatedDate:    "November 26th, 2020",
		Excerpt: "Gave a company-wide presentation on an introduction to Amazon Web " +
			"Service's big data offering called Elastic Map Reduce. This article reviews " +
			"the content I put together for that talk.",
		HTMLHold:     hh,
		ModifiedDate: "November 26th, 2020",
		PanelPicture: pp,
		PostID:       4,
		PostTitle:    "AWS's Elastic Map Reduce Offering",
		ShortTitle:   "Intro to AWS EMR",
		PostType:     "standard",
	}

	ddb.PutArticle("awsEMR", item)
}
