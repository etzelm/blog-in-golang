// Package ddb holds the DynamoDB client setup shared by the daemon's
// article pushers and table commands.
package ddb

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/etzelm/blog-in-golang/src/models"
	log "github.com/sirupsen/logrus"
)

// NewClient returns a DynamoDB client using the AWS_ACCESS_KEY_ID and
// AWS_SECRET_ACCESS_KEY environment credentials. An empty endpoint uses the
// AWS default; "http://localhost:8000" targets DynamoDB Local.
func NewClient(ctx context.Context, region, endpoint string) (*dynamodb.Client, error) {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(id, key, "")),
	)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	}), nil
}

// PutArticle writes item into the table named by the ARTICLES environment
// variable. It exits non-zero on failure so the CI article push fails loudly.
func PutArticle(name string, item models.Item) {
	ctx := context.Background()

	dbSvc, err := NewClient(ctx, "us-west-1", "")
	if err != nil {
		log.Fatal(err)
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling %s: %v", name, err)
	}

	log.Infof("Putting %s into DDB", name)
	_, err = dbSvc.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("ARTICLES")),
	})
	if err != nil {
		log.Fatalf("Got error calling PutItem for %s: %v", name, err)
	}
}
