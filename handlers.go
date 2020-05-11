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

//PostPage : Gets All Article Panels and Dynamically Displays index.html Template
func PostPage(c *gin.Context) {
	panels := getArticlePanels()

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

//CategoryPage : Gets Category Article Panels and Dynamically Displays index.html Template
func CategoryPage(c *gin.Context) {
	if category := c.Param("category"); category != "" {
		panels := getCategoryPageArticlePanels(category)

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

//ArticlePage : Queries DynamoDB for a Specific Article and Dynamically Displays article.html
func ArticlePage(c *gin.Context) {
	// Check if the article ID is valid
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		// Check if the article exists
		if article, err := getArticleByID(articleID); err == nil {
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
					// Set the HTTP status to 403 (Forbidden)
					http.StatusForbidden,
					// Use the error.html template
					"error.html",
					// Pass the data that the page uses
					gin.H{
						"title": "403 Server Error",
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

//AboutPage : Displays the static about.html page
func AboutPage(c *gin.Context) {
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

//ContactPage : Displays the static contact.html page for GET requests
func ContactPage(c *gin.Context) {
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"contact.html",
		// Pass the data that the page uses
		gin.H{
			"title": "Contact Me",
		},
	)
}

//ContactResponse : Saves the user's data in DynamoDB and displays static response.html
func ContactResponse(c *gin.Context) {
	var form ContactForm
	c.Bind(&form)

	if form.RobotCheck != 1 {
		c.HTML(
			// Set the HTTP status to 400 (Bad Request)
			http.StatusBadRequest,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "400 Server Error",
				"error": "Don't be a robot please!",
			},
		)
		return
	}

	name := template.HTMLEscapeString(form.Name)
	if m, _ := regexp.MatchString("^[ a-zA-Z0-9]+( +[a-zA-Z0-9]+)*$", name); !m {
		c.HTML(
			// Set the HTTP status to 400 (Bad Request)
			http.StatusBadRequest,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "400 Server Error",
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
		log.Println(err)
		c.HTML(
			// Set the HTTP status to 400 (Bad Request)
			http.StatusBadRequest,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "400 Server Error",
				"error": err.Error(),
			},
		)
		return
	}

	dbSvc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(form)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Contact"),
	}

	_, err = dbSvc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		c.HTML(
			// Set the HTTP status to 400 (Bad Request)
			http.StatusBadRequest,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "400 Server Error",
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

//ListingsAPI : Gets All Realtor Listings
func ListingsAPI(c *gin.Context) {
	listings := getRealtorListings()

	// Call the JSON method of the Context to return the results
	c.JSON(200, listings)

}

//ListingAPI : Gets A Realtor Listing
func ListingAPI(c *gin.Context) {

	if listing := c.Param("listing"); listing != "" {

		card := getRealtorListing(listing)

		// Call the JSON method of the Context to return the results
		c.JSON(200, card)

	} else {

		empty := []byte(``)

		// Call the JSON method of the Context to 404
		c.JSON(404, empty)

	}

}
