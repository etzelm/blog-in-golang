package models

import (
	"os"
	"testing"
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

	if listings != nil {
		_ = []Listing(listings)
	}
}
