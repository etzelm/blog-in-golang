package main

import (
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/etzelm/blog-in-golang/src/handlers"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

// Prometheus metrics. Labels:
//   - method:   HTTP method (GET, POST, …)
//   - route:    c.FullPath() — the route TEMPLATE (e.g. /article/:article_id),
//               never the raw path. Keeps label cardinality bounded and prevents
//               path parameters / query values from leaking into Prometheus.
//   - status:   HTTP status code as string.
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "blog_http_requests_total",
			Help: "Total HTTP requests served, labeled by method, route template, and status code.",
		},
		[]string{"method", "route", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "blog_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds, labeled by method, route template, and status code.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)
)

type envHook struct{ env string }

func (h *envHook) Levels() []log.Level { return log.AllLevels }
func (h *envHook) Fire(e *log.Entry) error { e.Data["env"] = h.env; return nil }

func init() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	envName := os.Getenv("ENV_NAME")
	if envName == "" {
		envName = "dev"
	}
	log.AddHook(&envHook{env: envName})

	if lokiURL := os.Getenv("LOKI_URL"); lokiURL != "" {
		opts := lokirus.NewLokiHookOptions().
			WithBasicAuth(os.Getenv("LOKI_USERNAME"), os.Getenv("LOKI_PASSWORD")).
			WithFormatter(&log.JSONFormatter{
				TimestampFormat: time.RFC3339Nano,
			}).
			WithStaticLabels(lokirus.Labels{
				"service": "blog-in-golang-gcp",
				"env":     "gcp",
				"source":  "lokirus",
			})
		hook := lokirus.NewLokiHookWithOpts(
			lokiURL,
			opts,
			log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel,
		)
		log.AddHook(hook)
	}

	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

var RandomOne int = randRange(1, 9)
var RandomTwo int = randRange(1, 9)

func main() {

	log.Info("Server is starting...")
	gin.SetMode(gin.ReleaseMode)
	// gin.New() instead of gin.Default() — Default() auto-registers Gin's
	// own text-format access logger and writes lines like
	//   [GIN] 2026/05/28 - 20:15:28 | 200 | 64.417µs | ::1 | GET "/healthz"
	// which Loki's `| json` parser can't read, defeating the structured-log
	// work in init(). Recovery() still gets us panic protection; the
	// equivalent of Default()'s Logger() lives in accessLogMiddleware below
	// and emits one JSON line per request via logrus.
	httpServer := gin.New()
	httpServer.Use(gin.Recovery())
	httpServer.LoadHTMLGlob("templates/*")
	// Middlewares MUST be registered before routes — Gin snapshots
	// engine.Handlers into each route's frozen handler chain at registration
	// time (combineHandlers). Use() called AFTER routes are registered only
	// affects the 404/405 handlers (rebuild404Handlers / rebuild405Handlers),
	// not the real routes — which is why CSP, X-Frame-Options, Permissions-
	// Policy, etc. were silently NOT being applied to any 200 response.
	//
	// metricsMiddleware is registered FIRST (in LoadMiddlewares) so it
	// observes total end-to-end request time, including downstream middleware.
	LoadMiddlewares(httpServer)
	LoadStaticFileRoutes(httpServer)
	LoadServerRoutes(httpServer)
	// Don't log the *gin.Engine as a structured field — it embeds
	// `gin.HandlerFunc` slices that json.Marshal can't serialize, so
	// JSONFormatter / lokirus drop the line with
	//   "failed to marshal fields to JSON, json: unsupported type: gin.HandlerFunc".
	log.Info("Gin server created.")

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

	// Liveness probe — used by the container healthcheck and any external
	// uptime monitors. Deliberately cheap: no DynamoDB / S3 / external calls.
	// Just confirms the process is alive and the HTTP server is responding.
	server.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Prometheus scrape endpoint. Gated by metricsAuth (Bearer ${METRICS_TOKEN}
	// in prod; open if METRICS_TOKEN is unset, for local dev). Wrapped with
	// gin.WrapH because promhttp.Handler() is a stdlib http.Handler.
	server.GET("/metrics", metricsAuth(), gin.WrapH(promhttp.Handler()))

}

// LoadMiddlewares loads third party and custom gin middlewares the server uses.
func LoadMiddlewares(server *gin.Engine) {

	server.Use(metricsMiddleware())
	server.Use(accessLogMiddleware())
	server.Use(securityHeadersMiddleware())
	server.Use(staticCacheMiddleware())
	server.Use(unauthorizedMiddleware())
	// Exclude /metrics — promhttp.Handler() already compresses when the client
	// sends Accept-Encoding: gzip (Prometheus does). Wrapping it here would
	// gzip the response a second time, and Prometheus's single decompression
	// pass would leave it staring at "\x1f" (inner gzip magic) and emit:
	//   "expected a valid start token, got \"\\x1f\" (\"INVALID\")"
	server.Use(gzip.Gzip(gzip.BestCompression, gzip.WithExcludedPaths([]string{"/metrics"})))

}

// accessLogMiddleware emits one structured JSON line per request via
// logrus, replacing Gin's text-format access logger. Fields stay bounded
// (route TEMPLATE not raw path; numeric status; integer-µs latency), so
// Loki streams don't blow up on cardinality. Level escalates to Warn at
// 4xx and Error at 5xx so `{service="..."} | json | level="error"` in
// LogQL surfaces actual failures rather than every healthcheck ping.
func accessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		route := c.FullPath() // route TEMPLATE — same source as metrics labels
		if route == "" {
			route = "unmatched"
		}
		status := c.Writer.Status()
		latencyMicros := time.Since(start).Microseconds()
		entry := log.WithFields(log.Fields{
			"client_ip":  c.ClientIP(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"route":      route,
			"status":     status,
			"latency_us": latencyMicros,
		})
		switch {
		case status >= 500:
			entry.Error("request")
		case status >= 400:
			entry.Warn("request")
		default:
			entry.Info("request")
		}
	}
}

// metricsMiddleware records request count + duration into the Prometheus
// counter/histogram vecs declared at package level. Skips /metrics itself so
// the scrape doesn't self-observe.
func metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		route := c.FullPath() // route TEMPLATE, not raw path — see vec docs
		if route == "" {
			route = "unmatched" // 404s without a registered template
		}
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		httpRequestsTotal.WithLabelValues(method, route, status).Inc()
		httpRequestDuration.WithLabelValues(method, route, status).Observe(time.Since(start).Seconds())
	}
}

// metricsAuth gates /metrics behind a bearer token when METRICS_TOKEN is set
// in the environment. If unset, /metrics is open — convenient for local dev,
// closed in prod by setting the secret. Token is captured at middleware
// construction time (process start) so we avoid a syscall per scrape.
func metricsAuth() gin.HandlerFunc {
	expected := os.Getenv("METRICS_TOKEN")
	return func(c *gin.Context) {
		if expected == "" {
			c.Next()
			return
		}
		if c.Request.Header.Get("Authorization") != "Bearer "+expected {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
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
		// X-XSS-Protection deliberately NOT set — Chrome removed the XSSAuditor in
		// 2019 and the header has been known to introduce its own side-channel
		// vulnerabilities. Modern guidance (OWASP, MDN) is to omit it and rely on
		// the CSP below. https://owasp.org/www-project-secure-headers/#x-xss-protection
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Content Security Policy (relaxed for external resources)
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' https://files.mitchelletzel.com https://accounts.google.com; " +
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://accounts.google.com; " +
			"img-src 'self' data: https: blob:; " +
			"font-src 'self' https://fonts.gstatic.com; " +
			"connect-src 'self' https://files.mitchelletzel.com https://accounts.google.com; " +
			"frame-src https://accounts.google.com; " +
			"frame-ancestors 'none';"
		c.Header("Content-Security-Policy", csp)

		c.Next()
	}
}

// randRange generates a random integer between min and max, inclusive.
func randRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}
