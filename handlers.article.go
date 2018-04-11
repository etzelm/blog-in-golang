package main

import (
	"net/http"
	"strconv"
	"strings"

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
	if article, err := getArticleByID(0, "About Me"); err == nil {
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

func GraphStore(c *gin.Context) {
	if article, err := getArticleByID(1,
		"Creating a Scalable, Fault Tolerant, & Strongly Consistent Graph Store API"); err == nil {
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

func getArticle(c *gin.Context) {
	// Check if the article ID is valid
	if articleInfo := c.Param("article_info"); articleInfo != "" {
		articleID, err := strconv.Atoi((strings.Split(articleInfo, "@"))[0])
		log.Info("Info: ", articleInfo)
		log.Info("ID: ", articleID)
		log.Info("Error: ", err)
		if title := (strings.Split(articleInfo, "@"))[1]; title != "" {
			log.Info("Title: ", title)
			// Check if the article exists
			if article, err := getArticleByID(articleID, title); err == nil {
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
			// If an invalid title is specified in the URL, abort with an error
			c.AbortWithStatus(http.StatusNotFound)
		}
	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}
