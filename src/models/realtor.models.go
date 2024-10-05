package models

import (
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// Listing : structure used to make DynamoDB data functional
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

// GetRealtorListings Get a list of all the current realtor listings
func GetRealtorListings() []Listing {
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

	expr, _ := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Listings"),
	}

	// Make the DynamoDB Query API call
	result, _ := dbSvc.Scan(params)

	listings := []Listing{}

	for _, i := range result.Items {
		listing := Listing{}

		err = dynamodbattribute.UnmarshalMap(i, &listing)

		if err != nil {
			log.Error("Got error unmarshalling:")
			log.Error(err.Error())
			return nil
		}

		listings = append(listings, listing)
	}

	sort.Slice(listings[:], func(i, j int) bool {
		return listings[i].LastModified > listings[j].LastModified
	})

	return listings
}

// GetRealtorListing Get a current realtor listing
func GetRealtorListing(listing string) []Listing {
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

	expr, _ := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Listings"),
	}

	// Make the DynamoDB Query API call
	result, _ := dbSvc.Scan(params)

	listings := []Listing{}

	for _, i := range result.Items {
		listing := Listing{}

		err = dynamodbattribute.UnmarshalMap(i, &listing)

		if err != nil {
			log.Error("Got error unmarshalling:")
			log.Error(err.Error())
			return nil
		}

		listings = append(listings, listing)
	}

	sort.Slice(listings[:], func(i, j int) bool {
		return listings[i].LastModified > listings[j].LastModified
	})

	return listings
}
