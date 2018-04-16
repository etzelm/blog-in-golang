package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

type article struct {
	ID       int    `json:"id"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Title    string `json:"title"`
	Blurb    string `json:"blurb"`
	Content  string `json:"content"`
}

type Item struct {
	ID   int      `json:"id"`
	Info ItemInfo `json:"info"`
}

type ItemInfo struct {
	Title    string `json:"title"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Blurb    string `json:"blurb"`
	Content  string `json:"content"`
}

type FeedbackForm struct {
	Name     string `form:"name" binding:"required"`
	Feedback string `form:"feedback" binding:"required"`
	X        int    `form:"x"`
}

type FeedbackItem struct {
	Name string       `json:"name"`
	Info FeedbackInfo `json:"info"`
}

type FeedbackInfo struct {
	Feedback string `json:"feedback"`
	X        int    `json:"x"`
}

// Return a list of all the articles
func getAllArticles() []article {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var my_credentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: my_credentials,
		Region:      aws.String("us-west-1"),
		Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	dbSvc := dynamodb.New(sess)

	filt := expression.Name("id").GreaterThanEqual(expression.Value(0))

	proj := expression.NamesList(expression.Name("info.title"), expression.Name("id"), expression.Name("info.blurb"),
		expression.Name("info.created"), expression.Name("info.modified"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Articles"),
	}

	// Make the DynamoDB Query API call
	result, err := dbSvc.Scan(params)

	temp := []article{}

	for _, i := range result.Items {
		item := Item{}
		article := article{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		article.ID = item.ID
		article.Title = item.Info.Title
		article.Created = item.Info.Created
		article.Modified = item.Info.Modified
		article.Blurb = item.Info.Blurb
		temp = append(temp, article)
	}

	return temp
}

func getArticleByID(id int) (*article, error) {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var my_credentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: my_credentials,
		Region:      aws.String("us-west-1"),
		Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dbSvc := dynamodb.New(sess)

	result, err := dbSvc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Articles"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	item := Item{}
	article := article{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	article.ID = item.ID
	article.Title = item.Info.Title
	article.Blurb = item.Info.Blurb
	article.Created = item.Info.Created
	article.Modified = item.Info.Modified
	article.Content = item.Info.Content

	return &article, nil
}

func getAllFeedback() []FeedbackForm {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var my_credentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: my_credentials,
		Region:      aws.String("us-west-1"),
		Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	dbSvc := dynamodb.New(sess)

	filt := expression.Name("info.x").Equal(expression.Value(1))

	proj := expression.NamesList(expression.Name("info.feedback"), expression.Name("name"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Feedback"),
	}

	// Make the DynamoDB Query API call
	result, err := dbSvc.Scan(params)

	temp := []FeedbackForm{}

	for _, i := range result.Items {
		item := FeedbackItem{}
		hold := FeedbackForm{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		hold.Name = item.Name
		hold.Feedback = item.Info.Feedback
		temp = append(temp, hold)
	}

	return temp
}
