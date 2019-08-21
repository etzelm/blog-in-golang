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

func googleSRE() {
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

	data, err := ioutil.ReadFile("googleSRE/articlePicture.html")
	ap := string(data)

	data, err = ioutil.ReadFile("googleSRE/panelPicture.html")
	pp := string(data)

	data, err = ioutil.ReadFile("googleSRE/googleSRE.html")
	hh := string(data)

	item := Item{
		ArticlePicture: ap,
		Author:         "<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>",
		Categories:     "Disciplines,Distributed Systems",
		CreatedDate:    "April 11th, 2018",
		Excerpt: "If you've ever smelled bad SCRUM you're going to want to learn how to internalize " +
			"these important lessons I've identified in the O'Reilly book, <i><a style=\"color:#9C6708;\" " +
			"href=\"https://landing.google.com/sre/book.html\" target=\"_blank\">Site Reliability Engineering: " +
			"How Google Runs Production Systems</a></i>.",
		HTMLHold:     hh,
		ModifiedDate: "August 20th, 2019",
		PanelPicture: pp,
		PostID:       0,
		PostTitle:    "How To Internalize Site Reliability Engineering's Top 3 Golden Lessons",
		ShortTitle:   "SRE Internalization - Mitchell Etzel",
		PostType:     "standard",
	}

	av, err := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Live-Articles"),
	}

	_, err = dbSvc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
