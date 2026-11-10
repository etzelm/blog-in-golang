module github.com/etzelm/blog-in-golang/daemon

go 1.27.3

replace github.com/etzelm/blog-in-golang => ../blog

require (
	github.com/aws/aws-sdk-go-v2 v1.47.8
	github.com/aws/aws-sdk-go-v2/config v1.33.15
	github.com/aws/aws-sdk-go-v2/credentials v1.20.15
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.21.19
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.73.1
	github.com/etzelm/blog-in-golang v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.10.2
)

require (
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.10.5 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.20.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.46.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.13.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.39.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.44.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.51.8 // indirect
	github.com/aws/smithy-go v1.28.8 // indirect
	golang.org/x/sys v0.49.0 // indirect
)
