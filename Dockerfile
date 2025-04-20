FROM golang:1.24.2-alpine3.21 AS builder

# Install git for dependency management
RUN apk add --no-cache git

# Copy go.mod and go.sum and download dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o blog-in-golang .

FROM alpine:3.21 AS runtime
RUN apk --no-cache add ca-certificates

# Create a non-root user and group
RUN addgroup -g 101 -S blog && adduser -u 101 -S blog -G blog

# Create and configure the certmagic directory for the blog user
RUN mkdir -p /home/blog/.local/share/certmagic && \
    chown -R blog:blog /home/blog/.local && \
    chmod -R 775 /home/blog/.local

WORKDIR /
COPY --from=builder /app/blog-in-golang .
COPY --from=builder /app/public/ /public/
COPY --from=builder /app/realtor/build/ /realtor/build/
COPY --from=builder /app/templates/ /templates/

# Switch to the blog user
USER blog

# Expose the application ports
EXPOSE 80
EXPOSE 443

# Set the entrypoint
ENTRYPOINT ["./blog-in-golang"]