package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LandingPage(c *gin.Context) {
	articles := getAllArticles()

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses
		gin.H{
			"title":   "Home Page",
			"payload": articles,
		},
	)

}

func AboutPage(c *gin.Context) {
	if article, err := getArticleByID(0); err == nil {
		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"article.html",
			// Pass the data that the page uses
			gin.H{
				"title": article.Title,
			},
		)
		//Write Stored HTML from mongoDB to article.html
		c.Writer.Write([]byte(article.Content))
	} else {
		// If the article is not found, abort with an error
		c.AbortWithError(http.StatusNotFound, err)
	}
}

func FeedbackPage(c *gin.Context) {
	feedbacks := getAllFeedback()

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"feedback.html",
		// Pass the data that the page uses
		gin.H{
			"title":   "Leave Feedback",
			"payload": feedbacks,
		},
	)
}

func FeedbackResponse(c *gin.Context) {
	var form FeedbackForm
	c.Bind(&form)

	name := template.HTMLEscapeString(form.Name)
	feedback := template.HTMLEscapeString(form.Feedback)
	if m, _ := regexp.MatchString("^[ a-zA-Z0-9]+( +[a-zA-Z0-9]+)*$", name); !m {
		c.AbortWithStatusJSON(400, "Name should contain only alphanumeric characters and spaces!")
		return
	}
	log.Info("Name: ", name)
	log.Info("Feedback: ", feedback)

	if form.X == 1 {
		log.Info("Public")
	} else {
		log.Info("Private")
	}

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
		return
	}

	dbSvc := dynamodb.New(sess)

	info := FeedbackInfo{
		Feedback: feedback,
		X:        form.X,
	}

	item := FeedbackItem{
		Name: name,
		Info: info,
	}

	av, err := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Feedback"),
	}

	_, err = dbSvc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	feedbacks := getAllFeedback()

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"response.html",
		// Pass the data that the page uses
		gin.H{
			"title":   "Thank You!",
			"payload": feedbacks,
		},
	)
}

func getArticle(c *gin.Context) {
	// Check if the article ID is valid
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		// Check if the article exists
		if article, err := getArticleByID(articleID); err == nil {
			// Call the HTML method of the Context to render a template
			c.HTML(
				// Set the HTTP status to 200 (OK)
				http.StatusOK,
				// Use the index.html template
				"article.html",
				// Pass the data that the page uses
				gin.H{
					"title": article.Title,
				},
			)
			//Write Stored HTML from mongoDB to article.html
			c.Writer.Write([]byte(article.Content))
		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}
