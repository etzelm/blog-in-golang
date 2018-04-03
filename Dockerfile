FROM golang:1.9.2-alpine3.7

# Need git for dep
RUN apk add --no-cache git

# Need dep to get dependencies
RUN go get github.com/golang/dep/cmd/dep

# Copy current dir (outside docker) to the proper directory (inside docker)
COPY . src/github.com/etzelm/blog-in-golang/

# Change current directory
WORKDIR src/github.com/etzelm/blog-in-golang/

# Get dependencies
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-in-golang .

# # Install assignment/code
# RUN go install .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/etzelm/blog-in-golang/blog-in-golang .
COPY . .
CMD ["./blog-in-golang"]  

EXPOSE 3000