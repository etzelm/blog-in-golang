package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Hello is the controller for the "/hello" route. Has to match query 'name' or
// address the user as 'user'.
func Hello(c *gin.Context) {
	name, hasName := c.GetQuery("name")
	log.WithFields(log.Fields{"name": name, "hasName": hasName}).Info("Hello request query string -->")
	if hasName == false {
		c.String(http.StatusOK, "Hello user!")
		return
	}
	c.String(http.StatusOK, "Hello %s!", name)
}

// CheckGet is used for all GET requests to the '/check' path.
func CheckGet(c *gin.Context) {
	c.String(http.StatusOK, "This is a GET request")
}

// CheckPost is used for all POST requests to the '/check' path.
func CheckPost(c *gin.Context) {
	c.String(http.StatusOK, "This is a POST request")
}

// CheckPut is used for all PUT requests to the '/check' path. Needed to explicitly
// issue a 405 (method not allowed) instead of the default 404 (not found).
func CheckPut(c *gin.Context) {
	c.AbortWithStatus(http.StatusMethodNotAllowed)
}

func LandingPage(c *gin.Context) {
	c.String(http.StatusOK, "Hello!")
}
