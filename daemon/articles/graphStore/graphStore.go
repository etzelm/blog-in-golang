package graphStore

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

// GraphStore function for addition/modification of GraphStore article
func GraphStore() {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(id, key, "")

	session, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Error(err)
		return
	}
	dbSvc := dynamodb.New(session)

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

	av, _ := dynamodbattribute.MarshalMap(item)

	table := os.Getenv("ARTICLES")
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	log.Info("Putting graphStore into DDB")
	_, err = dbSvc.PutItem(input)

	if err != nil {
		log.Error("Got error calling PutItem:")
		log.Error(err.Error())
	}
}
