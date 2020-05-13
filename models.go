package main

import (
	"fmt"
	"html"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

//ContactForm : structure used to grab user data from /contact POST requests
type ContactForm struct {
	Name       string `json:"name" form:"name" binding:"required"`
	Email      string `json:"email" form:"email" binding:"required"`
	Website    string `json:"website" form:"website"`
	Message    string `json:"message" form:"message" binding:"required"`
	RobotCheck int    `json:"robot" form:"robot"`
}

//Item : structure used to get data from DynamoDB requests
type Item struct {
	ArticlePicture string `json:"article-picture"`
	Author         string `json:"author"`
	Categories     string `json:"categories"`
	CreatedDate    string `json:"created-date"`
	Excerpt        string `json:"excerpt"`
	HTMLHold       string `json:"html-hold"`
	ModifiedDate   string `json:"modified-date"`
	PanelPicture   string `json:"panel-picture"`
	PostID         int    `json:"post-id"`
	PostTitle      string `json:"post-title"`
	ShortTitle     string `json:"short-title"`
	PostType       string `json:"post-type"`
}

//Article : structure used to make DynamoDB data functional
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

//Category : structure used to access data in HTML Templates
type Category struct {
	Category string `json:"category"`
}

//Listing : structure used to make DynamoDB data functional
type Listing struct {
	MLS          string   `json:"MLS"`
	Street1      string   `json:"Street1"`
	Street2      string   `json:"Street2"`
	City         string   `json:"City"`
	State        string   `json:"State"`
	ZipCode      string   `json:"Zip Code"`
	Neighborhood string   `json:"Neighborhood"`
	SalesPrice   string   `json:"Sales Price"`
	DateListed   string   `json:"Date Listed"`
	LastModified string   `json:"Last Modified"`
	Bedrooms     string   `json:"Bedrooms"`
	ListPhoto    string   `json:"List Photo"`
	PhotoArray   []string `json:"Photo Array"`
	Bathrooms    string   `json:"Bathrooms"`
	GarageSize   string   `json:"Garage Size"`
	SquareFeet   string   `json:"Square Feet"`
	LotSize      string   `json:"Lot Size"`
	Description  string   `json:"Description"`
	User         string   `json:"User"`
	Deleted      string   `json:"deleted"`
}

// Return a list of all the article panels for the Front Page
func getArticlePanels() []Article {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	dbSvc := dynamodb.New(sess)

	filt := expression.Name("post-id").GreaterThanEqual(expression.Value(0))

	proj := expression.NamesList(expression.Name("post-title"), expression.Name("post-id"), expression.Name("post-type"),
		expression.Name("author"), expression.Name("categories"), expression.Name("excerpt"),
		expression.Name("modified-date"), expression.Name("panel-picture"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Live-Articles"),
	}

	// Make the DynamoDB Query API call
	result, err := dbSvc.Scan(params)

	articles := []Article{}

	for _, i := range result.Items {
		item := Item{}
		article := Article{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
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

// Return a list of all the article panels for the Category Pages
func getCategoryPageArticlePanels(category string) []Article {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	dbSvc := dynamodb.New(sess)

	unescapedCategory := html.UnescapeString(category)

	filt := expression.Name("categories").Contains(unescapedCategory)

	proj := expression.NamesList(expression.Name("post-title"), expression.Name("post-id"), expression.Name("post-type"),
		expression.Name("author"), expression.Name("categories"), expression.Name("excerpt"),
		expression.Name("modified-date"), expression.Name("panel-picture"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Live-Articles"),
	}

	// Make the DynamoDB Query API call
	result, err := dbSvc.Scan(params)

	articles := []Article{}

	for _, i := range result.Items {
		item := Item{}
		article := Article{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
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

func getArticleByID(id int) (*Article, error) {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dbSvc := dynamodb.New(sess)

	result, err := dbSvc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Live-Articles"),
		Key: map[string]*dynamodb.AttributeValue{
			"post-id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	item := Item{}
	article := Article{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
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

// Get a list of all the current realtor listings
func getRealtorListings() []Listing {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	dbSvc := dynamodb.New(sess)

	filt := expression.Name("deleted").NotEqual(expression.Value("anything"))

	proj := expression.NamesList(expression.Name("MLS"), expression.Name("Street1"), expression.Name("Street2"),
		expression.Name("City"), expression.Name("State"), expression.Name("Zip Code"), expression.Name("User"),
		expression.Name("Neighborhood"), expression.Name("Sales Price"), expression.Name("Date Listed"),
		expression.Name("Last Modified"), expression.Name("Bedrooms"), expression.Name("List Photo"),
		expression.Name("Photo Array"), expression.Name("Bathrooms"), expression.Name("Garage Size"),
		expression.Name("Square Feet"), expression.Name("Lot Size"), expression.Name("Description"),
		expression.Name("deleted"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Listings"),
	}

	// Make the DynamoDB Query API call
	result, err := dbSvc.Scan(params)

	listings := []Listing{}

	for _, i := range result.Items {
		listing := Listing{}

		err = dynamodbattribute.UnmarshalMap(i, &listing)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		listings = append(listings, listing)
	}

	sort.Slice(listings[:], func(i, j int) bool {
		return listings[i].LastModified > listings[j].LastModified
	})

	return listings
}

// Get a current realtor listing
func getRealtorListing(listing string) []Listing {
	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	dbSvc := dynamodb.New(sess)

	filt := expression.Name("MLS").Equal(expression.Value(listing))

	proj := expression.NamesList(expression.Name("MLS"), expression.Name("Street1"), expression.Name("Street2"),
		expression.Name("City"), expression.Name("State"), expression.Name("Zip Code"), expression.Name("User"),
		expression.Name("Neighborhood"), expression.Name("Sales Price"), expression.Name("Date Listed"),
		expression.Name("Last Modified"), expression.Name("Bedrooms"), expression.Name("List Photo"),
		expression.Name("Photo Array"), expression.Name("Bathrooms"), expression.Name("Garage Size"),
		expression.Name("Square Feet"), expression.Name("Lot Size"), expression.Name("Description"),
		expression.Name("deleted"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Listings"),
	}

	// Make the DynamoDB Query API call
	result, err := dbSvc.Scan(params)

	listings := []Listing{}

	for _, i := range result.Items {
		listing := Listing{}

		err = dynamodbattribute.UnmarshalMap(i, &listing)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		listings = append(listings, listing)
	}

	sort.Slice(listings[:], func(i, j int) bool {
		return listings[i].LastModified > listings[j].LastModified
	})

	return listings
}
