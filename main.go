package main

import (
	"github.com/caddyserver/certmagic"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Server is starting...")

	gin.SetMode(gin.ReleaseMode)
	httpsServer := gin.Default()
	httpsServer.LoadHTMLGlob("templates/*")
	httpsServer.Use(static.Serve("/public", static.LocalFile("./public", true)))
	httpsServer.Use(static.Serve("/realtor", static.LocalFile("./realtor-site/build", true)))
	httpsServer.Use(static.Serve("/realtor/new", static.LocalFile("./realtor-site/build", true)))
	httpsServer.Use(static.Serve("/realtor/listing", static.LocalFile("./realtor-site/build", true)))
	LoadRoutes(httpsServer)

	log.WithField("server", httpsServer).Info("Default Gin server created.")

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = "etzelm@live.com"
	log.Info(certmagic.HTTPS([]string{"blog.mitchelletzel.com"}, httpsServer))

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
	server.GET("/listings", ListingsAPI)
	server.GET("/listing/:listing", ListingAPI)
	return server
}
