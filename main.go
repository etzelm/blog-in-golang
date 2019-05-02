package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

// IP_PORT is this computers IP address
var IP_PORT string

func main() {
	log.Info("Server is starting...")

	gin.SetMode(gin.ReleaseMode)
	httpsServer := gin.Default()
	httpsServer.LoadHTMLGlob("templates/*")
	httpsServer.Use(static.Serve("/public", static.LocalFile("./public", true)))
	LoadRoutes(httpsServer)
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("mitchelletzel.com"),
		Cache:      autocert.DirCache("certs"),
	}

	log.WithField("server", httpsServer).Info("Default Gin server created.")
	log.Info(autotls.RunWithManager(httpsServer, &m))
}

// LoadRoutes does exactly that... loads all routes for the server.
func LoadRoutes(server *gin.Engine) *gin.Engine {
	server.GET("/", LandingPage)
	server.GET("/about", AboutPage)
	server.GET("/feedback", FeedbackPage)
	server.POST("/feedback", FeedbackResponse)
	server.GET("/article/:article_id", getArticle)
	return server
}
