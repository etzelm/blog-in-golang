package handlers

import (
	"html/template"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/etzelm/blog-in-golang/src/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// AuthPage : Gets Authenticaton Details And Verify "Client"
func AuthPage(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"auth.html",
		// Pass the data that the page uses
		gin.H{
			"title": "Simple Auth Page",
		},
	)

}

// SecurePage : Display login details
func SecurePage(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")

	cookie1, _ := c.Cookie("userToken")
	cookie2, _ := c.Cookie("user")
	ip := c.ClientIP()

	if !CheckPasswordHash(cookie2, cookie1) {
		c.HTML(
			// Set the HTTP status to 401 (Unauthorized)
			http.StatusUnauthorized,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "401 (Unauthorized)",
			},
		)
		return
	}

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"secure.html",
		// Pass the data that the page uses
		gin.H{
			"title":   "Simple Secure Page",
			"payload": cookie2,
			"ip":      ip,
		},
	)

}

// AuthResponse : Saves the user's data in DynamoDB and displays static response.html
func AuthResponse(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	var form models.AuthForm
	c.Bind(&form)

	email := template.HTMLEscapeString(form.Email)
	if m, _ := regexp.MatchString("^[ a-zA-Z0-9]+(@[a-zA-Z0-9.]+)*$", email); !m {
		c.HTML(
			// Set the HTTP status to 401 (Unauthorized)
			http.StatusUnauthorized,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "401 (Unauthorized)",
				"error": "Email should match a standard format like etzelm@live.com",
			},
		)
		return
	}

	aid := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var myCredentials = credentials.NewStaticCredentials(aid, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: myCredentials,
		Region:      aws.String("us-west-1"),
		//Endpoint:    aws.String("http://localhost:8000"),
	})
	if err != nil {
		log.Error(err)
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

	dbSvc := dynamodb.New(sess)
	result, _ := dbSvc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Auth"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(form.Email),
			},
		},
	})

	authForm := models.AuthForm{}
	dynamodbattribute.UnmarshalMap(result.Item, &authForm)

	if CheckPasswordHash(form.Password, authForm.Password) {
		cipher, _ := HashPassword(form.Email)
		c.SetCookie("user", form.Email, 60*60*24, "/", "mitchelletzel.com", false, false)
		c.SetCookie("userToken", cipher, 60*60*24, "/", "mitchelletzel.com", false, false)
	} else {
		c.HTML(
			// Set the HTTP status to 400 (Bad Request)
			http.StatusBadRequest,
			// Use the error.html template
			"error.html",
			// Pass the data that the page uses
			gin.H{
				"title": "400 Client Error",
			},
		)
		return
	}

	c.Redirect(http.StatusFound, "/secure")

}

// https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
