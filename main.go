package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

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
	//httpsServer.Run("127.0.0.1:80")
}

// LoadRoutes does exactly that... loads all routes for the server.
func LoadRoutes(server *gin.Engine) *gin.Engine {
	server.GET("/", AboutPage)
	server.GET("/posts", PostPage)
	server.GET("/contact", ContactPage)
	server.POST("/contact", ContactResponse)
	server.GET("/article/:article_id", ArticlePage)
	server.GET("/category/:category", CategoryPage)
	return server
}
