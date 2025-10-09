package models

import (
	"context"
	"fmt"
	"html"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
)

// ContactForm : structure used to grab user data from /contact POST requests
type ContactForm struct {
	Name       string `json:"name" form:"name" binding:"required"`
	Email      string `json:"email" form:"email" binding:"required"`
	Website    string `json:"website" form:"website"`
	Message    string `json:"message" form:"message" binding:"required"`
	RobotCheck int    `json:"robot" form:"robot"`
	RobotNum   int    `json:"number" form:"number"`
}

// Item : structure used to get data from DynamoDB requests
type Item struct {
	ArticlePicture string `dynamodbav:"article-picture"`
	Author         string `dynamodbav:"author"`
	Categories     string `dynamodbav:"categories"`
	CreatedDate    string `dynamodbav:"created-date"`
	Excerpt        string `dynamodbav:"excerpt"`
	HTMLHold       string `dynamodbav:"html-hold"`
	ModifiedDate   string `dynamodbav:"modified-date"`
	PanelPicture   string `dynamodbav:"panel-picture"`
	PostID         int    `dynamodbav:"post-id"`
	PostTitle      string `dynamodbav:"post-title"`
	ShortTitle     string `dynamodbav:"short-title"`
	PostType       string `dynamodbav:"post-type"`
}

// Article : structure used to make DynamoDB data functional
type Article struct {
	ArticlePicture template.HTML `json:"article-picture"`
	Author         template.HTML `json:"author"`
	Categories     []Category    `json:"categories"`
	CreatedDate    string        `json:"created-date"`
	Excerpt        template.HTML `json:"excerpt"`
	HTMLHold       template.HTML `json:"html-hold"`
	ModifiedDate   string        `json:"modified-date"`
	PanelPicture   template.HTML `json:"panel-picture"`
	PostID         int           `json:"post-id"`
	PostTitle      string        `json:"post-title"`
	ShortTitle     string        `json:"short-title"`
	PostType       string        `json:"post-type"`
}

// Category : structure used to access data in HTML Templates
type Category struct {
	Category string `json:"category"`
}

// createDynamoDBClient creates a DynamoDB client with proper configuration
func createDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")

	var cfg aws.Config
	var err error

	if aid != "" && key != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion("us-west-1"),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(aid, key, "")),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion("us-west-1"),
		)
	}

	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

// GetArticlePanels Return a list of all the article panels for the Front Page
func GetArticlePanels() []Article {
	ctx := context.TODO()

	dbSvc, err := createDynamoDBClient(ctx)
	if err != nil {
		log.Error("Unable to create DynamoDB client:", err)
		return []Article{}
	}

	filt := expression.Name("post-id").GreaterThanEqual(expression.Value(0))

	proj := expression.NamesList(expression.Name("post-title"), expression.Name("post-id"), expression.Name("post-type"),
		expression.Name("author"), expression.Name("categories"), expression.Name("excerpt"),
		expression.Name("modified-date"), expression.Name("panel-picture"))

	expr, _ := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	tableName := os.Getenv("ARTICLES")
	if tableName == "" {
		tableName = "Test-Articles"
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := dbSvc.Scan(ctx, params)
	if err != nil {
		log.Error("Failed to scan DynamoDB:", err)
		return []Article{}
	}

	articles := []Article{}

	for _, i := range result.Items {
		item := Item{}
		article := Article{}

		err := attributevalue.UnmarshalMap(i, &item)

		if err != nil {
			log.Error("Got error unmarshalling:")
			log.Error(err.Error())
			return []Article{}
		}

		categories := []Category{}
		for _, category := range strings.Split(item.Categories, ",") {
			categories = append(categories, Category{category})
		}

		article.Author = template.HTML(item.Author)
		article.Categories = categories
		article.Excerpt = template.HTML(item.Excerpt)
		article.ModifiedDate = item.ModifiedDate
		article.PanelPicture = template.HTML(item.PanelPicture)
		article.PostID = item.PostID
		article.PostTitle = item.PostTitle
		article.PostType = item.PostType
		articles = append(articles, article)
	}

	sort.Slice(articles[:], func(i, j int) bool {
		return articles[i].PostID > articles[j].PostID
	})

	return articles
}

// GetCategoryPageArticlePanels Return a list of all the article panels for the Category Pages
func GetCategoryPageArticlePanels(category string) []Article {
	ctx := context.TODO()

	dbSvc, err := createDynamoDBClient(ctx)
	if err != nil {
		log.Error("Unable to create DynamoDB client:", err)
		return []Article{}
	}

	unescapedCategory := html.UnescapeString(category)

	filt := expression.Name("categories").Contains(unescapedCategory)

	proj := expression.NamesList(expression.Name("post-title"), expression.Name("post-id"), expression.Name("post-type"),
		expression.Name("author"), expression.Name("categories"), expression.Name("excerpt"),
		expression.Name("modified-date"), expression.Name("panel-picture"))

	expr, _ := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	tableName := os.Getenv("ARTICLES")
	if tableName == "" {
		tableName = "Test-Articles"
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := dbSvc.Scan(ctx, params)
	if err != nil {
		log.Error("Failed to scan DynamoDB:", err)
		return []Article{}
	}

	articles := []Article{}

	for _, i := range result.Items {
		item := Item{}
		article := Article{}

		err := attributevalue.UnmarshalMap(i, &item)

		if err != nil {
			log.Error("Got error unmarshalling:")
			log.Error(err.Error())
			return []Article{}
		}

		categories := []Category{}
		for _, category := range strings.Split(item.Categories, ",") {
			categories = append(categories, Category{category})
		}

		article.Author = template.HTML(item.Author)
		article.Categories = categories
		article.Excerpt = template.HTML(item.Excerpt)
		article.ModifiedDate = item.ModifiedDate
		article.PanelPicture = template.HTML(item.PanelPicture)
		article.PostID = item.PostID
		article.PostTitle = item.PostTitle
		article.PostType = item.PostType
		articles = append(articles, article)
	}

	sort.Slice(articles[:], func(i, j int) bool {
		return articles[i].PostID < articles[j].PostID
	})

	return articles
}

// GetArticleByID gets an article from DDB by id number
func GetArticleByID(id int) (*Article, error) {
	ctx := context.TODO()

	dbSvc, err := createDynamoDBClient(ctx)
	if err != nil {
		log.Error("Unable to create DynamoDB client:", err)
		return nil, err
	}

	table := os.Getenv("ARTICLES")
	if table == "" {
		table = "Test-Articles"
	}
	result, err := dbSvc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"post-id": &types.AttributeValueMemberN{
				Value: strconv.Itoa(id),
			},
		},
	})

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("article with ID %d not found", id)
	}

	item := Item{}
	article := Article{}

	err = attributevalue.UnmarshalMap(result.Item, &item)

	if err != nil {
		log.Error("Failed to unmarshal Record:", err)
		return nil, fmt.Errorf("failed to unmarshal record: %v", err)
	}

	categories := []Category{}
	for _, category := range strings.Split(item.Categories, ",") {
		categories = append(categories, Category{category})
	}

	article.ArticlePicture = template.HTML(item.ArticlePicture)
	article.Author = template.HTML(item.Author)
	article.Categories = categories
	article.CreatedDate = item.CreatedDate
	article.HTMLHold = template.HTML(item.HTMLHold)
	article.ModifiedDate = item.ModifiedDate
	article.PostID = item.PostID
	article.PostTitle = item.PostTitle
	article.ShortTitle = item.ShortTitle
	article.PostType = item.PostType

	return &article, nil
}
