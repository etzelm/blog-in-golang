package handlers

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/etzelm/blog-in-golang/src/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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

		aid := os.Getenv("AWS_ACCESS_KEY_ID")
		key := os.Getenv("AWS_SECRET_ACCESS_KEY")
		var myCredentials = credentials.NewStaticCredentials(aid, key, "")

		sess, _ := session.NewSession(&aws.Config{
			Credentials: myCredentials,
			Region:      aws.String("us-west-1"),
			//Endpoint:    aws.String("http://localhost:8000"),
		})

		dbSvc := dynamodb.New(sess)

		av, _ := dynamodbattribute.MarshalMap(listing)

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Listings"),
		}

		_, err := dbSvc.PutItem(input)

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

		aid := os.Getenv("AWS_ACCESS_KEY_ID")
		key := os.Getenv("AWS_SECRET_ACCESS_KEY")
		creds := credentials.NewStaticCredentials(aid, key, "")

		cfg := aws.NewConfig().WithRegion("us-west-1").WithCredentials(creds)
		svc := s3.New(session.New(), cfg)

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
				Bucket:        aws.String("blog-in-golang"),
				Key:           aws.String(path),
				Body:          fileBytes,
				ContentLength: aws.Int64(size),
				ContentType:   aws.String(fileType),
			}
			resp, _ := svc.PutObject(params)

			log.Debug("response: ", awsutil.StringValue(resp))
		}

		empty := []byte(``)
		c.JSON(http.StatusOK, empty)

	} else {

		empty := []byte(``)

		// Call the JSON method of the Context to 404
		c.JSON(404, empty)

	}

}
