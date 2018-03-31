package main

import (
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	port = ":50051"
)

// IP_PORT is this computers IP address
var IP_PORT string

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
