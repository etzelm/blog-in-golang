package reactRealtor

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/etzelm/blog-in-golang/src/models"
	log "github.com/sirupsen/logrus"
)

// ReactRealtor function for addition/modification of ReactRealtor article
func ReactRealtor() {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(id, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Error(err)
		return
	}
	dbSvc := dynamodb.New(sess)

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

	av, _ := dynamodbattribute.MarshalMap(item)

	table := os.Getenv("ARTICLES")
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	log.Info("Putting reactRealtor into DDB")
	_, err = dbSvc.PutItem(input)

	if err != nil {
		log.Error("Got error calling PutItem:")
		log.Error(err.Error())
	}
}
