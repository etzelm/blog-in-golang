package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func silenceLogrus(t *testing.T) {
	originalOut := logrus.StandardLogger().Out
	logrus.SetOutput(io.Discard)
	t.Cleanup(func() {
		logrus.SetOutput(originalOut)
	})
}

func setupTestRouterWithHTMLTemplates(t *testing.T, templates map[string]string) (*gin.Engine, *httptest.ResponseRecorder, string) {
	silenceLogrus(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	tempDir := t.TempDir()

	for name, content := range templates {
		templatePath := filepath.Join(tempDir, name)
		if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create dummy template %s: %v", templatePath, err)
		}
	}
	router.LoadHTMLGlob(filepath.Join(tempDir, "*.html"))
	recorder := httptest.NewRecorder()
	return router, recorder, tempDir
}

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

func TestContactResponse_Success(t *testing.T) {
	silenceLogrus(t)
	dummyTemplates := map[string]string{
		"response.html": "<html><head><title>{{.title}}</title></head><body>Thank You For Contacting Us!</body></html>",
		"error.html":    "<html><head><title>{{.title}}</title></head><body>Error: {{.error}}</body></html>",
	}
	router, recorder, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)

	testRandomOne := 7
	testRandomTwo := 3
	expectedSum := testRandomOne + testRandomTwo

	originalAccessKeyID, accessKeyIDSet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	originalSecretKey, secretKeySet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	os.Setenv("AWS_ACCESS_KEY_ID", "DUMMY_AWS_ACCESS_KEY_ID")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "DUMMY_AWS_SECRET_ACCESS_KEY")

	defer func() {
		if accessKeyIDSet {
			os.Setenv("AWS_ACCESS_KEY_ID", originalAccessKeyID)
		} else {
			os.Unsetenv("AWS_ACCESS_KEY_ID")
		}
		if secretKeySet {
			os.Setenv("AWS_SECRET_ACCESS_KEY", originalSecretKey)
		} else {
			os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		}
	}()

	router.POST("/contact", ContactResponse(&testRandomOne, &testRandomTwo))

	formData := url.Values{}
	formData.Set("name", "Valid User")
	formData.Set("email", "valid@example.com")
	formData.Set("website", "example.com")
	formData.Set("message", "This is a valid message.")
	formData.Set("robot", "1")
	formData.Set("number", strconv.Itoa(expectedSum))

	req, err := http.NewRequest(http.MethodPost, "/contact", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError && recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d or %d for successful path (may fail at AWS); got %d. Response body: %s", http.StatusOK, http.StatusInternalServerError, recorder.Code, recorder.Body.String())
	}

	if recorder.Code == http.StatusOK {
		expectedBodySubstring := "Thank You For Contacting Us!"
		if !strings.Contains(recorder.Body.String(), expectedBodySubstring) {
			t.Errorf("Expected body to contain %q, got %q", expectedBodySubstring, recorder.Body.String())
		}
		expectedTitle := "Thank You!"
		if !strings.Contains(recorder.Body.String(), "<title>"+expectedTitle+"</title>") {
			t.Errorf("Expected title to contain %q, got %q", expectedTitle, recorder.Body.String())
		}
	} else if recorder.Code == http.StatusInternalServerError {
		if !strings.Contains(recorder.Body.String(), "Error:") {
			t.Errorf("Expected error page content for AWS failure, got: %s", recorder.Body.String())
		}
	}

	expectedCacheControl := "no-cache"
	actualCacheControl := recorder.Header().Get("Cache-Control")
	assert.Equal(t, expectedCacheControl, actualCacheControl, "Cache-Control header mismatch")
}

func TestContactResponse_RobotCheckFail(t *testing.T) {
	silenceLogrus(t)
	dummyTemplates := map[string]string{
		"error.html": "<html><head><title>{{.title}}</title></head><body>Error Message: {{.error}}</body></html>",
	}
	router, _, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)

	testRandomOne := 3
	testRandomTwo := 5
	expectedSum := testRandomOne + testRandomTwo

	router.POST("/contact", ContactResponse(&testRandomOne, &testRandomTwo))

	formData1 := url.Values{}
	formData1.Set("name", "Test User")
	formData1.Set("email", "test@example.com")
	formData1.Set("message", "Hello World")
	formData1.Set("robot", "0")
	formData1.Set("number", strconv.Itoa(expectedSum))

	req1, err1 := http.NewRequest(http.MethodPost, "/contact", strings.NewReader(formData1.Encode()))
	if err1 != nil {
		t.Fatalf("Couldn't create request: %v\n", err1)
	}
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder1 := httptest.NewRecorder()
	router.ServeHTTP(recorder1, req1)

	if recorder1.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d for robot=0; got %d. Response body: %s", http.StatusUnauthorized, recorder1.Code, recorder1.Body.String())
	}
	expectedBodySubstring1 := "Error Message: Don&#39;t be a robot please!"
	if !strings.Contains(recorder1.Body.String(), expectedBodySubstring1) {
		t.Errorf("Expected body for robot=0 to contain %q, got %q", expectedBodySubstring1, recorder1.Body.String())
	}
	expectedTitle1 := "401 (Unauthorized)"
	if !strings.Contains(recorder1.Body.String(), "<title>"+expectedTitle1+"</title>") {
		t.Errorf("Expected title for robot=0 to contain %q, got %q", expectedTitle1, recorder1.Body.String())
	}

	formData2 := url.Values{}
	formData2.Set("name", "Test User")
	formData2.Set("email", "test@example.com")
	formData2.Set("message", "Hello World")
	formData2.Set("robot", "1")
	formData2.Set("number", strconv.Itoa(expectedSum+1))

	req2, err2 := http.NewRequest(http.MethodPost, "/contact", strings.NewReader(formData2.Encode()))
	if err2 != nil {
		t.Fatalf("Couldn't create request: %v\n", err2)
	}
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder2 := httptest.NewRecorder()
	router.ServeHTTP(recorder2, req2)

	if recorder2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d for incorrect sum; got %d. Response body: %s", http.StatusUnauthorized, recorder2.Code, recorder2.Body.String())
	}
	expectedBodySubstring2 := "Error Message: Don&#39;t be a robot please!"
	if !strings.Contains(recorder2.Body.String(), expectedBodySubstring2) {
		t.Errorf("Expected body for incorrect sum to contain %q, got %q", expectedBodySubstring2, recorder2.Body.String())
	}
	expectedTitle2 := "401 (Unauthorized)"
	if !strings.Contains(recorder2.Body.String(), "<title>"+expectedTitle2+"</title>") {
		t.Errorf("Expected title for incorrect sum to contain %q, got %q", expectedTitle2, recorder2.Body.String())
	}
}

func TestPostPage_Simple(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "dummy-test-articles-table-for-postpage")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	dummyTemplates := map[string]string{
		"index.html": "<html><head><title>{{.title}}</title></head><body>{{range .payload}}<div>{{.PostTitle}}</div>{{else}}No Posts{{end}}</body></html>",
	}
	router, recorder, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)

	router.GET("/posts", PostPage)

	req, err := http.NewRequest(http.MethodGet, "/posts", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	expectedBodySubstringTitle := "<title>Blog Posts</title>"
	if !strings.Contains(recorder.Body.String(), expectedBodySubstringTitle) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstringTitle, recorder.Body.String())
	}

	if !strings.Contains(recorder.Body.String(), "No Posts") && !strings.Contains(recorder.Body.String(), "<div>") {
		t.Errorf("Expected 'No Posts' or post content, but found neither. Body: %s", recorder.Body.String())
	}

	expectedCacheControl := "public, max-age=31536000"
	actualCacheControl := recorder.Header().Get("Cache-Control")
	if actualCacheControl != expectedCacheControl {
		t.Errorf("Expected Cache-Control header %q, got %q", expectedCacheControl, actualCacheControl)
	}
}

func TestCategoryPage_ErrorOnNoPanels(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "dummy-test-articles-table-for-categoryerror")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	dummyTemplates := map[string]string{
		"index.html": "<html><head><title>{{.title}}</title></head><body>Category: {{.category}} {{if .payload}}Have Payload{{else}}No Payload{{end}}</body></html>",
		"error.html": "<html><head><title>{{.title}}</title></head><body>Error Details: {{.error}}</body></html>",
	}
	router, recorder, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)
	router.GET("/category/:category", CategoryPage)

	req, err := http.NewRequest(http.MethodGet, "/category/nonexistentcategory", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
	}

	expectedBodySubstring := "Error Details: Please provide a valid category"
	if !strings.Contains(recorder.Body.String(), expectedBodySubstring) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstring, recorder.Body.String())
	}
	expectedTitle := "404 (Not Found)"
	if !strings.Contains(recorder.Body.String(), "<title>"+expectedTitle+"</title>") {
		t.Errorf("Expected title to contain %q, got %q", expectedTitle, recorder.Body.String())
	}
}

func TestCategoryPage_RealPanels(t *testing.T) {
	silenceLogrus(t)
	os.Setenv("ARTICLES", "Test-Articles")

	dummyTemplates := map[string]string{
		"index.html": "<html><head><title>{{.title}}</title></head><body>Category: {{.category}} {{if .payload}}Have Payload{{else}}No Payload{{end}}</body></html>",
		"error.html": "<html><head><title>{{.title}}</title></head><body>Error Details: {{.error}}</body></html>",
	}
	router, recorder, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)
	router.GET("/category/:category", CategoryPage)

	req, err := http.NewRequest(http.MethodGet, "/category/NonExistentCategory", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
	}
}

func TestContactResponse_InvalidFormData(t *testing.T) {
	silenceLogrus(t)
	dummyTemplates := map[string]string{
		"error.html":    "<html><head><title>{{.title}}</title></head><body>Error Message: {{.error}}</body></html>",
		"response.html": "<html><head><title>{{.title}}</title></head><body>Success!</body></html>",
	}
	router, _, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)

	testRandomOne := 3
	testRandomTwo := 5

	router.POST("/contact", ContactResponse(&testRandomOne, &testRandomTwo))

	testCases := []struct {
		name           string
		formData       url.Values
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "MissingName",
			formData: url.Values{
				"email":   {"test@example.com"},
				"message": {"Hello"},
				"robot":   {"1"},
				"number":  {strconv.Itoa(testRandomOne + testRandomTwo)},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error Message: Invalid form data.",
		},
		{
			name: "MissingEmail",
			formData: url.Values{
				"name":    {"Test"},
				"message": {"Hello"},
				"robot":   {"1"},
				"number":  {strconv.Itoa(testRandomOne + testRandomTwo)},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error Message: Invalid form data.",
		},
		{
			name: "MissingMessage",
			formData: url.Values{
				"name":   {"Test"},
				"email":  {"test@example.com"},
				"robot":  {"1"},
				"number": {strconv.Itoa(testRandomOne + testRandomTwo)},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error Message: Invalid form data.",
		},
		{
			name: "InvalidName",
			formData: url.Values{
				"name":    {"!@#$"},
				"email":   {"test@example.com"},
				"message": {"Hello"},
				"robot":   {"1"},
				"number":  {strconv.Itoa(testRandomOne + testRandomTwo)},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Error Message: Name should contain only alphanumeric characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/contact", strings.NewReader(tc.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("Test %s: Expected status %d, got %d. Response body: %s", tc.name, tc.expectedStatus, recorder.Code, recorder.Body.String())
			}
			if !strings.Contains(recorder.Body.String(), tc.expectedBody) {
				t.Errorf("Test %s: Expected body to contain %q, got %q", tc.name, tc.expectedBody, recorder.Body.String())
			}
		})
	}
}

func TestArticlePage_InvalidArticleID(t *testing.T) {
	silenceLogrus(t)
	dummyTemplates := map[string]string{
		"error.html":   "<html><head><title>{{.title}}</title></head><body>Error: {{.error}}</body></html>",
		"article.html": "<html><head><title>{{.title}}</title></head><body>Article Content</body></html>",
	}
	router, _, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)

	router.GET("/article/:article_id", ArticlePage)

	testCases := []struct {
		name           string
		articleID      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "InvalidArticleIDFormat",
			articleID:      "invalid",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Error: Please provide a valid Article ID.",
		},
		{
			name:           "EmptyArticleID",
			articleID:      "",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/article/"+tc.articleID, nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("Test %s: Expected status %d, got %d. Response body: %s", tc.name, tc.expectedStatus, recorder.Code, recorder.Body.String())
			}
			if !strings.Contains(recorder.Body.String(), tc.expectedBody) {
				t.Errorf("Test %s: Expected body to contain %q, got %q", tc.name, tc.expectedBody, recorder.Body.String())
			}
		})
	}
}

func TestArticlePage_ArticleNotFound(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "dummy-test-articles-table-for-articlenotfound")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	dummyTemplates := map[string]string{
		"error.html":   "<html><head><title>{{.title}}</title></head><body>Error: {{.error}}</body></html>",
		"article.html": "<html><head><title>{{.title}}</title></head><body>Article Content</body></html>",
	}
	router, _, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)

	router.GET("/article/:article_id", ArticlePage)

	req, _ := http.NewRequest(http.MethodGet, "/article/9999", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent ID, got %d. Response body: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "Error: Please provide a valid Article ID.") {
		t.Errorf("Expected body to contain 'Error: Please provide a valid Article ID.', got %q", recorder.Body.String())
	}
}

func TestArticlePage_InvalidPostType(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "Test-Articles")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	dummyTemplates := map[string]string{
		"error.html":   "<html><head><title>{{.title}}</title></head><body>Error: {{.error}}</body></html>",
		"article.html": "<html><head><title>{{.title}}</title></head><body>Article Content</body></html>",
	}
	router, _, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)

	router.GET("/article/:article_id", ArticlePage)

	req, _ := http.NewRequest(http.MethodGet, "/article/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized && recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status %d (or %d if article 1 not found/misconfigured); got %d. Response body: %s", http.StatusUnauthorized, http.StatusNotFound, recorder.Code, recorder.Body.String())
	}

	if recorder.Code == http.StatusUnauthorized {
		if !strings.Contains(recorder.Body.String(), "Error: Please provide a valid Article ID.") {
			t.Errorf("Expected body for invalid PostType to contain 'Error: Please provide a valid Article ID.', got %q", recorder.Body.String())
		}
	} else if recorder.Code == http.StatusNotFound {
		if !strings.Contains(recorder.Body.String(), "Error: Please provide a valid Article ID.") {
			t.Errorf("Expected body for article not found to contain 'Error: Please provide a valid Article ID.', got %q", recorder.Body.String())
		}
	}
}

func TestAboutPage_Simple(t *testing.T) {
	silenceLogrus(t)
	dummyTemplateContent := "<html><head><title>About Page</title></head><body>About Us</body></html>"
	templateFileName := "about.html"

	router, recorder, _ := setupTestRouterWithHTMLTemplate(t, templateFileName, dummyTemplateContent)

	router.GET("/", AboutPage)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	expectedBodySubstring := "About Us"
	if !strings.Contains(recorder.Body.String(), expectedBodySubstring) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstring, recorder.Body.String())
	}

	expectedCacheControl := "public, max-age=31536000"
	actualCacheControl := recorder.Header().Get("Cache-Control")
	if actualCacheControl != expectedCacheControl {
		t.Errorf("Expected Cache-Control header %q, got %q", expectedCacheControl, actualCacheControl)
	}
}

func TestContactPage_Simple(t *testing.T) {
	silenceLogrus(t)
	dummyTemplateContent := "<html><head><title>Contact Page</title></head><body>Contact Us: {{.RandomOne}} + {{.RandomTwo}}</body></html>"
	templateFileName := "contact.html"

	router, recorder, _ := setupTestRouterWithHTMLTemplate(t, templateFileName, dummyTemplateContent)

	testRandomOne := 3
	testRandomTwo := 5

	router.GET("/contact", ContactPage(&testRandomOne, &testRandomTwo))

	req, err := http.NewRequest(http.MethodGet, "/contact", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	expectedBodySubstring := "Contact Us: 3 + 5"
	if !strings.Contains(recorder.Body.String(), expectedBodySubstring) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstring, recorder.Body.String())
	}

	expectedCacheControl := "no-cache"
	actualCacheControl := recorder.Header().Get("Cache-Control")
	if actualCacheControl != expectedCacheControl {
		t.Errorf("Expected Cache-Control header %q, got %q", expectedCacheControl, actualCacheControl)
	}
}

func TestCategoryPage_SuccessWithValidCategory(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "Test-Articles")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	dummyTemplates := map[string]string{
		"index.html": "<html><head><title>{{.title}}</title></head><body>Category: {{.category}} {{if .payload}}Have Payload{{else}}No Payload{{end}}</body></html>",
		"error.html": "<html><head><title>{{.title}}</title></head><body>Error Details: {{.error}}</body></html>",
	}
	router, recorder, _ := setupTestRouterWithHTMLTemplates(t, dummyTemplates)
	router.GET("/category/:category", CategoryPage)

	req, err := http.NewRequest(http.MethodGet, "/category/Distributed%20Systems", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	expectedBodySubstring := "Category: Distributed Systems"
	if !strings.Contains(recorder.Body.String(), expectedBodySubstring) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstring, recorder.Body.String())
	}

	expectedCacheControl := "public, max-age=31536000"
	actualCacheControl := recorder.Header().Get("Cache-Control")
	if actualCacheControl != expectedCacheControl {
		t.Errorf("Expected Cache-Control header %q, got %q", expectedCacheControl, actualCacheControl)
	}
}
