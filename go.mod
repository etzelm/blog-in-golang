module github.com/etzelm/blog-in-golang

go 1.14

require (
	github.com/aws/aws-sdk-go v1.33.0
	github.com/caddyserver/certmagic v0.10.13
	github.com/gin-contrib/cache v1.1.0
	github.com/gin-contrib/static v0.0.0-20191128031702-f81c604d8ac2
	github.com/gin-gonic/gin v1.7.7
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120 // indirect
	golang.org/x/sys v0.0.0-20200515095857-1151b9dac4a9 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace golang.org/x/crypto v0.0.0-20180222182404-49796115aa4b => github.com/golang/crypto v0.0.0-20180222182404-49796115aa4b
