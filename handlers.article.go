package main

import (
	"html/template"
	"net/http"
	"regexp"
	"strconv"

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
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"feedback.html",
		// Pass the data that the page uses
		gin.H{
			"title": "Leave Feedback",
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
