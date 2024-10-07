package handlers

import (
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
	"github.com/etzelm/blog-in-golang/src/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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
			c.HTML(
				// Set the HTTP status to 404 (Not Found)
				http.StatusNotFound,
				// Use the error.html template
				"error.html",
				// Pass the data that the page uses
				gin.H{
					"title": "404 Server Error",
					"error": "Please provide a valid category",
				},
			)

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
		c.HTML(
			// Set the HTTP status to 404 (Not Found)
			http.StatusNotFound,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "404 Server Error",
				"error": "Please provide a valid category",
			},
		)
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
				c.HTML(
					// Set the HTTP status to 401 (Unauthorized)
					http.StatusUnauthorized,
					// Use the error.html template
					"error.html",
					// Pass the data that the page uses
					gin.H{
						"title": "401 (Unauthorized)",
						"error": "Please provide a valid Article ID.",
					},
				)
			}
		} else {
			// If the article is not found, abort with an error
			c.HTML(
				// Set the HTTP status to 404 (Not Found)
				http.StatusNotFound,
				// Use the error.html template
				"error.html",
				// Pass the data that the page uses
				gin.H{
					"title": "404 Server Error",
					"error": "Please provide a valid Article ID.",
				},
			)
		}
	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.HTML(
			// Set the HTTP status to 404 (Not Found)
			http.StatusNotFound,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "404 Server Error",
				"error": "Please provide a valid Article ID.",
			},
		)
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
		c.Bind(&form)

		if form.RobotCheck != 1 || form.RobotNum != *numOne+*numTwo {
			c.HTML(
				// Set the HTTP status to 401 (Unauthorized)
				http.StatusUnauthorized,
				// Use the error.html template
				"error.html",
				// Pass the data that the page uses
				gin.H{
					"title": "401 (Unauthorized)",
					"error": "Don't be a robot please!",
				},
			)
			return
		}

		name := template.HTMLEscapeString(form.Name)
		if m, _ := regexp.MatchString("^[ a-zA-Z0-9]+( +[a-zA-Z0-9]+)*$", name); !m {
			c.HTML(
				// Set the HTTP status to 401 (Unauthorized)
				http.StatusUnauthorized,
				// Use the error.html template
				"error.html",
				// Pass the data that the page uses
				gin.H{
					"title": "401 (Unauthorized)",
					"error": "Name should contain only alphanumeric characters and spaces!",
				},
			)
			return
		}

		aid := os.Getenv("AWS_ACCESS_KEY_ID")
		key := os.Getenv("AWS_SECRET_ACCESS_KEY")
		var myCredentials = credentials.NewStaticCredentials(aid, key, "")

		sess, err := session.NewSession(&aws.Config{
			Credentials: myCredentials,
			Region:      aws.String("us-west-1"),
			//Endpoint:    aws.String("http://localhost:8000"),
		})
		if err != nil {
			log.Error(err)
			c.HTML(
				// Set the HTTP status to 500 (Internal Server Error)
				http.StatusInternalServerError,
				// Use the error.html template
				"error.html",
				// Pass the data that the page uses
				gin.H{
					"title": "500 Internal Server Error",
					"error": err.Error(),
				},
			)
			return
		}

		dbSvc := dynamodb.New(sess)

		av, _ := dynamodbattribute.MarshalMap(form)

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Contact"),
		}

		_, err = dbSvc.PutItem(input)

		if err != nil {
			log.Error("Got error calling PutItem:")
			log.Error(err.Error())
			c.HTML(
				// Set the HTTP status to 500 (Internal Server Error)
				http.StatusInternalServerError,
				// Use the error.html template
				"error.html",
				// Pass the data that the page uses
				gin.H{
					"title": "500 Internal Server Error",
					"error": err.Error(),
				},
			)
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
