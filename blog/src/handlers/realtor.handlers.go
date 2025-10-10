package handlers

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/etzelm/blog-in-golang/src/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// createS3Client creates an S3 client with proper configuration
func createS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := createAWSConfig(ctx)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

// ListingsGETAPI : Gets All Realtor Listings
func ListingsGETAPI(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	listings := models.GetRealtorListings()

	// Call the JSON method of the Context to return the results
	c.JSON(200, listings)

}

// ListingGETAPI : Gets A Realtor Listing
func ListingGETAPI(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")

	if listing := c.Param("listing"); listing != "" {

		card := models.GetRealtorListing(listing)

		// Call the JSON method of the Context to return the results
		c.JSON(200, card)

	} else {

		empty := []byte(``)

		// Call the JSON method of the Context to 404
		c.JSON(404, empty)

	}

}

// ListingPOSTAPI : Saves the user's data in DynamoDB and displays static response.html
func ListingPOSTAPI(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")

	if key := c.Param("key"); key == "HowMuchDoesSecurityCost" {

		var listing models.Listing
		c.BindJSON(&listing)

		ctx := context.TODO()
		dbSvc, err := createDynamoDBClient(ctx)
		if err != nil {
			log.Error("Unable to create DynamoDB client:", err)
			c.HTML(
				http.StatusInternalServerError,
				"error.html",
				gin.H{
					"title": "500 Internal Server Error",
					"error": err.Error(),
				},
			)
			return
		}

		av, err := attributevalue.MarshalMap(listing)
		if err != nil {
			log.Error("Error marshalling listing:", err)
			c.HTML(
				http.StatusInternalServerError,
				"error.html",
				gin.H{
					"title": "500 Internal Server Error",
					"error": err.Error(),
				},
			)
			return
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Listings"),
		}

		_, err = dbSvc.PutItem(ctx, input)

		if err != nil {
			log.Error("Got error calling PutItem:")
			log.Error(err.Error())
			c.HTML(
				// Set the HTTP status to 500 (Internal Server Error)
				http.StatusInternalServerError,
				// Use the error.html template
				"error.html",
				// Pass the data that the page uses
				gin.H{
					"title": "500 Internal Server Error",
					"error": err.Error(),
				},
			)
			return
		}

		success := []byte(`{"status":"success"}`)

		// Call the JSON method of the Context to return the results
		c.JSON(200, success)

	} else {

		empty := []byte(``)

		// Call the JSON method of the Context to 404
		c.JSON(404, empty)

	}

}

// UploadImagePOSTAPI : Upload an image to S3
func UploadImagePOSTAPI(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")

	if user := c.Param("user"); user != "" {

		ctx := context.TODO()
		svc, err := createS3Client(ctx)
		if err != nil {
			log.Error("Unable to create S3 client:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		form, _ := c.MultipartForm()

		files := form.File["file"]

		for _, file := range files {

			f, err := file.Open()

			if err != nil {
				log.Error(err)
			}

			defer f.Close()

			size := file.Size
			buffer := make([]byte, size)

			f.Read(buffer)
			fileBytes := bytes.NewReader(buffer)
			fileType := http.DetectContentType(buffer)
			path := "/media/" + user + "/" + file.Filename
			params := &s3.PutObjectInput{
				Bucket:      aws.String("blog-in-golang"),
				Key:         aws.String(path),
				Body:        fileBytes,
				ContentType: aws.String(fileType),
			}
			resp, err := svc.PutObject(ctx, params)
			if err != nil {
				log.Error("Error uploading to S3:", err)
			} else {
				log.Debug("S3 upload response:", resp)
			}
		}

		empty := []byte(``)
		c.JSON(http.StatusOK, empty)

	} else {

		empty := []byte(``)

		// Call the JSON method of the Context to 404
		c.JSON(404, empty)

	}

}
