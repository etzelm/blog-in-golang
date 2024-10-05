package googleSRE

import (
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/etzelm/blog-in-golang/src/models"
	log "github.com/sirupsen/logrus"
)

// GoogleSRE function for addition/modification of GoogleSRE article
func GoogleSRE() {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(id, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return
	}
	dbSvc := dynamodb.New(sess)

	data, _ := ioutil.ReadFile("articles/googleSRE/articlePicture.html")
	ap := string(data)

	data, _ = ioutil.ReadFile("articles/googleSRE/panelPicture.html")
	pp := string(data)

	data, _ = ioutil.ReadFile("articles/googleSRE/googleSRE.html")
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

	av, _ := dynamodbattribute.MarshalMap(item)

	table := os.Getenv("ARTICLES")
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	log.Info("Putting googleSRE into DDB")
	_, err = dbSvc.PutItem(input)

	if err != nil {
		log.Error("Got error calling PutItem:")
		log.Error(err.Error())
		os.Exit(1)
	}
}
