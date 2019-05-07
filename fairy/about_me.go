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
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return
	}
	dbSvc := dynamodb.New(sess)

	blurb := "Just a quick blurb about me and this blog"
	created := "March 31st, 2018"
	modified := "May 6th, 2019"
	//		Start of the green well that backgrounds About Me post
	hold := "<div class=\"well\" style=\"background-color:#DFF0D8;\">" +
		//Start of the container for the face picture and education panel
		"<div class=\"container-fluid\">" +
		//Only used one row for the container
		"<div class=\"row\">" +
		//Start of the column with face picture in it
		"<div class=\"col-md-5 form-group\" style=\"text-align: center;\">" +
		//Spacing for face picture
		"<br>" +
		//Actual face picture
		"<img src=\"/public/face.png\" alt=\"My Face\" height=\"250\" width=\"250\"></div>" +
		//Start of the column with education panel in it
		"<div class=\"col-md-6 form-group\"><br>" +
		//Start of the education panel
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">Education:</div>" +
		//Start of the list
		"<ul class=\"list-group\" style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\"><b>Bachelors of Science in Computer Science</b><br>" +
		"UCSC, Santa Cruz, CA<br>" +
		"Attended: June '15 to March '18</li>" +
		"<li class=\"list-group-item\"><b>Associates of Science in Computer Science</b><br>" +
		"DVC, Pleasant Hill, CA<br>" +
		"Attended: Janurary '13 to June '15</li>" +
		//End of the container for the face picture and education panel
		"</div></div></div></div>" +
		//Start of the About Me paragraph and contact info
		"<div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\" style=\"color:#A619FF;font-size: 1.17em\">About Me & Contact Info:</div>" +
		"<ul class=\"list-group\" style=\"font-size: 1.07em;\">" +
		"<li class=\"list-group-item\">" +
		"&emsp;&emsp;By now you've hopefully guessed that my name is Mitchell Etzel and as a college graduate, " +
		"as well as a recent member of the cyber-security industry, I've decided to maintain this blog to help foster and " +
		"share some of my experiences in the hope that others may find them useful. The goal is to dedicate this website " +
		"to the study of distributed systems and cyber-security, in addition to all topics I deem related along the way. " +
		"Distributed systems are all around us and in most of the web services we use nowadays. Considering that they have " +
		"only gotten more popular over the last twenty years, I believe they deserve a more thorough and discerning study. " +
		"That also means that as an industry, through this study, we are quickly discovering that the web services we use " +
		"in our day to day life are not as secure as we once imagined them to be. With different exploits like Spectre and " +
		"Meltdown, boundaries that we once thought were insurmontable, like those between two virtual machines running on " +
		"the same host, are now demonstratively not so impenetrable. One of the best goals that this blog could help to " +
		"achieve is shedding a better light on how to make distributed systems more secure than they currently are.</li>" +

		"<li class=\"list-group-item\" style=\"text-align: center;\">" +
		"<a href=\"https://github.com/etzelm\"target=\"_blank\">" +
		"<img src=\"/public/github.png\" alt=\"Github\"height=\"45\" width=\"45\"></a>" +
		"&emsp;&emsp;<a href=\"mailto:etzelm@live.com\">" +
		"<img src=\"/public/email.png\" alt=\"Email\"height=\"30\" width=\"45\"></a>" +
		"&emsp;&emsp;<a href=\"https://www.linkedin.com/in/etzelm/\" target=\"_blank\">" +
		"<img src=\"/public/linkedin.png\" alt=\"LinkedIn\"height=\"45\" width=\"45\"></a></li>" +

		"</div>" +
		"</div>" +
		"<p style=\"text-align: center;\">Source code for this website can be located: " +
		"<a style=\"color:#A619FF;\" href=\"https://github.com/etzelm/blog-in-golang\" target=\"_blank\">" +
		"here</a></p>" +
		"</div>"

	info := ItemInfo{
		Title:    "About Me",
		Created:  created,
		Modified: modified,
		Blurb:    blurb,
		Content:  hold,
	}

	item := Item{
		ID:   0,
		Info: info,
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
}
