FROM golang:1.23.1-alpine3.20 AS builder

# Need git for dep
RUN apk add --no-cache git

# Copy go.mod and go.sum and download dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy current dir (outside docker) to inside docker and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o blog-in-golang .

FROM alpine:3.20 AS runtime
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /app/blog-in-golang .
COPY --from=0 /app/public/ /public/
COPY --from=0 /app/realtor/build/ /realtor/build/
COPY --from=0 /app/templates/ /templates/

# Create a non-root user
RUN addgroup -S blog && adduser -S blog -G blog
USER blog

# Expose the application ports
EXPOSE 80
EXPOSE 443

# Set the entrypoint
ENTRYPOINT ["./blog-in-golang"]