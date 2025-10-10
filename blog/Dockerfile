# Stage 1: Go Builder
FROM golang:1.24.6-alpine3.22 AS golang-builder

# Update package index and upgrade all packages to latest versions
RUN apk update && apk upgrade

# Install git for dependency management
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY blog/go.mod blog/go.sum ./
RUN go mod download

# Copy the rest of the Go source code
COPY blog/ .

# Test the Go application
# Declare ARGs for AWS creds
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
# Set environment variables and run tests in the same RUN command
RUN export AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID && \
    export AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY && \
    go test -v -coverprofile=coverage.out ./...
RUN go tool cover -func=coverage.out

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o blog-in-golang .

# Development stage - keeps source code for testing
FROM golang-builder AS development
# This stage has all the source code and dependencies for running tests

# Stage 2: React Builder
FROM node:20-alpine AS react-builder

# Declare ARG for Google API Key
ARG GAPI

# Set the working directory for the React app
WORKDIR /app/realtor

# Copy package.json and yarn.lock to leverage Docker cache
COPY realtor/package.json realtor/yarn.lock ./

# Install dependencies using yarn
RUN yarn install --frozen-lockfile

# Copy the rest of the React app source code
COPY realtor/ ./

# Replace placeholder with actual Google API Key using the build argument
# Ensure GAPI is set during docker build (e.g., --build-arg GAPI=$GAPI)
RUN sed -i "s/ThisIsSupposedToBeAnId/$GAPI/g" src/App.jsx
RUN sed -i "s/ThisIsSupposedToBeAnId/$GAPI/g" index.html

# Run the JS-level Unit Tests
RUN yarn test --coverage

# Build the React application
RUN yarn build

# Stage 3: Final Runtime Image
FROM alpine:3.22 AS runtime

# Update package index and upgrade all packages to latest versions for security
RUN apk update && apk upgrade

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create a non-root user and group for security
RUN addgroup -g 101 -S blog && adduser -u 101 -S blog -G blog

# Create and configure the certmagic directory for the blog user
# Ensure the non-root user owns the directory where certmagic stores data
RUN mkdir -p /home/blog/.local/share/certmagic && \
    chown -R blog:blog /home/blog/.local && \
    chmod -R 775 /home/blog/.local

# Set the final working directory
WORKDIR /app

# Copy the built Go binary from the golang-builder stage
COPY --from=golang-builder /app/blog-in-golang .

# Copy the static Go assets (public, templates) from the golang-builder stage
COPY --from=golang-builder /app/public/ /app/public/
COPY --from=golang-builder /app/templates/ /app/templates/

# Copy the built React app from the react-builder stage
COPY --from=react-builder /app/realtor/build/ /app/realtor/build/

# Ensure the final application directory is owned by the non-root user
RUN chown -R blog:blog /app

# Switch to the non-root user
USER blog

# Expose the necessary ports
EXPOSE 80
EXPOSE 443

# Set the entrypoint for the container
ENTRYPOINT ["./blog-in-golang"]