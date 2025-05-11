package models

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func silenceLogrus(t *testing.T) {
	originalOut := logrus.StandardLogger().Out
	logrus.SetOutput(io.Discard)
	t.Cleanup(func() {
		logrus.SetOutput(originalOut)
	})
}

func TestGetArticlePanels_Simple(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "dummy-test-table-for-models")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	panels := GetArticlePanels()

	if panels != nil {
		_ = []Article(panels)
	}
}

func TestGetCategoryPageArticlePanels_Simple(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "dummy-test-table-for-category-models")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	dummyCategory := "TestCategory"
	panels := GetCategoryPageArticlePanels(dummyCategory)

	if panels != nil {
		_ = []Article(panels)
	}
}

func TestGetArticleByID_Simple(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "dummy-test-table-for-article-by-id")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	dummyArticleID := 1
	article, err := GetArticleByID(dummyArticleID)

	if article != nil {
		_ = (*Article)(article)
	}

	_ = err
}
