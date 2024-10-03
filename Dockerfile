FROM golang:1.23.1-alpine3.20

# Need git for dep
RUN apk add --no-cache git

# Copy current dir (outside docker) to inside docker
COPY . /go/src/github.com/etzelm/blog-in-golang/
WORKDIR /go/src/github.com/etzelm/blog-in-golang/

# Build go server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o blog-in-golang .

FROM alpine:3.20 
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/etzelm/blog-in-golang/blog-in-golang .
COPY --from=0 /go/src/github.com/etzelm/blog-in-golang/public/ /public/
COPY --from=0 /go/src/github.com/etzelm/blog-in-golang/realtor/build/ /realtor/build/
COPY --from=0 /go/src/github.com/etzelm/blog-in-golang/templates/ /templates/
CMD ["./blog-in-golang"]

EXPOSE 80
EXPOSE 443