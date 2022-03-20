package infraCode

import (
	"fmt"
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

// infraCode function for addition/modification of infraCode article
func InfraCode() {
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

	data, _ := ioutil.ReadFile("articles/infraCode/articlePicture.html")
	ap := string(data)

	data, _ = ioutil.ReadFile("articles/infraCode/panelPicture.html")
	pp := string(data)

	data, _ = ioutil.ReadFile("articles/infraCode/infraCode.html")
	hh := string(data)

	item := models.Item{
		ArticlePicture: ap,
		Author:         "<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>",
		Categories:     "Disciplines,Distributed Systems",
		CreatedDate:    "March 19th, 2022",
		Excerpt: "CloudFormation vs. CDK vs. Serverless Framework vs. Terraform: If you've ever had to deploy a " +
			"repeatable set of AWS Infrastructure then it's likely that you've come across these tools. Let's " +
			"jump into some of their strengths and weaknesses.",
		HTMLHold:     hh,
		ModifiedDate: "March 20th, 2022",
		PanelPicture: pp,
		PostID:       5,
		PostTitle:    "A look at Infrastructure as Code in the AWS Cloud",
		ShortTitle:   "IAC Compare and Contrast",
		PostType:     "standard",
	}

	av, _ := dynamodbattribute.MarshalMap(item)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Test-Articles"),
	}

	_, err = dbSvc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
