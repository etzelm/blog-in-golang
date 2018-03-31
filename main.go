package main

import (
	"net"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// cancelAll can cancel all currently running goroutines
var cancelAll context.CancelFunc

// BoolInt is an integer type alias to create an 'enum' type for responses
type BoolInt int

// BoolInt constants are used to represent an integer for a specific bool value
// such that False = 0, and True = 1
const (
	FALSE BoolInt = iota
	TRUE
)

const (
	port = ":50051"
)

// StatusMsg is a type alias to allow for 'enum' style type
type StatusMsg int

// The constants to represent a Status message
const (
	SUCCESS StatusMsg = iota
	ERROR
)

// Represents the string representation of the 'status' field in some
// responses
var statuses = [...]string{
	"success",
	"error",
}

// IP is this computers IP address
var IP string

var IP_PORT string

func init() {
	log.SetFormatter(&log.TextFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

}

func launchGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// Register reflection service on gRPC server.

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	log.Info("Server is starting...")
	go launchGrpcServer()
	log.Info("finished launchGrpcServer.")

	IP_PORT := os.Getenv("ip_port")
	if IP_PORT == "" {
		IP_PORT = "127.0.0.1:8080"
	}
	log.Info("IP_PORT: ", IP_PORT)

	server := gin.Default()
	log.WithField("server", server).Info("Default Gin server create.")
	LoadRoutes(server)
	server.Run(IP_PORT)

}

// LoadRoutes does exactly that... loads all routes for the server.
func LoadRoutes(server *gin.Engine) *gin.Engine {
	server.GET("/", LandingPage)

	server.GET("/hello", Hello)

	// All '/check' routes are grouped for convenience/clarity.
	check := server.Group("/check")
	{
		check.GET("", CheckGet)
		check.POST("", CheckPost)
		check.PUT("", CheckPut)
	}

	return server
}
