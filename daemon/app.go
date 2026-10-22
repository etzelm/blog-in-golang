package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/etzelm/blog-in-golang/daemon/articles/awsEMR"
	"github.com/etzelm/blog-in-golang/daemon/articles/googleSRE"
	"github.com/etzelm/blog-in-golang/daemon/articles/graphStore"
	"github.com/etzelm/blog-in-golang/daemon/articles/infraCode"
	"github.com/etzelm/blog-in-golang/daemon/articles/reactRealtor"
	"github.com/etzelm/blog-in-golang/daemon/ddb"
	log "github.com/sirupsen/logrus"
)

func main() {
	//Switch Case Depends on Article ID Number
	switch args := os.Args[1:]; args[0] {
	case "0":
		log.Info("Pushing changes for GraphStore")
		graphStore.GraphStore()
	case "1":
		log.Info("Pushing changes for GoogleSRE")
		googleSRE.GoogleSRE()
	case "2":
		log.Info("Pushing changes for ReactRealtor")
		reactRealtor.ReactRealtor()
	case "3":
		log.Info("Pushing changes for AmazonEMR")
		awsEMR.AmazonEMR()
	case "4":
		log.Info("Pushing changes for InfraCode")
		infraCode.InfraCode()
	case "42":
		createTable()
	//Execute Order 66 Meme
	case "66":
		dropTable()
	default:
		fmt.Printf("No Input Given")
	}
}

func createTable() {
	ctx := context.Background()
	dbSvc, err := ddb.NewClient(ctx, "us-west-1", "http://localhost:8000")
	if err != nil {
		log.Error(err)
		return
	}
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Feedback"),
	}

	_, err = dbSvc.CreateTable(ctx, input)

	if err != nil {
		log.Error("Got error calling CreateTable:")
		log.Error(err.Error())
		return
	}

	result, err := dbSvc.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Tables:")
	for _, table := range result.TableNames {
		log.Info(table)
	}
}

func dropTable() {
	ctx := context.Background()
	dbSvc, err := ddb.NewClient(ctx, "us-west-2", "http://localhost:8000")
	if err != nil {
		log.Error(err)
		return
	}

	_, err = dbSvc.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String("Test-Articles")})

	if err != nil {
		log.Error("Got error calling DeleteTable:")
		log.Error(err.Error())
	}

	result, err := dbSvc.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Current Tables:")
	for _, table := range result.TableNames {
		log.Info(table)
	}
}
