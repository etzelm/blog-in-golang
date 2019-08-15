package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

func graphStore() {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	log.Info("id: ", id)
	log.Info("key: ", key)
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

	data, err := ioutil.ReadFile("graphStore/articlePicture.html")
	ap := string(data)

	data, err = ioutil.ReadFile("graphStore/panelPicture.html")
	pp := string(data)

	data, err = ioutil.ReadFile("graphStore/graphStore.html")
	hh := string(data)

	item := Item{
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
		PostID:       2,
		PostTitle:    "Scalable, Fault Tolerant, & Strongly Consistent Graph Store API",
		ShortTitle:   "Fault Tolerant Graph Store API",
		PostType:     "standard",
	}

	av, err := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("New-Articles"),
	}

	_, err = dbSvc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
