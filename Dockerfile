FROM golang:1.14.2-alpine3.11

# Need git for dep
RUN apk add --no-cache git

# Copy current dir (outside docker) to the proper directory (inside docker)
COPY . src/github.com/etzelm/blog-in-golang/

# Change current directory
WORKDIR src/github.com/etzelm/blog-in-golang/

# Get dependencies
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-in-golang .

# # Install assignment/code
# RUN go install .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/etzelm/blog-in-golang/blog-in-golang .
COPY . .
CMD ["./blog-in-golang"]

EXPOSE 80
EXPOSE 443