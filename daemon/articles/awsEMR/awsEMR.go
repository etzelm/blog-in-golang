package awsEMR

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

// AmazonEMR function for addition/modification of AmazonEMR article
func AmazonEMR() {
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

	av, _ := dynamodbattribute.MarshalMap(item)

	table := os.Getenv("ARTICLES")
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	log.Info("Putting awsEMR into DDB")
	_, err = dbSvc.PutItem(input)

	if err != nil {
		log.Error("Got error calling PutItem:")
		log.Error(err.Error())
	}
}
