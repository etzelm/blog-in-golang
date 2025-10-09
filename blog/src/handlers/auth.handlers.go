package handlers

import (
	"context"
	"html/template"
	"net/http"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/etzelm/blog-in-golang/src/models"
	"github.com/gin-gonic/gin"
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

	ctx := context.TODO()
	dbSvc, err := createDynamoDBClient(ctx)
	if err != nil {
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

	result, err := dbSvc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Auth"),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{
				Value: form.Email,
			},
		},
	})
	if err != nil {
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

	authForm := models.AuthForm{}
	err = attributevalue.UnmarshalMap(result.Item, &authForm)
	if err != nil {
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
