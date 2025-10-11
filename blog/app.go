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
	"github.com/gin-contrib/gzip"
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
	LoadStaticFileRoutes(httpServer)
	LoadServerRoutes(httpServer)
	LoadMiddlewares(httpServer)
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

// LoadStaticFileRoutes loads all api routes that serve static paths to server files.
func LoadStaticFileRoutes(server *gin.Engine) {

	server.StaticFile("/robots.txt", "./public/robots.txt")
	server.StaticFile("/sitemap.xml", "./public/sitemap.xml")
	server.StaticFile("/favicon.ico", "./public/images/favicon.ico")
	server.Use(static.Serve("/public", static.LocalFile("./public", true)))
	server.Use(static.Serve("/realtor", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/new", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/search", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/listing", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/my-listing", static.LocalFile("./realtor/build", true)))
	server.Use(static.Serve("/realtor/my-listings", static.LocalFile("./realtor/build", true)))

}

// LoadServerRoutes loads all the custom api calls I've written for the server.
func LoadServerRoutes(server *gin.Engine) {

	store := persistence.NewInMemoryStore(365 * 24 * time.Hour)
	server.GET("/", cache.CachePage(store, 365*24*time.Hour, handlers.AboutPage))
	server.GET("/posts", cache.CachePage(store, 365*24*time.Hour, handlers.PostPage))
	server.GET("/article/:article_id", cache.CachePage(store, 365*24*time.Hour, handlers.ArticlePage))
	server.GET("/category/:category", cache.CachePage(store, 365*24*time.Hour, handlers.CategoryPage))
	server.GET("/contact", handlers.ContactPage(&RandomOne, &RandomTwo))
	server.POST("/contact", handlers.ContactResponse(&RandomOne, &RandomTwo))
	server.GET("/listing/:listing", handlers.ListingGETAPI)
	server.GET("/listings", handlers.ListingsGETAPI)
	server.GET("/auth", handlers.AuthPage)
	server.POST("/auth", handlers.AuthResponse)
	server.GET("/secure", handlers.SecurePage)
	server.POST("/listings/add/:key", handlers.ListingPOSTAPI)
	server.POST("/upload/image/:user", handlers.UploadImagePOSTAPI)

}

// LoadMiddlewares loads third party and custom gin middlewares the server uses.
func LoadMiddlewares(server *gin.Engine) {

	server.Use(securityHeadersMiddleware())
	server.Use(staticCacheMiddleware())
	server.Use(unauthorizedMiddleware())
	server.Use(gzip.Gzip(gzip.BestCompression))

}

// staticCacheMiddleware adds optimized caching headers for static files with proper cache strategies
func staticCacheMiddleware() gin.HandlerFunc {
	staticPrefixes := []string{
		"/public/", "/favicon.ico", "/robots.txt", "/sitemap.xml",
		"/realtor/js/", "/realtor/css/", "/realtor/images/", "/realtor/static/",
	}
	return func(c *gin.Context) {
		for _, prefix := range staticPrefixes {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				// Enhanced caching with ETag support
				c.Header("Cache-Control", "public, max-age=31536000, immutable")
				c.Header("Vary", "Accept-Encoding")
				// Add security headers for static assets
				c.Header("X-Content-Type-Options", "nosniff")
				break
			}
		}
		c.Next()
	}
}

// unauthorizedMiddleware blocks access to common malicious request paths
func unauthorizedMiddleware() gin.HandlerFunc {
	blockedPatterns := []string{
		"wp-includes", "wp-content", "wp-login", ".git", "admin", ".php",
	}
	return func(c *gin.Context) {
		for _, pattern := range blockedPatterns {
			if strings.Contains(c.Request.URL.Path, pattern) {
				c.AbortWithStatus(401)
				return
			}
		}
		c.Next()
	}
}

// securityHeadersMiddleware adds comprehensive security and performance headers
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Content Security Policy (relaxed for external resources)
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' https://files.mitchelletzel.com; " +
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
			"img-src 'self' data: https: blob:; " +
			"font-src 'self' https://fonts.gstatic.com; " +
			"connect-src 'self' https://files.mitchelletzel.com; " +
			"frame-ancestors 'none';"
		c.Header("Content-Security-Policy", csp)

		c.Next()
	}
}

// randRange generates a random integer between min and max, inclusive.
func randRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}
