package main

import (
	"os"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/etzelm/blog-in-golang/src/handlers"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Info("Server is starting...")
	gin.SetMode(gin.ReleaseMode)
	httpServer := gin.Default()

	httpServer.LoadHTMLGlob("templates/*")
	LoadStaticFolderRoutes(httpServer)
	LoadServerRoutes(httpServer)
	log.WithField("server", httpServer).Info("Default Gin server created.")

	env := os.Getenv("DEPLOYMENT")
	domain := os.Getenv("DOMAIN")
	if env == "NAS" {
		httpServer.RunTLS(":8080", "/config/ssl/live/"+domain+"/fullchain.pem", "/config/ssl/live/"+domain+"/privkey.pem")
	} else if env == "GCP" {
		certmagic.DefaultACME.Agreed = true
		certmagic.DefaultACME.Email = "etzelm@live.com"
		log.Info(certmagic.HTTPS([]string{domain}, httpServer))
	} else {
		httpServer.Run()
	}

}

// LoadStaticFolderRoutes loads all api routes that serve a static server folder.
func LoadStaticFolderRoutes(server *gin.Engine) *gin.Engine {

	server.Use(static.Serve("/public", static.LocalFile("./public", true)))
	server.Use(static.Serve("/realtor", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/new", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/search", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/listing", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/my-listing", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/my-listings", static.LocalFile("./realtor/build", true)))
	return server

}

// LoadServerRoutes does exactly that... loads all api routes for the server.
func LoadServerRoutes(server *gin.Engine) *gin.Engine {

	store := persistence.NewInMemoryStore(24 * time.Hour)
	server.GET("/", cache.CachePage(store, 7*24*time.Hour, handlers.AboutPage))
	server.GET("/posts", cache.CachePage(store, 24*time.Hour, handlers.PostPage))
	server.GET("/contact", cache.CachePage(store, 7*24*time.Hour, handlers.ContactPage))
	server.POST("/contact", handlers.ContactResponse)
	server.GET("/article/:article_id", cache.CachePage(store, 24*time.Hour, handlers.ArticlePage))
	server.GET("/category/:category", cache.CachePage(store, 24*time.Hour, handlers.CategoryPage))
	server.GET("/listing/:listing", handlers.ListingGETAPI)
	server.GET("/listings", handlers.ListingsGETAPI)
	server.GET("/auth", handlers.AuthPage)
	server.POST("/auth", handlers.AuthResponse)
	server.GET("/secure", handlers.SecurePage)
	server.POST("/listings/add/:key", handlers.ListingPOSTAPI)
	server.POST("/upload/image/:user", handlers.UploadImagePOSTAPI)
	return server

}
