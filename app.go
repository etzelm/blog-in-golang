package main

import (
	"math/rand/v2"
	"os"
	"strings"
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
		for range time.Tick(time.Hour * 3) {
			go func() {
				RandomOne = randRange(1, 9)
				RandomTwo = randRange(1, 9)
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

	server.Use(staticCacheMiddleware())
	server.Use(static.Serve("/public", static.LocalFile("./public", true)))
	server.Use(static.Serve("/favicon.ico", static.LocalFile("./public/images/favicon.ico", false)))
	server.Use(static.Serve("/robots.txt", static.LocalFile("./public/robots.txt", false)))
	server.Use(static.Serve("/sitemap.xml", static.LocalFile("./public/sitemap.xml", false)))
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

	store := persistence.NewInMemoryStore(365 * 24 * time.Hour)
	server.GET("/", cache.CachePage(store, 365*24*time.Hour, handlers.AboutPage))
	server.GET("/posts", cache.CachePage(store, 365*24*time.Hour, handlers.PostPage))
	server.GET("/article/:article_id", cache.CachePage(store, 365*24*time.Hour, handlers.ArticlePage))
	server.GET("/category/:category", cache.CachePage(store, 365*24*time.Hour, handlers.CategoryPage))

	server.Use(unauthorizedMiddleware())
	server.GET("/contact", handlers.ContactPage(&RandomOne, &RandomTwo))
	server.POST("/contact", handlers.ContactResponse(&RandomOne, &RandomTwo))
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

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine if request path is for statically served files
		if strings.HasPrefix(c.Request.URL.Path, "/public/") ||
			strings.HasPrefix(c.Request.URL.Path, "/favicon.ico") ||
			strings.HasPrefix(c.Request.URL.Path, "/robots.txt") ||
			strings.HasPrefix(c.Request.URL.Path, "/sitemap.xml") ||
			strings.HasPrefix(c.Request.URL.Path, "/realtor/js/") ||
			strings.HasPrefix(c.Request.URL.Path, "/realtor/css/") ||
			strings.HasPrefix(c.Request.URL.Path, "/realtor/images/") ||
			strings.HasPrefix(c.Request.URL.Path, "/realtor/static/") {
			// Apply the Cache-Control header to the static files
			c.Header("Cache-Control", "public, max-age=31536000")
		}
		// Continue to the next middleware or handler
		c.Next()
	}
}

func unauthorizedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine if request path is unused pattern common in scans
		if strings.Contains(c.Request.URL.Path, "wp-includes") ||
			strings.Contains(c.Request.URL.Path, "wp-content") ||
			strings.Contains(c.Request.URL.Path, "wp-login") ||
			strings.Contains(c.Request.URL.Path, ".git") ||
			strings.Contains(c.Request.URL.Path, "admin") ||
			strings.Contains(c.Request.URL.Path, "php") {
			// Abort the gin context while returning 401
			c.AbortWithStatus(401)
			return
		}
		// Continue to the next middleware or handler
		c.Next()
	}
}
