package infraCode

import (
	"os"

	"github.com/etzelm/blog-in-golang/daemon/ddb"
	"github.com/etzelm/blog-in-golang/src/models"
)

// infraCode function for addition/modification of infraCode article
func InfraCode() {
	data, _ := os.ReadFile("articles/infraCode/articlePicture.html")
	ap := string(data)

	data, _ = os.ReadFile("articles/infraCode/panelPicture.html")
	pp := string(data)

	data, _ = os.ReadFile("articles/infraCode/infraCode.html")
	hh := string(data)

	item := models.Item{
		ArticlePicture: ap,
		Author:         "<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>",
		Categories:     "Disciplines,Distributed Systems",
		CreatedDate:    "March 19th, 2022",
		Excerpt: "CloudFormation vs. CDK vs. Serverless Framework vs. Terraform: If you've ever had to deploy a " +
			"repeatable set of AWS Infrastructure then it's likely that you've come across these tools. Let's " +
			"jump into some of their strengths and weaknesses.",
		HTMLHold:     hh,
		ModifiedDate: "March 20th, 2022",
		PanelPicture: pp,
		PostID:       5,
		PostTitle:    "A look at Infrastructure as Code in the AWS Cloud",
		ShortTitle:   "IAC Compare and Contrast",
		PostType:     "standard",
	}

	ddb.PutArticle("infraCode", item)
}
