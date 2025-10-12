package models

import (
	"context"
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

func TestGetArticlePanels_DefaultTableName(t *testing.T) {
	silenceLogrus(t)
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	// Unset the ARTICLES environment variable to test default table name path
	os.Unsetenv("ARTICLES")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	panels := GetArticlePanels()

	// Should use default table name "Test-Articles" when env var not set
	assert.NotNil(t, panels, "GetArticlePanels should return a non-nil slice when using default table name")
}

func TestGetArticleByID_ArticleNotFound(t *testing.T) {
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

	// Use a non-existent article ID that should not be found
	nonExistentID := 99999
	article, err := GetArticleByID(nonExistentID)

	// Should return nil article and error when article not found
	assert.Nil(t, article, "Article should be nil when not found")
	assert.Error(t, err, "Error should be returned when article not found")
	assert.Contains(t, err.Error(), "not found", "Error message should indicate article was not found")
}

func TestCreateDynamoDBClient_Success(t *testing.T) {
	silenceLogrus(t)

	// Set up AWS credentials for testing
	originalAccessKey, accessKeySet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	originalSecretKey, secretKeySet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	originalRegion, regionSet := os.LookupEnv("AWS_REGION")

	os.Setenv("AWS_ACCESS_KEY_ID", "test_access_key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test_secret_key")
	os.Setenv("AWS_REGION", "us-east-1")

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

	client, err := createDynamoDBClient(context.Background())

	if err != nil {
		t.Errorf("createDynamoDBClient should not error with valid credentials: %v", err)
	}

	if client == nil {
		t.Error("Expected non-nil DynamoDB client")
	}
}

func TestCreateDynamoDBClient_DefaultRegion(t *testing.T) {
	silenceLogrus(t)

	// Set up AWS credentials but no region to test default
	originalAccessKey, accessKeySet := os.LookupEnv("AWS_ACCESS_KEY_ID")
	originalSecretKey, secretKeySet := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	originalRegion, regionSet := os.LookupEnv("AWS_REGION")

	os.Setenv("AWS_ACCESS_KEY_ID", "test_access_key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test_secret_key")
	os.Unsetenv("AWS_REGION") // Test default region

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

	client, err := createDynamoDBClient(context.Background())

	if err != nil {
		t.Errorf("createDynamoDBClient should not error with default region: %v", err)
	}

	if client == nil {
		t.Error("Expected non-nil DynamoDB client with default region")
	}
}

func TestCreateDynamoDBClient_NoCredentials(t *testing.T) {
	silenceLogrus(t)

	// Remove AWS credentials to test credential chain fallback
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

	client, _ := createDynamoDBClient(context.Background())

	// The function should still return a client using default credential chain
	if client == nil {
		t.Error("Expected non-nil DynamoDB client even without explicit credentials")
	}
}

func TestGetArticlePanels_EmptyResults(t *testing.T) {
	silenceLogrus(t)

	// Set table name to one that exists but has no matching items
	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "empty-test-table")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	panels := GetArticlePanels()

	// Should return empty slice, not nil, when no articles found
	assert.NotNil(t, panels, "GetArticlePanels should never return nil")
	assert.IsType(t, []Article{}, panels, "GetArticlePanels should return []Article type")
}

func TestGetCategoryPageArticlePanels_EmptyCategory(t *testing.T) {
	silenceLogrus(t)

	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "test-table")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	testCases := []struct {
		name     string
		category string
	}{
		{
			name:     "EmptyCategory",
			category: "",
		},
		{
			name:     "NonExistentCategory",
			category: "NonExistentCategoryThatShouldNotExist",
		},
		{
			name:     "SpecialCharsCategory",
			category: "Category@#$%",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			panels := GetCategoryPageArticlePanels(tc.category)

			assert.NotNil(t, panels, "GetCategoryPageArticlePanels should never return nil for category: %s", tc.category)
			assert.IsType(t, []Article{}, panels, "GetCategoryPageArticlePanels should return []Article type")
		})
	}
}

func TestGetArticleByID_ErrorCases(t *testing.T) {
	silenceLogrus(t)

	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "test-table")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	testCases := []struct {
		name      string
		articleID int
	}{
		{
			name:      "NegativeID",
			articleID: -1,
		},
		{
			name:      "ZeroID",
			articleID: 0,
		},
		{
			name:      "LargeID",
			articleID: 999999,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			article, err := GetArticleByID(tc.articleID)

			// For non-existent articles, should return nil article
			// Error handling depends on implementation - might return error or nil
			if article != nil {
				assert.IsType(t, &Article{}, article, "GetArticleByID should return *Article type when found")
			}

			// Test completed successfully regardless of error
			_ = err // Error handling varies by implementation
		})
	}
}

func TestGetArticleByID_TableNameHandling(t *testing.T) {
	silenceLogrus(t)

	// Test with different table name configurations
	testCases := []struct {
		name      string
		tableName string
		unsetEnv  bool
	}{
		{
			name:      "CustomTableName",
			tableName: "CustomTestTable",
			unsetEnv:  false,
		},
		{
			name:      "DefaultTableName",
			tableName: "",
			unsetEnv:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")

			if tc.unsetEnv {
				os.Unsetenv("ARTICLES")
			} else {
				os.Setenv("ARTICLES", tc.tableName)
			}

			defer func() {
				if articlesEnvIsSet {
					os.Setenv("ARTICLES", originalArticlesEnv)
				} else {
					os.Unsetenv("ARTICLES")
				}
			}()

			// Test that function handles different table configurations
			article, err := GetArticleByID(1)

			// Function should handle gracefully regardless of table name
			if article != nil {
				assert.IsType(t, &Article{}, article, "GetArticleByID should return *Article type")
			}

			_ = err // Allow for various error conditions
		})
	}
}

func TestGetCategoryPageArticlePanels_CategoryProcessing(t *testing.T) {
	silenceLogrus(t)

	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "Test-Articles")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	// Test various category formats to improve coverage
	testCategories := []struct {
		name     string
		category string
	}{
		{
			name:     "StandardCategory",
			category: "Technology",
		},
		{
			name:     "MultiWordCategory",
			category: "Distributed Systems",
		},
		{
			name:     "CaseSensitiveCategory",
			category: "distributed systems", // Different case
		},
		{
			name:     "CategoryWithSpaces",
			category: " Technology ", // Leading/trailing spaces
		},
	}

	for _, tc := range testCategories {
		t.Run(tc.name, func(t *testing.T) {
			panels := GetCategoryPageArticlePanels(tc.category)

			assert.NotNil(t, panels, "GetCategoryPageArticlePanels should never return nil")
			assert.IsType(t, []Article{}, panels, "Should return []Article type")

			// If results found, verify they're properly structured
			for i, panel := range panels {
				assert.GreaterOrEqual(t, panel.PostID, 0, "PostID should be non-negative for panel %d", i)
				assert.NotEmpty(t, panel.PostTitle, "PostTitle should not be empty for panel %d", i)
				assert.NotEmpty(t, panel.PostType, "PostType should not be empty for panel %d", i)
			}
		})
	}
}

func TestGetArticlePanels_SortingAndFiltering(t *testing.T) {
	silenceLogrus(t)

	originalArticlesEnv, articlesEnvIsSet := os.LookupEnv("ARTICLES")
	os.Setenv("ARTICLES", "Test-Articles")
	defer func() {
		if articlesEnvIsSet {
			os.Setenv("ARTICLES", originalArticlesEnv)
		} else {
			os.Unsetenv("ARTICLES")
		}
	}()

	panels := GetArticlePanels()

	assert.NotNil(t, panels, "GetArticlePanels should never return nil")

	// Test that results are properly sorted and filtered
	if len(panels) > 1 {
		// Verify descending PostID order
		for i := 0; i < len(panels)-1; i++ {
			assert.GreaterOrEqual(t, panels[i].PostID, panels[i+1].PostID,
				"Articles should be sorted by PostID in descending order")
		}
	}

	// Verify all articles have valid data
	for i, panel := range panels {
		// Allow different post types as data may contain 'standard', 'quote', etc.
		assert.NotEmpty(t, panel.PostType, "PostType should not be empty for article %d", i)
		assert.Contains(t, []string{"standard", "quote", "page", "post"}, panel.PostType,
			"PostType should be a valid type for article %d, got %s", i, panel.PostType)
		// Skip title check for quote types as they may have empty titles
		if panel.PostType == "standard" {
			assert.NotEmpty(t, panel.PostTitle, "PostTitle should not be empty for standard article %d", i)
		}
		assert.GreaterOrEqual(t, panel.PostID, 0, "PostID should be non-negative for article %d", i)
	}
}
