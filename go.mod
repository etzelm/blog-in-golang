module github.com/etzelm/blog-in-golang

go 1.14

require (
	github.com/aws/aws-sdk-go v1.23.0
	github.com/caddyserver/certmagic v0.10.12
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-contrib/static v0.0.0-20180301030858-73da7037e716
	github.com/gin-gonic/gin v1.3.0
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/ugorji/go/codec v0.0.0-20181209151446-772ced7fd4c2 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
)

replace golang.org/x/crypto v0.0.0-20180222182404-49796115aa4b => github.com/golang/crypto v0.0.0-20180222182404-49796115aa4b

replace github.com/ugorji/go => github.com/ugorji/go/codec v1.1.7
