# Stage 1: Go Builder
# Pin the builder to the native build platform and cross-compile to the target
# arch (see the build RUN below) — fast, vs. emulating the whole toolchain.
FROM --platform=$BUILDPLATFORM golang:1.26.4-alpine3.22 AS golang-builder

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

# Test the Go application. The tests hit DynamoDB, so they need AWS creds —
# but those creds must not persist in the image's layer history. Mount them as
# BuildKit secrets (only visible during this single RUN, never written into a
# layer). Replaces the prior ARG/ENV approach, which fired the docker buildx
# `SecretsUsedInArgOrEnv` warnings on every CI build.
# CLI side (studio-workflow.yml):
#   --secret id=aws_access_key_id,env=AWS_ACCESS_KEY_ID
#   --secret id=aws_secret_access_key,env=AWS_SECRET_ACCESS_KEY
RUN --mount=type=secret,id=aws_access_key_id,env=AWS_ACCESS_KEY_ID,required=true \
    --mount=type=secret,id=aws_secret_access_key,env=AWS_SECRET_ACCESS_KEY,required=true \
    go test -v -coverprofile=coverage.out ./...
RUN go tool cover -func=coverage.out

# Build the Go application for the target platform. BuildKit populates
# TARGETOS/TARGETARCH automatically: linux/arm64 on the Studio's native build,
# linux/amd64 under `docker buildx build --platform linux/amd64` for the GCP
# tarball. CGO is disabled so this is a pure cross-compile.
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -a -installsuffix cgo -o blog-in-golang .

# Stage 2: React Builder
# JS build is arch-independent; keep it on the native platform for speed.
FROM --platform=$BUILDPLATFORM node:26-alpine AS react-builder

# node:N-alpine dropped the pre-bundled yarn around v22+; install yarn 1.x
# explicitly (matches the yarn 1.x lockfile format) so the existing yarn-based
# build commands still work without migrating to corepack/yarn 4.
RUN npm install -g corepack@latest && corepack enable

# Declare ARG for Google API Key
ARG GAPI

# Set the working directory for the React app
WORKDIR /app/realtor

# Copy package.json and yarn.lock to leverage Docker cache
COPY realtor/package.json realtor/yarn.lock realtor/.yarnrc.yml ./

# Install dependencies using yarn
RUN yarn install --immutable

# Copy the rest of the React app source code
COPY realtor/ ./

# Replace placeholder with actual Google API Key using the build argument
# Ensure GAPI is set during docker build (e.g., --build-arg GAPI=$GAPI)
RUN sed -i "s/ThisIsSupposedToBeAnId/$GAPI/g" src/App.jsx

# Save the coverage report for later extraction
RUN yarn test --coverage --reporter=verbose --silent > test-results.txt 2>&1
RUN cat test-results.txt

# Build the React application
RUN yarn build

# Development stage - keeps source code for testing
FROM golang-builder AS development
# This stage has all the source code and dependencies for running tests
# Copy React test results from react-builder stage
COPY --from=react-builder /app/realtor/test-results.txt /app/react-test-results.txt
COPY --from=react-builder /app/realtor/coverage/ /app/react-coverage/

# Stage 3: Final Runtime Image
FROM alpine:3.22 AS runtime

# Update package index and upgrade all packages to latest versions for security
# Combine RUN commands to reduce layers and improve build cache efficiency
RUN apk update && apk upgrade && \
    apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

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

# Copy the résumé data (single source for the about page + the distributed PDF)
# so the AboutPage handler can read data/resume.json at runtime.
COPY --from=golang-builder /app/data/ /app/data/

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
