package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

func google_SRE() {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	log.Info("id: ", id)
	log.Info("key: ", key)
	var my_credentials = credentials.NewStaticCredentials(id, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: my_credentials,
		Region:      aws.String("us-west-1"),
		Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return
	}
	dbSvc := dynamodb.New(sess)

	blurb := "Using the popular O'Reilly book, I try to gain a better understanding of what SRE really means"
	created := "April 11th, 2018"
	modified := "April 12th, 2018"
	hold := "<div class=\"well\" style=\"background-color:#DFF0D8;\">" +
		//Start of the About Me paragraph and contact info
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.6em\">" +
		"A Review & Implementation of Google's Site Reliability Engineering:</div>" +
		"<ul class=\"list-group\">" +
		"<li class=\"list-group-item\">" +
		"<h4>&emsp;&emsp;Paying homage to who many consider to be at the forefront of distributed systems, I have deciced " +
		"to review and implement some of the concepts put forth in the O'Reilly book, <i>" +
		"<a style=\"color:#9C6708;\" href=\"https://landing.google.com/sre/book.html\" target=\"_blank\">" +
		"Site Reliability Engineering: How Google Runs Production Systems</a></i>. This work will be done concurrently " +
		"with my development of a personal distributed system for graph storage to help foster both projects and I have " +
		"attached a link to that source code at the bottom of this introduction. " +
		"" +
		"</h4></li>" +
		"<li class=\"list-group-item\" style=\"text-align: center;\">" +
		"<h4>Source code for this project can be located here:&emsp;&emsp;" +
		"<a href=\"https://github.com/etzelm/consistent-graph-store-api\" target=\"_blank\">" +
		"<img src=\"/public/github.png\" alt=\"Github\"  height=\"45\" width=\"45\"></a></h4>" +
		"</div>" +
		//Start of the container for

		"</div>"

	d_input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String("2"),
			},
			"title": {
				S: aws.String("A Review & Implementation of Google's Site Reliability Engineering"),
			},
		},
		TableName: aws.String("Articles"),
	}

	_, err = dbSvc.DeleteItem(d_input)

	if err != nil {
		fmt.Println("Got error calling DeleteItem")
		fmt.Println(err.Error())
		return
	}

	info := ItemInfo{
		Created:  created,
		Modified: modified,
		Blurb:    blurb,
		Content:  hold,
	}

	item := Item{
		ID:    2,
		Title: "A Review & Implementation of Google's Site Reliability Engineering",
		Info:  info,
	}

	av, err := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Articles"),
	}

	_, err = dbSvc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	result, err := dbSvc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Tables:")
	for _, table := range result.TableNames {
		log.Println(*table)
	}
}
