package handlers // This test file is for the 'handlers' package

import (
	"net/http"
	"strings"
	"testing" // Import the standard Go testing package
)

// TestCheckPasswordHash verifies the CheckPasswordHash and HashPassword functions.
func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	otherPassword := "otherpassword"

	// Generate a hash for "password123" using your HashPassword function
	knownHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password '%s': %v", password, err)
	}

	// Generate a hash for "otherpassword" to use as a different valid hash
	differentValidHash, err := HashPassword(otherPassword)
	if err != nil {
		t.Fatalf("Failed to hash password '%s': %v", otherPassword, err)
	}

	malformedHash := "not_a_valid_bcrypt_hash"

	testCases := []struct {
		name           string
		testPassword   string // Renamed to avoid confusion with the 'password' variable above
		hashToTest     string // Renamed for clarity
		expectedResult bool
	}{
		{"CorrectPasswordCorrectHash", password, knownHash, true},
		{"IncorrectPasswordCorrectHash", "wrongpassword", knownHash, false},
		{"CorrectPasswordDifferentHash", password, differentValidHash, false}, // "password123" should not match hash of "otherpassword"
		{"CorrectPasswordMalformedHash", password, malformedHash, false},
		{"EmptyPasswordKnownHash", "", knownHash, false},
		{"PasswordEmptyHash", password, "", false},
		{"OtherPasswordCorrectHash", otherPassword, differentValidHash, true}, // Test "otherpassword" against its own hash
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the CheckPasswordHash function from the current (handlers) package
			result := CheckPasswordHash(tc.testPassword, tc.hashToTest)
			if result != tc.expectedResult {
				t.Errorf("CheckPasswordHash(%q, %q) = %v; want %v", tc.testPassword, tc.hashToTest, result, tc.expectedResult)
			}
		})
	}
}

func TestAuthPage_Simple(t *testing.T) {
	dummyTemplateContent := "<html><head><title>Test Auth</title></head><body>Auth Title: {{.title}}</body></html>"
	templateFileName := "auth.html"

	router, recorder, _ := setupTestRouterWithHTMLTemplate(t, templateFileName, dummyTemplateContent)

	router.GET("/auth", AuthPage)

	req, err := http.NewRequest(http.MethodGet, "/auth", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	expectedBodySubstring := "Auth Title: Simple Auth Page"
	if !strings.Contains(recorder.Body.String(), expectedBodySubstring) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstring, recorder.Body.String())
	}

	expectedCacheControl := "no-cache"
	actualCacheControl := recorder.Header().Get("Cache-Control")
	if actualCacheControl != expectedCacheControl {
		t.Errorf("Expected Cache-Control header %q, got %q", expectedCacheControl, actualCacheControl)
	}
}

func TestSecurePage_ValidCookie_Simple(t *testing.T) {
	dummyTemplateContent := "<html><head><title>Test Secure</title></head><body>User: {{.payload}}, IP: {{.ip}}</body></html>"
	templateFileName := "secure.html"

	router, recorder, _ := setupTestRouterWithHTMLTemplate(t, templateFileName, dummyTemplateContent)

	router.GET("/secure", SecurePage)

	testUser := "test@example.com"
	hashedUserToken, err := HashPassword(testUser)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, "/secure", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	req.RemoteAddr = "192.0.2.1:1234"
	req.AddCookie(&http.Cookie{Name: "user", Value: testUser})
	req.AddCookie(&http.Cookie{Name: "userToken", Value: hashedUserToken})

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d. Response body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	expectedBodySubstringUser := "User: " + testUser
	if !strings.Contains(recorder.Body.String(), expectedBodySubstringUser) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstringUser, recorder.Body.String())
	}

	expectedBodySubstringIP := "IP: 192.0.2.1"
	if !strings.Contains(recorder.Body.String(), expectedBodySubstringIP) {
		t.Errorf("Expected body to contain %q, got %q", expectedBodySubstringIP, recorder.Body.String())
	}

	expectedCacheControl := "no-cache"
	actualCacheControl := recorder.Header().Get("Cache-Control")
	if actualCacheControl != expectedCacheControl {
		t.Errorf("Expected Cache-Control header %q, got %q", expectedCacheControl, actualCacheControl)
	}
}
