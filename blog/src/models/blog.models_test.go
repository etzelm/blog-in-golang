package models

import (
	"html/template"
	"io"
	"os"
	"testing"

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

func TestGetArticleByID_SuccessfulFetchAndMap(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	tableName := "Test-Articles"
	os.Setenv("ARTICLES", tableName)
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	articleIDToFetch := 0
	article, err := GetArticleByID(articleIDToFetch)

	if err != nil {
		t.Fatalf("GetArticleByID returned an error, but a successful fetch was expected. Error: %v. Ensure table '%s' exists and item ID %d is present with valid data (like from daemon/articles/graphStore).",
			err, tableName, articleIDToFetch)
	}

	assert.NoError(t, err, "Expected no error for a successful fetch of an existing, valid article")
	assert.NotNil(t, article, "Article should not be nil for a successful fetch")
	if article != nil {
		assert.Equal(t, articleIDToFetch, article.PostID, "PostID should match the requested ID")
		assert.Equal(t, "Scalable, Fault Tolerant, & Strongly Consistent Graph Store API", article.PostTitle, "PostTitle mismatch")
		assert.Equal(t, template.HTML("<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>"), article.Author, "Author mismatch")
		assert.Equal(t, "Fault Tolerant Graph Store API", article.ShortTitle, "ShortTitle mismatch")
		assert.Equal(t, "standard", article.PostType, "PostType mismatch")
		assert.Equal(t, "April 10th, 2018", article.CreatedDate, "CreatedDate mismatch")
		assert.Equal(t, "August 10th, 2019", article.ModifiedDate, "ModifiedDate mismatch")
		expectedCategories := []Category{
			{Category: "Distributed Systems"},
			{Category: "My Projects"},
		}
		assert.ElementsMatch(t, expectedCategories, article.Categories, "Categories mismatch")
		assert.NotEmpty(t, string(article.HTMLHold), "HTMLHold should be populated")
		assert.NotEmpty(t, string(article.ArticlePicture), "ArticlePicture should be populated")
	}
}
