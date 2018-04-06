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

func how_its_made() {
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

	blurb := "A reflection on what it took to make this website and its content, plus maybe some updates along the way"
	created := "April 1st, 2018"
	modified := "April 2nd, 2018"
	hold := "<h3>Who Would Want to Write a Blog in Go?</h3>" +
		"<h4>&emsp;&emsp;Hopefully you've guessed by now that my name is Mitchell " +
		"Etzel and as a recent college graduate I've decided to start this blog to " +
		"help foster and share my further eduction. I take these actions in the " +
		"hope that they might help others along a similar journey one day. The " +
		"goal is to dedicate this website to distributed systems and things " +
		"related to them as well. Although I make promises toward the fact that " +
		"I will probably get a little bit distracted with side projects along the " +
		"way, here and there.</h4>" +
		"</div>"

	d_input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String("0"),
			},
			"title": {
				S: aws.String("Who Would Want to Write a Blog in Golang Anyways"),
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
		ID:    1,
		Title: "Who Would Want to Write a Blog in Golang Anyways",
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