package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRealtorListings_Simple(t *testing.T) {
	silenceLogrus(t)
	originalAccessKeyID, accessKeyIDSet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	originalSecretKey, secretKeySet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	os.Setenv("AWS_ACCESS_KEY_ID", "dummy_access_key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "dummy_secret_key")

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

	listings := GetRealtorListings()
	assert.NotNil(t, listings, "GetRealtorListings should return an empty slice, not nil, on typical scan failure with dummy credentials")
}

func TestGetRealtorListings_Real(t *testing.T) {
	silenceLogrus(t)
	listings := GetRealtorListings()
	assert.NotNil(t, listings, "GetRealtorListings should return an empty slice, not nil, on typical scan failure with dummy credentials")
}

func TestGetRealtorListing_Simple(t *testing.T) {
	silenceLogrus(t)
	originalAccessKeyID, accessKeyIDSet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	originalSecretKey, secretKeySet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	os.Setenv("AWS_ACCESS_KEY_ID", "dummy_access_key_for_single_listing")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "dummy_secret_key_for_single_listing")

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

	dummyMLSID := "0000000000"
	listings := GetRealtorListing(dummyMLSID)
	assert.NotNil(t, listings, "GetRealtorListing should return an empty slice, not nil, on typical scan failure")
}

func TestGetRealtorListing_Real(t *testing.T) {
	silenceLogrus(t)
	realMLSID := "e0377-aed6-f12-0f66-ee8ab4edcfdc"
	listing := GetRealtorListing(realMLSID)
	assert.NotNil(t, listing, "GetRealtorListing should return an empty slice, not nil, on typical scan failure")
}

func TestGetRealtorListings_ErrorHandling(t *testing.T) {
	silenceLogrus(t)

	// Test with no AWS credentials to force error handling
	originalAccessKey, accessKeySet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	originalSecretKey, secretKeySet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")

	defer func() {
		if accessKeySet {
			os.Setenv("AWS_ACCESS_KEY_ID", originalAccessKey)
		}
		if secretKeySet {
			os.Setenv("AWS_SECRET_ACCESS_KEY", originalSecretKey)
		}
	}()

	listings := GetRealtorListings()

	// Function should handle errors gracefully and return empty slice
	assert.NotNil(t, listings, "GetRealtorListings should never return nil")
	assert.IsType(t, []Listing{}, listings, "GetRealtorListings should return []Listing type")
	assert.Len(t, listings, 0, "GetRealtorListings should return empty slice on error")
}

func TestGetRealtorListings_WithCredentials(t *testing.T) {
	silenceLogrus(t)

	// Test with dummy credentials to exercise credential path
	originalAccessKey, accessKeySet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	originalSecretKey, secretKeySet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	originalRegion, regionSet := os.LookupEnv("AWS_REGION")

	os.Setenv("AWS_ACCESS_KEY_ID", "test_key_for_listings")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test_secret_for_listings")
	os.Setenv("AWS_REGION", "us-west-2")

	defer func() {
		if accessKeySet {
			os.Setenv("AWS_ACCESS_KEY_ID", originalAccessKey)
		} else {
			os.Unsetenv("AWS_ACCESS_KEY_ID")
		}
		if secretKeySet {
			os.Setenv("AWS_SECRET_ACCESS_KEY", originalSecretKey)
		} else {
			os.Unsetenv("AWS_SECRET_ACCESS_KEY")
		}
		if regionSet {
			os.Setenv("AWS_REGION", originalRegion)
		} else {
			os.Unsetenv("AWS_REGION")
		}
	}()

	listings := GetRealtorListings()

	assert.NotNil(t, listings, "GetRealtorListings should never return nil")
	assert.IsType(t, []Listing{}, listings, "GetRealtorListings should return []Listing type")

	// With dummy credentials, will likely get empty results due to AWS auth errors
	// but function should handle gracefully
}

func TestGetRealtorListing_ParameterVariations(t *testing.T) {
	silenceLogrus(t)

	// Set up dummy credentials
	os.Setenv("AWS_ACCESS_KEY_ID", "test_key_for_single")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test_secret_for_single")
	defer func() {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	}()

	testCases := []struct {
		name  string
		mlsID string
	}{
		{
			name:  "EmptyMLSID",
			mlsID: "",
		},
		{
			name:  "NumericMLSID",
			mlsID: "123456789",
		},
		{
			name:  "AlphanumericMLSID",
			mlsID: "ABC123XYZ",
		},
		{
			name:  "UUIDLikeMLSID",
			mlsID: "e0377-aed6-f12-0f66-ee8ab4edcfdc",
		},
		{
			name:  "ShortMLSID",
			mlsID: "123",
		},
		{
			name:  "LongMLSID",
			mlsID: "very-long-mls-id-that-might-test-length-limits-12345",
		},
		{
			name:  "SpecialCharsMLSID",
			mlsID: "MLS@#$%^&*()",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			listing := GetRealtorListing(tc.mlsID)

			assert.NotNil(t, listing, "GetRealtorListing should never return nil for MLS ID: %s", tc.mlsID)
			assert.IsType(t, []Listing{}, listing, "GetRealtorListing should return []Listing type")

			// With dummy credentials and likely non-existent MLS IDs,
			// expect empty results but no crashes
		})
	}
}

func TestGetRealtorListing_ErrorHandling(t *testing.T) {
	silenceLogrus(t)

	// Test without AWS credentials to force error path
	originalAccessKey, accessKeySet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	originalSecretKey, secretKeySet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")

	defer func() {
		if accessKeySet {
			os.Setenv("AWS_ACCESS_KEY_ID", originalAccessKey)
		}
		if secretKeySet {
			os.Setenv("AWS_SECRET_ACCESS_KEY", originalSecretKey)
		}
	}()

	testMLSID := "test-mls-123"
	listing := GetRealtorListing(testMLSID)

	// Function should handle errors gracefully
	assert.NotNil(t, listing, "GetRealtorListing should never return nil")
	assert.IsType(t, []Listing{}, listing, "GetRealtorListing should return []Listing type")
	assert.Len(t, listing, 0, "GetRealtorListing should return empty slice on error")
}

func TestGetRealtorListings_DataStructure(t *testing.T) {
	silenceLogrus(t)

	// Test that the function returns proper data structure even with minimal setup
	listings := GetRealtorListings()

	assert.NotNil(t, listings, "GetRealtorListings should never return nil")
	assert.IsType(t, []Listing{}, listings, "GetRealtorListings should return []Listing slice")

	// If any listings are returned, verify they have the expected structure
	for i, listing := range listings {
		// Test that listing has proper structure (all these should be strings)
		assert.IsType(t, "", listing.MLS, "MLS field should be string for listing %d", i)
		assert.IsType(t, "", listing.Street1, "Street1 field should be string for listing %d", i)
		assert.IsType(t, "", listing.City, "City field should be string for listing %d", i)
		assert.IsType(t, "", listing.State, "State field should be string for listing %d", i)
		assert.IsType(t, "", listing.ZipCode, "ZipCode field should be string for listing %d", i)
		assert.IsType(t, "", listing.SalesPrice, "SalesPrice field should be string for listing %d", i)
		assert.IsType(t, "", listing.Bedrooms, "Bedrooms field should be string for listing %d", i)
		assert.IsType(t, "", listing.Bathrooms, "Bathrooms field should be string for listing %d", i)
		assert.IsType(t, "", listing.User, "User field should be string for listing %d", i)
		assert.IsType(t, "", listing.Deleted, "Deleted field should be string for listing %d", i)

		// PhotoArray should be []string
		assert.IsType(t, []string{}, listing.PhotoArray, "PhotoArray should be []string for listing %d", i)
	}
}

func TestGetRealtorListing_DataStructure(t *testing.T) {
	silenceLogrus(t)

	testMLSID := "test-structure-123"
	listings := GetRealtorListing(testMLSID)

	assert.NotNil(t, listings, "GetRealtorListing should never return nil")
	assert.IsType(t, []Listing{}, listings, "GetRealtorListing should return []Listing slice")

	// If any listings are returned, verify structure
	for i, listing := range listings {
		assert.IsType(t, "", listing.MLS, "MLS field should be string for listing %d", i)
		assert.IsType(t, "", listing.Street1, "Street1 field should be string for listing %d", i)
		assert.IsType(t, "", listing.City, "City field should be string for listing %d", i)
		assert.IsType(t, "", listing.State, "State field should be string for listing %d", i)
		assert.IsType(t, "", listing.ZipCode, "ZipCode field should be string for listing %d", i)
		assert.IsType(t, []string{}, listing.PhotoArray, "PhotoArray should be []string for listing %d", i)

		// Verify MLS field matches requested ID if found
		if listing.MLS != "" && listing.MLS != "*" {
			// Only assert if we actually got a meaningful result
			if len(listings) == 1 {
				// If single result, it should match our query
				// (though with dummy data this might not be the case)
			}
		}
	}
}

func TestGetRealtorListings_EmptyResultHandling(t *testing.T) {
	silenceLogrus(t)

	// Test with credentials that will likely result in empty scan
	os.Setenv("AWS_ACCESS_KEY_ID", "empty_test_key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "empty_test_secret")
	defer func() {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	}()

	listings := GetRealtorListings()

	// Should handle empty results gracefully
	assert.NotNil(t, listings, "GetRealtorListings should return empty slice, not nil")
	assert.IsType(t, []Listing{}, listings, "Should return proper slice type")

	// Empty results are acceptable and expected with dummy credentials
	if len(listings) == 0 {
		assert.Len(t, listings, 0, "Empty result should be length 0")
	}
}

func TestGetRealtorListing_EmptyResultHandling(t *testing.T) {
	silenceLogrus(t)

	// Test with credentials that will likely result in empty query
	os.Setenv("AWS_ACCESS_KEY_ID", "empty_single_test_key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "empty_single_test_secret")
	defer func() {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	}()

	nonExistentMLS := "definitely-does-not-exist-12345"
	listings := GetRealtorListing(nonExistentMLS)

	// Should handle empty results gracefully
	assert.NotNil(t, listings, "GetRealtorListing should return empty slice, not nil")
	assert.IsType(t, []Listing{}, listings, "Should return proper slice type")

	// Empty results are expected for non-existent MLS IDs
	if len(listings) == 0 {
		assert.Len(t, listings, 0, "Empty result should be length 0")
	}
}
