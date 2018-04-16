package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
