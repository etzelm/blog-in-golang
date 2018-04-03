package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// IP_PORT is this computers IP address
var IP_PORT string

func main() {
	log.Info("Server is starting...")

	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	IP_PORT := os.Getenv("ip_port")
	if IP_PORT == "" {
		IP_PORT = ":3000"
	}
	log.Info("IP_PORT: ", IP_PORT)
	log.Info("id: ", id)
	log.Info("key: ", key)

	var my_credentials = credentials.NewStaticCredentials(id, key, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: my_credentials,
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:8000")})
	if err != nil {
		log.Println(err)
		return
	}
	dbSvc := dynamodb.New(sess)

	result, err := dbSvc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Tables:")
	for _, table := range result.TableNames {
		log.Println(*table)
	}

	server := gin.Default()
	log.WithField("server", server).Info("Default Gin server create.")
	server.LoadHTMLGlob("templates/*")
	server.Use(static.Serve("/public", static.LocalFile("./public", true)))
	LoadRoutes(server)
	server.Run(IP_PORT)
}

// LoadRoutes does exactly that... loads all routes for the server.
func LoadRoutes(server *gin.Engine) *gin.Engine {
	server.GET("/", LandingPage)
	server.GET("/about", AboutPage)
	server.GET("/article/view/:article_id", getArticle)
	return server
}
