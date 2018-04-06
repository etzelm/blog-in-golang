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

func about_me() {
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

	blurb := "Just a quick blurb about me and this blog"
	created := "March 31st, 2018"
	modified := "April 2nd, 2018"
	//		Start of the green well that backgrounds About Me post
	hold := "<div class=\"well\" style=\"background-color:#DFF0D8;\">" +
		//Start of the container for the face picture and education panel
		"<div class=\"container\">" +
		//Only used one row for the container
		"<div class=\"row\">" +
		//Start of the column with face picture in it
		"<div class=\"col-md-5 form-group\">" +
		//Spacing for face picture
		"<br><br>&emsp;&emsp;&emsp;&emsp;" +
		//Actual face picture
		"<img src=\"/public/face.png\" alt=\"My Face\" height=\"250\" width=\"250\"></div>" +
		//Start of the column with education panel in it
		"<div class=\"col-md-6 form-group\"><br>" +
		//Start of the education panel
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.6em\">Education:</div>" +
		//Start of the list
		"<ul class=\"list-group\">" +
		"<li class=\"list-group-item\"><h4>&emsp;<b>&emsp;Bachelors of Science in Computer Science</b></h4>" +
		"<h4>&emsp;&emsp;University of California, Santa Cruz</h4>" +
		"<h4>&emsp;&emsp;Attended from June 2015 to March 2018</h4></li>" +
		"<li class=\"list-group-item\"><h4>&emsp;&emsp;<b>Associates of Science in Computer Science</b></h4>" +
		"<h4>&emsp;&emsp;Diablo Valley College, Pleasant Hill, CA</h4>" +
		"<h4>&emsp;&emsp;Attended from January 2013 to June 2015</h4></li>" +
		//End of the container for the face picture and education panel
		"</div></div></div></div>" +
		//Start of the About Me paragraph and contact info
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.6em\">About Me & Contact Info:</div>" +
		"<ul class=\"list-group\">" +
		"<li class=\"list-group-item\">" +
		"<h4>&emsp;&emsp;By now you've hopefully guessed that my name is Mitchell " +
		"Etzel and as a recent college graduate I've decided to start this blog to " +
		"help foster and share my further eduction. I take these actions in the " +
		"hope that they might help others along a similar journey one day. The " +
		"goal is to dedicate this website to the study of distributed systems and things " +
		"related to them as well. Distributed systems are all around us and in most of the " +
		"web services we use nowadays. Considering that they have only gotten more popular over the " +
		"last twenty years, I believe they deserve a more thorough and discerning study.</h4></li>" +
		"<li class=\"list-group-item\">" +
		"<h4>&emsp;&emsp;<a style=\"color:#9C6708;\" href=\"mailto:etzelm@live.com\">" +
		"<img src=\"/public/email.png\" alt=\"Email\"  height=\"30\" width=\"45\"></a>" +
		"&emsp;&emsp;<a style=\"color:#9C6708;\" href=\"https://github.com/etzelm\">" +
		"<img src=\"/public/github.png\" alt=\"Github\"  height=\"45\" width=\"45\"></a>" +
		"&emsp;&emsp;<a style=\"color:#9C6708;\" href=\"https://www.linkedin.com/in/etzelm/\">" +
		"<img src=\"/public/linkedin.png\" alt=\"LinkedIn\"  height=\"45\" width=\"45\"></a></h4></h4></li>" +
		"</div>" +
		"</div>" +
		"</div>"

	d_input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String("0"),
			},
			"title": {
				S: aws.String("About Me"),
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
		ID:    0,
		Title: "About Me",
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
