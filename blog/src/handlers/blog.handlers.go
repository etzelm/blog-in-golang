package handlers

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/etzelm/blog-in-golang/src/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// createAWSConfig creates AWS configuration with proper credentials
func createAWSConfig(ctx context.Context) (aws.Config, error) {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")

	var cfg aws.Config
	var err error

	if aid != "" && key != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion("us-west-1"),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(aid, key, "")),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion("us-west-1"),
		)
	}

	return cfg, err
}

// createDynamoDBClient creates a DynamoDB client with proper configuration
func createDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := createAWSConfig(ctx)
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}

// PostPage : Gets All Article Panels and Dynamically Displays index.html Template
func PostPage(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000")
	panels := models.GetArticlePanels()

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses
		gin.H{
			"title":   "Blog Posts",
			"payload": panels,
		},
	)

}

// CategoryPage : Gets Category Article Panels and Dynamically Displays index.html Template
func CategoryPage(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000")
	if category := c.Param("category"); category != "" {
		panels := models.GetCategoryPageArticlePanels(category)

		if len(panels) <= 0 {
			// If an invalid category is specified in the URL, abort with an error
			renderErrorPage(c, 404, "404 (Not Found)", "Please provide a valid category")
			return
		}

		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"index.html",
			// Pass the data that the page uses
			gin.H{
				"payload":    panels,
				"category":   category,
				"title":      category,
				"IsCategory": true,
			},
		)

	} else {
		// If an invalid category is specified in the URL, abort with an error
		renderErrorPage(c, 404, "404 (Not Found)", "Please provide a valid category")
	}

}

// ArticlePage : Queries DynamoDB for a Specific Article and Dynamically Displays article.html
func ArticlePage(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000")
	// Check if the article ID is valid
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		// Check if the article exists
		if article, err := models.GetArticleByID(articleID); err == nil {
			// Check the post type for appropriateness
			if article.PostType != "quote" && article.PostType != "" {
				// Call the HTML method of the Context to render a template
				c.HTML(
					// Set the HTTP status to 200 (OK)
					http.StatusOK,
					// Use the index.html template
					"article.html",
					// Pass the data that the page uses
					gin.H{
						"title":   article.ShortTitle,
						"payload": article,
					},
				)
			} else {
				// If the article is not appropriate, abort with an error
				renderErrorPage(c, 401, "401 (Unauthorized)", "Please provide a valid Article ID.")
			}
		} else {
			// If the article is not found, abort with an error
			renderErrorPage(c, 404, "404 (Not Found)", "Please provide a valid Article ID.")
		}
	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		renderErrorPage(c, 404, "404 (Not Found)", "Please provide a valid Article ID.")
	}
}

// AboutPage : Displays the static about.html page
func AboutPage(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000")
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the about.html template
		"about.html",
		// Pass the data that the page uses
		gin.H{
			"title": "Mitchell Etzel",
		},
	)
}

// ContactPage : Displays the static contact.html page for GET requests
func ContactPage(numOne *int, numTwo *int) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"contact.html",
			// Pass the data that the page uses
			gin.H{
				"title":     "Contact Me",
				"RandomOne": *numOne,
				"RandomTwo": *numTwo,
			},
		)
	}
	return gin.HandlerFunc(fn)
}

// ContactResponse : Saves the user's data in DynamoDB and displays static response.html
func ContactResponse(numOne *int, numTwo *int) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		var form models.ContactForm
		if err := c.Bind(&form); err != nil {
			renderErrorPage(c, 400, "400 (Bad Request)", "Invalid form data.")
			return
		}

		if form.RobotCheck != 1 || form.RobotNum != *numOne+*numTwo {
			renderErrorPage(c, 401, "401 (Unauthorized)", "Don't be a robot please!")
			return
		}

		name := template.HTMLEscapeString(form.Name)
		m, err := regexp.MatchString("^[ a-zA-Z0-9]+( +[a-zA-Z0-9]+)*$", name)
		if err != nil || !m {
			renderErrorPage(c, 401, "401 (Unauthorized)", "Name should contain only alphanumeric characters & spaces!")
			return
		}

		ctx := context.TODO()
		dbSvc, err := createDynamoDBClient(ctx)
		if err != nil {
			log.Error("Unable to create DynamoDB client:", err)
			renderErrorPage(c, 500, "500 Internal Server Error", err.Error())
			return
		}

		av, err := attributevalue.MarshalMap(form)
		if err != nil {
			log.Error("Error marshalling form:", err)
			renderErrorPage(c, 500, "500 Internal Server Error", err.Error())
			return
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Contact"),
		}

		_, err = dbSvc.PutItem(ctx, input)

		if err != nil {
			log.Error("Got error calling PutItem:")
			log.Error(err.Error())
			renderErrorPage(c, 500, "500 Internal Server Error", err.Error())
			return
		}

		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"response.html",
			// Pass the data that the page uses
			gin.H{
				"title": "Thank You!",
			},
		)
	}
	return gin.HandlerFunc(fn)
}

func renderErrorPage(c *gin.Context, statusCode int, title, message string) {
	c.HTML(
		statusCode,
		"error.html",
		gin.H{
			"title": title,
			"error": message,
		},
	)
}
