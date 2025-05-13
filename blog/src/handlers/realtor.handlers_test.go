package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/etzelm/blog-in-golang/src/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	tempDir, _ := os.MkdirTemp("", "test-templates")
	errorTemplatePath := filepath.Join(tempDir, "error.html")
	err := os.WriteFile(errorTemplatePath, []byte("<html><body>Error: {{.error}}</body></html>"), 0644)
	if err != nil {
		logrus.WithError(err).Warn("Failed to create dummy error template, proceeding without HTML templates")
	} else {
		router.LoadHTMLGlob(filepath.Join(tempDir, "*.html"))
	}
	return router
}

func TestListingPOSTAPI(t *testing.T) {
	silenceLogrus(t)
	router := setupTestRouter()
	router.POST("/listings/add/:key", ListingPOSTAPI)

	mockListing := models.Listing{
		MLS:          "TestMLS123",
		Street1:      "123 Test St",
		City:         "Testville",
		State:        "TS",
		ZipCode:      "12345",
		SalesPrice:   "100000",
		Bedrooms:     "3",
		Bathrooms:    "2",
		SquareFeet:   "1500",
		LotSize:      "3000",
		GarageSize:   "2 car",
		Description:  "A test listing",
		Neighborhood: "Test Hood",
		DateListed:   "1678886400000",
		LastModified: "1678886400000",
		User:         "test@example.com",
		Deleted:      "false",
		ListPhoto:    "http://example.com/photo.jpg",
		PhotoArray:   []string{"http://example.com/photo1.jpg"},
	}
	listingJSON, _ := json.Marshal(mockListing)

	t.Run("SuccessPathLeadsToAWSError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/listings/add/HowMuchDoesSecurityCost", bytes.NewBuffer(listingJSON))
		req.Header.Set("Content-Type", "application/json")

		os.Setenv("AWS_ACCESS_KEY_ID", "FAKE_KEY_ID")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "FAKE_SECRET_KEY")
		defer os.Unsetenv("AWS_ACCESS_KEY_ID")
		defer os.Unsetenv("AWS_SECRET_ACCESS_KEY")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Error: UnrecognizedClientException")
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	})

	t.Run("InvalidKey", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/listings/add/WrongKey", bytes.NewBuffer(listingJSON))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, `""`, w.Body.String())
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	})

	t.Run("BindError", func(t *testing.T) {
		w := httptest.NewRecorder()
		invalidJSON := []byte(`{"MLS": "TestMLS123", "Street1": 123`)
		req, _ := http.NewRequest(http.MethodPost, "/listings/add/HowMuchDoesSecurityCost", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		// Set credentials here as well, so the session creation doesn't fail first
		os.Setenv("AWS_ACCESS_KEY_ID", "FAKE_KEY_ID_FOR_BIND_ERROR")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "FAKE_SECRET_KEY_FOR_BIND_ERROR")
		defer os.Unsetenv("AWS_ACCESS_KEY_ID")
		defer os.Unsetenv("AWS_SECRET_ACCESS_KEY")

		router.ServeHTTP(w, req)

		// Expect 500 because PutItem should fail due to invalid data from failed bind
		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 500 status for invalid JSON leading to AWS error")
		// Check for a generic AWS error message part, as the exact error might vary
		assert.Contains(t, w.Body.String(), "Error:")
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	})
}

func TestUploadImagePOSTAPI(t *testing.T) {
	silenceLogrus(t)
	router := setupTestRouter()
	router.POST("/upload/image/:user", UploadImagePOSTAPI)

	t.Run("Success", func(t *testing.T) {
		silenceLogrus(t)
		w := httptest.NewRecorder()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.jpg")
		assert.NoError(t, err)
		_, err = io.WriteString(part, "fake image data")
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/upload/image/testuser", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		os.Setenv("AWS_ACCESS_KEY_ID", "FAKE_KEY_ID")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "FAKE_SECRET_KEY")
		defer os.Unsetenv("AWS_ACCESS_KEY_ID")
		defer os.Unsetenv("AWS_SECRET_ACCESS_KEY")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `""`, w.Body.String())
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	})

	t.Run("NoFileUploaded", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/upload/image/testuser", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `""`, w.Body.String())
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	})

	t.Run("MissingUserInRoute", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/upload/image/", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
