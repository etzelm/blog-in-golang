package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/etzelm/blog-in-golang/daemon/articles/awsEMR"
	"github.com/etzelm/blog-in-golang/daemon/articles/googleSRE"
	"github.com/etzelm/blog-in-golang/daemon/articles/graphStore"
	"github.com/etzelm/blog-in-golang/daemon/articles/infraCode"
	"github.com/etzelm/blog-in-golang/daemon/articles/reactRealtor"
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
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(id, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Error(err)
		return
	}

	dbSvc := dynamodb.New(sess)
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Feedback"),
	}

	_, err = dbSvc.CreateTable(input)

	if err != nil {
		log.Error("Got error calling CreateTable:")
		log.Error(err.Error())
		return
	}

	result, err := dbSvc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Tables:")
	for _, table := range result.TableNames {
		log.Info(*table)
	}
}

func dropTable() {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(id, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:8000")})
	if err != nil {
		log.Error(err)
		return
	}
	dbSvc := dynamodb.New(sess)

	hold := "Test-Articles"
	input := &dynamodb.DeleteTableInput{TableName: &hold}

	_, err = dbSvc.DeleteTable(input)

	if err != nil {
		log.Error("Got error calling DeleteTable:")
		log.Error(err.Error())
	}

	result, err := dbSvc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Current Tables:")
	for _, table := range result.TableNames {
		log.Info(*table)
	}
}
