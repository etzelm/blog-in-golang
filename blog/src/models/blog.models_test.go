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
	assert.NotNil(t, panels, "GetArticlePanels should return a non-nil slice")
}

func TestGetArticlePanels_DataProcessingAndSorting(t *testing.T) {
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

	panels := GetArticlePanels()

	assert.NotNil(t, panels, "GetArticlePanels returned nil, expected a slice of articles.")

	if len(panels) > 0 {
		// Test sorting - articles should be sorted by PostID in descending order
		for i := 0; i < len(panels)-1; i++ {
			assert.GreaterOrEqual(t, panels[i].PostID, panels[i+1].PostID,
				"Articles should be sorted by PostID in descending order. Article at index %d has PostID %d, next has %d",
				i, panels[i].PostID, panels[i+1].PostID)
		}

		// Test data processing - verify first article has expected fields populated
		firstArticle := panels[0]
		assert.NotEmpty(t, firstArticle.PostTitle, "PostTitle should be populated")
		assert.NotEmpty(t, firstArticle.PostType, "PostType should be populated")
		assert.NotEmpty(t, string(firstArticle.Author), "Author should be populated")
		assert.NotEmpty(t, firstArticle.Categories, "Categories should be populated")
		assert.NotEmpty(t, string(firstArticle.Excerpt), "Excerpt should be populated")
		assert.NotEmpty(t, firstArticle.ModifiedDate, "ModifiedDate should be populated")
		assert.NotEmpty(t, string(firstArticle.PanelPicture), "PanelPicture should be populated")
		assert.GreaterOrEqual(t, firstArticle.PostID, 0, "PostID should be non-negative")
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
	assert.NotNil(t, panels, "GetCategoryPageArticlePanels should return a non-nil slice even for a category with no articles")
}

func TestGetCategoryPageArticlePanels_DataProcessingAndSorting(t *testing.T) {
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

	categoryToTest := "Distributed Systems"
	expectedPostIDs := []int{0, 2, 4, 5}
	expectedNumberOfArticles := len(expectedPostIDs)

	panels := GetCategoryPageArticlePanels(categoryToTest)

	assert.NotNil(t, panels, "GetCategoryPageArticlePanels returned nil, expected a slice of articles.")
	assert.Len(t, panels, expectedNumberOfArticles, "Expected %d articles for category '%s', but got %d.", expectedNumberOfArticles, categoryToTest, len(panels))

	if len(panels) == expectedNumberOfArticles {
		for i, panel := range panels {
			assert.Equal(t, expectedPostIDs[i], panel.PostID, "Article at index %d has PostID %d, expected %d. Sorting might be incorrect.", i, panel.PostID, expectedPostIDs[i])
		}

		articleZero := panels[0]
		assert.Equal(t, 0, articleZero.PostID, "First article's PostID mismatch.")
		assert.Equal(t, "Scalable, Fault Tolerant, & Strongly Consistent Graph Store API", articleZero.PostTitle, "PostTitle mismatch for PostID 0.")
		assert.Equal(t, template.HTML("<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>"), articleZero.Author, "Author mismatch for PostID 0.")
		assert.Equal(t, "standard", articleZero.PostType, "PostType mismatch for PostID 0.")

		var actualCategories []string
		for _, cat := range articleZero.Categories {
			actualCategories = append(actualCategories, cat.Category)
		}
		expectedCategoriesForArticleZero := []string{"Distributed Systems", "My Projects"}
		assert.ElementsMatch(t, expectedCategoriesForArticleZero, actualCategories, "Categories mismatch for PostID 0.")
		assert.Len(t, articleZero.Categories, len(expectedCategoriesForArticleZero), "Incorrect number of categories for PostID 0.")
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

	assert.NoError(t, err, "GetArticleByID returned an error for PostID %d, but a successful fetch was expected. Ensure table '%s' exists and is populated.", articleIDToFetch, tableName)
	assert.NotNil(t, article, "Article should not be nil for a successful fetch of PostID %d from table '%s'.", articleIDToFetch, tableName)

	if article != nil {
		assert.Equal(t, articleIDToFetch, article.PostID, "PostID should match the requested ID.")
		assert.Equal(t, "Scalable, Fault Tolerant, & Strongly Consistent Graph Store API", article.PostTitle, "PostTitle mismatch.")
		assert.Equal(t, template.HTML("<a style=\"color:#9C6708;\" href=\"/\">Mitchell Etzel</a>"), article.Author, "Author mismatch.")
		assert.Equal(t, "Fault Tolerant Graph Store API", article.ShortTitle, "ShortTitle mismatch.")
		assert.Equal(t, "standard", article.PostType, "PostType mismatch.")
		assert.Equal(t, "April 10th, 2018", article.CreatedDate, "CreatedDate mismatch.")
		assert.Equal(t, "August 10th, 2019", article.ModifiedDate, "ModifiedDate mismatch.")

		expectedCategories := []Category{
			{Category: "Distributed Systems"},
			{Category: "My Projects"},
		}
		assert.ElementsMatch(t, expectedCategories, article.Categories, "Categories mismatch.")
		assert.NotEmpty(t, string(article.HTMLHold), "HTMLHold should be populated.")
		assert.NotEmpty(t, string(article.ArticlePicture), "ArticlePicture should be populated.")
	}
}

func TestGetArticlePanels_ErrorHandling(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	// Set to a non-existent table to trigger DynamoDB scan error
	os.Setenv("ARTICLES", "non-existent-table-that-should-cause-error")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	panels := GetArticlePanels()

	// When DynamoDB scan fails, function should return empty slice
	assert.NotNil(t, panels, "GetArticlePanels should return a non-nil slice even when scan fails")
	assert.Len(t, panels, 0, "GetArticlePanels should return empty slice when scan fails")
}
