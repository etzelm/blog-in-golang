module github.com/etzelm/blog-in-golang/daemon

go 1.26.5

replace github.com/etzelm/blog-in-golang => ../blog

require (
	github.com/aws/aws-sdk-go v1.55.8
	github.com/etzelm/blog-in-golang v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.10.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.43.5 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.35 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.35 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.60 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.8.60 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.37 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.63.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.36.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.12.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.36 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.5.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.5 // indirect
	github.com/aws/smithy-go v1.27.7 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)
