package main

import (
	"net/http"
	"net/http/httptest"
	"testing" // Import the standard Go testing package

	"github.com/gin-gonic/gin"
)

// TestRandRange checks if the randRange function generates a number
// within the specified min and max boundaries.
func TestRandRange(t *testing.T) {
	// Define test cases
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

	// Iterate over each test case
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

// TestStaticCacheMiddleware tests the staticCacheMiddleware.
func TestStaticCacheMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode) // Set Gin to test mode to suppress debug prints
	router := gin.New()
	router.Use(staticCacheMiddleware())

	router.GET("/public/test.css", func(c *gin.Context) {
		c.String(http.StatusOK, "test css")
	})
	router.GET("/api/data", func(c *gin.Context) {
		c.String(http.StatusOK, "api data")
	})

	reqStatic, _ := http.NewRequest(http.MethodGet, "/public/test.css", nil)
	rrStatic := httptest.NewRecorder()
	router.ServeHTTP(rrStatic, reqStatic)

	if status := rrStatic.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for static path: got %v want %v",
			status, http.StatusOK)
	}
	expectedCacheControl := "public, max-age=31536000"
	if rrStatic.Header().Get("Cache-Control") != expectedCacheControl {
		t.Errorf("handler returned unexpected Cache-Control header for static path: got %v want %v",
			rrStatic.Header().Get("Cache-Control"), expectedCacheControl)
	}

	reqNonStatic, _ := http.NewRequest(http.MethodGet, "/api/data", nil)
	rrNonStatic := httptest.NewRecorder()
	router.ServeHTTP(rrNonStatic, reqNonStatic)

	if status := rrNonStatic.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for non-static path: got %v want %v",
			status, http.StatusOK)
	}
	if rrNonStatic.Header().Get("Cache-Control") == expectedCacheControl {
		t.Errorf("handler returned Cache-Control header for non-static path, but should not have: got %v",
			rrNonStatic.Header().Get("Cache-Control"))
	}
}

// TestUnauthorizedMiddleware tests the unauthorizedMiddleware.
func TestUnauthorizedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode) // Set Gin to test mode
	router := gin.New()
	router.Use(unauthorizedMiddleware())

	// Register handlers for paths that are expected to return 200 OK
	router.GET("/allowed/path", func(c *gin.Context) {
		c.String(http.StatusOK, "allowed")
	})
	router.GET("/another/safe/path", func(c *gin.Context) {
		c.String(http.StatusOK, "safe path allowed")
	})

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string // Optional: check body for 200 OK responses
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
				t.Errorf("handler returned wrong status code for path %s: got %v want %v",
					tc.path, status, tc.expectedStatus)
			}

			if tc.expectedStatus == http.StatusOK && tc.expectedBody != "" {
				if rr.Body.String() != tc.expectedBody {
					t.Errorf("handler returned unexpected body for allowed path %s: got %v want %v",
						tc.path, rr.Body.String(), tc.expectedBody)
				}
			}
		})
	}
}
