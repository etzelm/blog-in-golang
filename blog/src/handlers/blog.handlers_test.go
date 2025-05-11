package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouterWithHTMLTemplate(t *testing.T, templateName, templateContent string) (*gin.Engine, *httptest.ResponseRecorder, string) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, templateName)

	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create dummy template %s: %v", templatePath, err)
	}

	router.LoadHTMLGlob(filepath.Join(tempDir, "*.html"))

	recorder := httptest.NewRecorder()

	return router, recorder, tempDir
}

func TestAboutPage_Simple(t *testing.T) {
	dummyTemplateContent := "<html><head><title>Test About</title></head><body>Content: {{.title}}</body></html>"
	templateFileName := "about.html"

	router, recorder, _ := setupTestRouterWithHTMLTemplate(t, templateFileName, dummyTemplateContent)

	router.GET("/about", AboutPage)

	req, err := http.NewRequest(http.MethodGet, "/about", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	expectedBodySubstring := "Content: Mitchell Etzel"
	if !strings.Contains(recorder.Body.String(), expectedBodySubstring) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstring, recorder.Body.String())
	}

	expectedCacheControl := "public, max-age=31536000"
	actualCacheControl := recorder.Header().Get("Cache-Control")
	if actualCacheControl != expectedCacheControl {
		t.Errorf("Expected Cache-Control header %q, got %q", expectedCacheControl, actualCacheControl)
	}
}
