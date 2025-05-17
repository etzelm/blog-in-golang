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
