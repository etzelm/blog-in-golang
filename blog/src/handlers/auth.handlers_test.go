package handlers // This test file is for the 'handlers' package

import (
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
