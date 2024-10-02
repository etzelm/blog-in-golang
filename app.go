package main

import (
	"math/rand/v2"
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

var RandomOne int = randRange(1, 9)
var RandomTwo int = randRange(1, 9)

func main() {

	log.Info("Server is starting...")
	gin.SetMode(gin.ReleaseMode)
	httpServer := gin.Default()

	httpServer.LoadHTMLGlob("templates/*")
	LoadStaticFolderRoutes(httpServer)
	LoadServerRoutes(httpServer)
	log.WithField("server", httpServer).Info("Default Gin server created.")

	go func() {
		for range time.Tick(time.Minute * 3) {
			go func() {
				RandomOne = randRange(1, 9)
				log.Info("Changing RandomOne: ", RandomOne)
				RandomTwo = randRange(1, 9)
				log.Info("Changing RandomTwo: ", RandomTwo)
			}()
		}
	}()

	env := os.Getenv("DEPLOYMENT")
	domain := os.Getenv("DOMAIN")
	if env == "NAS" || env == "GCP" {
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
	server.GET("/contact", handlers.ContactPage(RandomOne, RandomTwo))
	server.POST("/contact", handlers.ContactResponse(RandomOne, RandomTwo))
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

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
