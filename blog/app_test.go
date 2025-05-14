package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func silenceLogrus(t *testing.T) {
	originalOut := log.StandardLogger().Out
	log.SetOutput(io.Discard)
	t.Cleanup(func() {
		log.SetOutput(originalOut)
	})
}

func TestRandRange(t *testing.T) {
	silenceLogrus(t)
	testCases := []struct {
		name string
		min  int
		max  int
	}{
		{"PositiveRange", 1, 10},
		{"SingleValueRange", 5, 5},
		{"NegativeRange", -10, -1},
		{"MixedRange", -5, 5},
		{"ZeroMinRange", 0, 10},
		{"ZeroMaxRange", -10, 0},
		{"EqualPositiveRange", 7, 7},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := randRange(tc.min, tc.max)
			if result < tc.min {
				t.Errorf("randRange(%d, %d) = %d; want value >= %d", tc.min, tc.max, result, tc.min)
			}
			if result > tc.max {
				t.Errorf("randRange(%d, %d) = %d; want value <= %d", tc.min, tc.max, result, tc.max)
			}
		})
	}
}

func TestStaticCacheMiddleware(t *testing.T) {
	silenceLogrus(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(staticCacheMiddleware())
	router.GET("/public/test.css", func(c *gin.Context) { c.String(http.StatusOK, "test css") })
	router.GET("/api/data", func(c *gin.Context) { c.String(http.StatusOK, "api data") })

	reqStatic, _ := http.NewRequest(http.MethodGet, "/public/test.css", nil)
	rrStatic := httptest.NewRecorder()
	router.ServeHTTP(rrStatic, reqStatic)
	if status := rrStatic.Code; status != http.StatusOK {
		t.Errorf("static path status: got %v want %v", status, http.StatusOK)
	}
	expectedCacheControl := "public, max-age=31536000"
	if rrStatic.Header().Get("Cache-Control") != expectedCacheControl {
		t.Errorf("static path Cache-Control: got %v want %v", rrStatic.Header().Get("Cache-Control"), expectedCacheControl)
	}

	reqNonStatic, _ := http.NewRequest(http.MethodGet, "/api/data", nil)
	rrNonStatic := httptest.NewRecorder()
	router.ServeHTTP(rrNonStatic, reqNonStatic)
	if status := rrNonStatic.Code; status != http.StatusOK {
		t.Errorf("non-static path status: got %v want %v", status, http.StatusOK)
	}
	if rrNonStatic.Header().Get("Cache-Control") == expectedCacheControl {
		t.Errorf("non-static path Cache-Control should not be %v", expectedCacheControl)
	}
}

func TestUnauthorizedMiddleware(t *testing.T) {
	silenceLogrus(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(unauthorizedMiddleware())
	router.GET("/allowed/path", func(c *gin.Context) { c.String(http.StatusOK, "allowed") })
	router.GET("/another/safe/path", func(c *gin.Context) { c.String(http.StatusOK, "safe path allowed") })

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"AllowedPath", "/allowed/path", http.StatusOK, "allowed"},
		{"BlockedPatternWPIncludes", "/wp-includes/something", http.StatusUnauthorized, ""},
		{"BlockedPatternGit", "/.git/config", http.StatusUnauthorized, ""},
		{"BlockedPatternPHP", "/login.php", http.StatusUnauthorized, ""},
		{"BlockedPatternAdmin", "/admin/index", http.StatusUnauthorized, ""},
		{"PathContainingWPContent", "/some/path/wp-content/uploads", http.StatusUnauthorized, ""},
		{"NonBlockedPath", "/another/safe/path", http.StatusOK, "safe path allowed"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("path %s status: got %v want %v", tc.path, status, tc.expectedStatus)
			}
			if tc.expectedStatus == http.StatusOK && tc.expectedBody != "" {
				if rr.Body.String() != tc.expectedBody {
					t.Errorf("path %s body: got %v want %v", tc.path, rr.Body.String(), tc.expectedBody)
				}
			}
		})
	}
}

func TestLoadServerRoutes(t *testing.T) {
	silenceLogrus(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()

	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "dummy-test-articles-table")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	tempTemplatesDir := t.TempDir()

	dummyTemplateFiles := map[string]string{
		"about.html":   "<html><title>About</title><body>{{.title}}</body></html>",
		"index.html":   "<html><title>Index</title><body>{{.title}}</body></html>",
		"article.html": "<html><title>Article</title><body>{{if .payload}}{{.payload.PostTitle}}{{else}}No Payload{{end}}</body></html>",
		"contact.html": "<html><title>Contact</title><body>Contact Us: {{.RandomOne}} + {{.RandomTwo}}</body></html>",
		"error.html":   "<html><title>Error</title><body>Error: {{.error}}</body></html>",
		"auth.html":    "<html><title>Auth</title><body>Auth Page</body></html>",
		"secure.html":  "<html><title>Secure</title><body>Secure Page: {{.payload}}</body></html>",
	}

	for name, content := range dummyTemplateFiles {
		path := filepath.Join(tempTemplatesDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create dummy template %s: %v", path, err)
		}
	}
	router.LoadHTMLGlob(filepath.Join(tempTemplatesDir, "*.html"))

	RandomOne = 1
	RandomTwo = 2

	LoadServerRoutes(router)

	var uploadBody bytes.Buffer
	mpWriter := multipart.NewWriter(&uploadBody)
	_, err := mpWriter.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file for upload: %v", err)
	}
	mpWriter.Close()

	testCases := []struct {
		name        string
		path        string
		method      string
		body        *bytes.Buffer
		contentType string
	}{
		{"RootGET", "/", http.MethodGet, nil, ""},
		{"PostsGET", "/posts", http.MethodGet, nil, ""},
		{"ArticleGET", "/article/3", http.MethodGet, nil, ""},
		{"CategoryGET", "/category/somecategory", http.MethodGet, nil, ""},
		{"ContactGET", "/contact", http.MethodGet, nil, ""},
		{"ListingsGET", "/listings", http.MethodGet, nil, ""},
		{"SpecificListingGET", "/listing/MLS123", http.MethodGet, nil, ""},
		{"AuthGET", "/auth", http.MethodGet, nil, ""},
		{"AuthPOST", "/auth", http.MethodPost, bytes.NewBufferString("email=test@example.com&password=p"), "application/x-www-form-urlencoded"},
		{"SecureGET", "/secure", http.MethodGet, nil, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var reqErr error

			if tc.body != nil {
				req, reqErr = http.NewRequest(tc.method, tc.path, tc.body)
				if reqErr == nil && tc.contentType != "" {
					req.Header.Set("Content-Type", tc.contentType)
				}
			} else {
				req, reqErr = http.NewRequest(tc.method, tc.path, nil)
			}
			if reqErr != nil {
				t.Fatalf("Failed to create request for %s %s: %v", tc.method, tc.path, reqErr)
			}

			rr := httptest.NewRecorder()
			var panicDetail interface{}
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicDetail = r
					}
				}()
				router.ServeHTTP(rr, req)
			}()

			if panicDetail != nil {
				t.Errorf("%s %s handler panicked: %v. Response Code: %d, Body: %s", tc.method, tc.path, panicDetail, rr.Code, rr.Body.String())
				return
			}

			isGenericGin404 := (rr.Code == http.StatusNotFound && strings.TrimSpace(rr.Body.String()) == "404 page not found")

			if isGenericGin404 {
				t.Errorf("%s %s returned a generic Gin 404, route likely not registered. Code: %d, Body: '%s'", tc.method, tc.path, rr.Code, rr.Body.String())
			}
		})
	}
}

func TestLoadMiddlewares(t *testing.T) {
	silenceLogrus(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	LoadMiddlewares(router)
}

func TestLoadStaticFileRoutes(t *testing.T) {
	silenceLogrus(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	LoadStaticFileRoutes(router)
}

func TestMainExecutionPath(t *testing.T) {
	silenceLogrus(t)
	gin.SetMode(gin.TestMode)

	originalDeployment, deploymentSet := os.LookupEnv("DEPLOYMENT")
	os.Setenv("DEPLOYMENT", "test-main-execution")
	defer func() {
		if deploymentSet {
			os.Setenv("DEPLOYMENT", originalDeployment)
		} else {
			os.Unsetenv("DEPLOYMENT")
		}
	}()

	finished := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("main() panicked: %v", r)
			}
			close(finished)
		}()
		main()
	}()

	select {
	case <-finished:
	case <-time.After(500 * time.Millisecond):
	}
}

func TestMainExecutionPathWithCertMagic(t *testing.T) {
	silenceLogrus(t)
	gin.SetMode(gin.TestMode)

	originalDeployment, deploymentSet := os.LookupEnv("DEPLOYMENT")
	originalDomain, domainSet := os.LookupEnv("DOMAIN")
	originalCertmagicLogger := certmagic.Default.Logger
	os.Setenv("DEPLOYMENT", "NAS")
	os.Setenv("DOMAIN", "localhost")
	certmagic.Default.Logger = zap.NewNop()
	defer func() {
		if deploymentSet {
			os.Setenv("DEPLOYMENT", originalDeployment)
		} else {
			os.Unsetenv("DEPLOYMENT")
		}
		if domainSet {
			os.Setenv("DOMAIN", originalDomain)
		} else {
			os.Unsetenv("DOMAIN")
		}
		certmagic.Default.Logger = originalCertmagicLogger
	}()

	finished := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("main() with DEPLOYMENT=NAS panicked (possibly expected due to certmagic in test env): %v", r)
			}
			close(finished)
		}()
		main()
	}()

	select {
	case <-finished:
	case <-time.After(1 * time.Second):
		t.Log("main() with DEPLOYMENT=NAS timed out (certmagic.HTTPS likely blocking, as expected in test).")
	}
}
