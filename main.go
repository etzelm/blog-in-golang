package main

import (
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// IP_PORT is this computers IP address
var IP_PORT string

func main() {
	log.Info("Server is starting...")

	IP_PORT := os.Getenv("ip_port")
	if IP_PORT == "" {
		IP_PORT = ":3000"
	}
	log.Info("IP_PORT: ", IP_PORT)

	server := gin.Default()
	log.WithField("server", server).Info("Default Gin server create.")
	server.LoadHTMLGlob("templates/*")
	server.Use(static.Serve("/public", static.LocalFile("./public", true)))
	LoadRoutes(server)
	server.Run(IP_PORT)
}

// LoadRoutes does exactly that... loads all routes for the server.
func LoadRoutes(server *gin.Engine) *gin.Engine {
	server.GET("/", LandingPage)
	server.GET("/about", AboutPage)
	server.GET("/graph_store", GraphStore)
	server.GET("/article/view/:article_info", getArticle)
	return server
}
