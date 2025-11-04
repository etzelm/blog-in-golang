# blog-in-golang

![Go](https://img.shields.io/badge/Go-1.24.6+-00ADD8.svg?style=flat&logo=go)
![React](https://img.shields.io/badge/React-19.2+-61DAFB.svg?style=flat&logo=react)
![Vite](https://img.shields.io/badge/Vite-7.1+-646CFF.svg?style=flat&logo=vite)
![Testing](https://img.shields.io/badge/Testing-Vitest-green.svg?style=flat&logo=vitest)
![Code Coverage](https://img.shields.io/badge/Code%20Coverage-%3E90%25-green.svg?style=flat&logo=vitest)
![AWS](https://img.shields.io/badge/AWS-ACM%20%7C%20CloudFront%20%7C%20DynamoDB%20%7C%20Route53%20%7C%20S3-FF9900.svg?style=flat&logo=amazonaws)
![GCP Compute Engine](https://img.shields.io/badge/GCP-Compute%20Engine-4285F4.svg?style=flat&logo=google-cloud)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

A modern full-stack web application featuring a blog platform and realtor listing service. Built with **Go** (Gin) for the backend API, **React with Vite** for the sample frontend application, and **AWS services** (DynamoDB/S3) for cloud storage. Includes comprehensive authentication, blog management, and real estate CRUD operations with search functionality.

## Table of Contents

- [Quick Start](#quick-start)
- [Features](#features)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)

## Architecture Overview

```text
┌───────────────────┐    ┌───────────────────┐    ┌─────────────────┐
│   Frontend(User)  │    │    Go Backend     │    │   AWS Services  │
│                   │    │   (NAS/GCP VMs)   │    │                 │
│ ┌───────────────┐ │    │ ┌───────────────┐ │    │ ┌─────────────┐ │
│ │   CloudFront  │ │    │ │               │ │    │ │  DynamoDB   │ │
│ │     (CDN)     │◄┼────┼►│  Gin Server   │◄┼────┼►│ • Articles  │ │
│ └───────────────┘ │    │ │               │ │    │ │ • Listings  │ │
│                   │    │ │ • Blog API    │ │    │ │ • Sessions  │ │
│ ┌───────────────┐ │    │ │ • Realtor API │ │    │ └─────────────┘ │
│ │   React SPA   │ │    │ │ • Auth API    │ │    │                 │
│ │    (Vite)     │ │    │ │ • Templates   │ │    │ ┌─────────────┐ │
│ │               │◄┼────┼►│               │◄┼────┼►│     S3      │ │
│ │ • Realtor App │ │    │ └───────────────┘ │    │ │ • Images    │ │
│ │ • Auth Flow   │ │    │                   │    │ │ • CSS/JS    │ │
│ │ • CRUD Ops    │ │    │ ┌───────────────┐ │    │ │ • PDFs      │ │
│ └───────────────┘ │    │ │   CertMagic   │ │    │ └─────────────┘ │
│        ▲          │    │ │   (Auto SSL)  │ │    │                 │
│        │          │    │ └───────────────┘ │    │                 │
└────────┼──────────┘    └───────────────────┘    └─────────────────┘
         │
         ▼
┌──────────────────┐
│   Google OAuth   │
│       API        │
└──────────────────┘
```

## Quick Start

**Prerequisites:** Go 1.24.6+, Node.js 18+, Docker, AWS Account

### Local Development

```bash
# Clone and setup
git clone https://github.com/etzelm/blog-in-golang.git
cd blog-in-golang
cp .env.example .env  # Edit with your AWS credentials

# Start frontend
cd realtor && yarn install && yarn dev  # http://localhost:3000

# Start backend (new terminal)
cd ../blog && go mod download && go run app.go  # http://localhost:8080
```

### Docker Deployment

```bash
# Set environment variables
export AWS_ACCESS_KEY_ID='your_key'
export AWS_SECRET_ACCESS_KEY='your_secret'
export GAPI='your_google_client_id.apps.googleusercontent.com'

# Production build
docker build --build-arg AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
  --build-arg AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
  --build-arg GAPI="${GAPI}" -t blog:develop -f blog/Dockerfile .

# Development build (keeps source for testing)
docker build --build-arg AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
  --build-arg AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
  --build-arg GAPI="${GAPI}" --target development \
  -t blog:develop-test -f blog/Dockerfile .

# Run with Docker Compose
docker compose --file blog/docker-compose.yml up -d

# Run tests in container
docker compose -f blog/docker-compose.yml run --rm blog-test
```

**Live Demo:** [mitchelletzel.com](https://mitchelletzel.com) | [Realtor App](https://mitchelletzel.com/realtor)

## Features

- **Blog Platform**: Go templates with DynamoDB storage and categorization
- **Realtor SPA**: React/Vite app with property listings, image uploads, and search
- **Authentication**: Google OAuth2 + custom Go auth with bcrypt
- **Performance**: Gin middleware, Gzip compression, multi-level caching
- **Security**: Request filtering, automatic HTTPS via CertMagic
- **Cloud Native**: AWS (DynamoDB, S3, CloudFront) + GCP deployment
- **CI/CD**: GitHub Actions with automated testing and deployment
- **Deployment**: Docker Compose on GCP/NAS with CloudFront CDN

## API Documentation

### Blog API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Homepage with article panels |
| `GET` | `/category/:name` | Articles by category |
| `GET` | `/article/:id` | Individual article view |
| `GET` | `/about` | About page |
| `GET` | `/contact` | Contact form |
| `POST` | `/contact` | Submit contact form |
| `GET` | `/auth` | Authentication page |
| `POST` | `/auth` | User authentication |
| `GET` | `/secure` | Protected content |

### Realtor API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/listings` | Get all property listings |
| `GET` | `/api/listing/:id` | Get specific listing |
| `POST` | `/api/listing` | Create/update listing |
| `POST` | `/api/upload` | Upload property images |

### Response Formats

```json
// Listing Response
{
  "id": "uuid-string",
  "address": "123 Main St",
  "price": 350000,
  "bedrooms": 3,
  "bathrooms": 2,
  "sqft": 1800,
  "images": ["image1.jpg", "image2.jpg"],
  "description": "Beautiful home...",
  "created_at": "2025-01-01T00:00:00Z"
}
```

## Testing

### Backend Tests (Go)

```bash
cd blog

# Run all tests with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run specific test package
go test -v ./src/handlers
```

### Frontend Tests (React/Vitest)

```bash
cd realtor

# Run all tests
yarn test

# Run tests with coverage
yarn test --coverage

# Run tests in watch mode
yarn test:watch
```

### Docker Testing

```bash
# Build test image
docker build --build-arg AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
  --build-arg AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
  --build-arg GAPI="${GAPI}" --target development \
  -t blog:develop-test -f blog/Dockerfile .

# Run all tests with coverage
docker compose -f blog/docker-compose.yml run --rm blog-test

# Run specific tests
docker run --rm -v $(pwd)/blog:/app -w /app \
  -e AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
  -e AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
  golang:1.24.6-alpine3.22 go test -run TestSpecificFunction -v ./src/handlers/
```

### Current Test Coverage

- **Backend (Go)**: 82.7% overall coverage
- **Frontend (React)**: 92.6% statement coverage
- **Integration Tests**: Automated via GitHub Actions

## Development

### Environment Setup

Create `.env` with your AWS credentials:

```bash
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
AWS_REGION=us-east-1
ARTICLES=Test-Articles          # Use Live-Articles for production
LISTINGS=Test-Listings          # Use Live-Listings for production
GAPI=your_google_client_id.apps.googleusercontent.com

# Optional (Production)
DISTRIBUTION_ID1=your_cloudfront_distribution_id
DISTRIBUTION_ID2=your_backup_distribution_id
```

### AWS Setup Requirements

1. **DynamoDB Tables**: Create `Test-Articles` and `Test-Listings` tables
2. **S3 Bucket**: For image storage (realtor listings)
3. **IAM User**: With DynamoDB and S3 permissions
4. **CloudFront**: Optional, for CDN and SSL termination

### Common Commands

```bash
# Backend
cd blog && go run app.go
go test -v -coverprofile=coverage.out ./...

# Frontend
cd realtor && yarn dev
yarn test --coverage

# Docker Development
docker compose -f blog/docker-compose.yml up --force-recreate -d  # Start services
docker compose -f blog/docker-compose.yml run --rm blog-test      # Run tests
docker compose down && docker system prune -af                   # Clean up

# Docker Troubleshooting
docker build --no-cache --build-arg AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
  --build-arg AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
  --build-arg GAPI="${GAPI}" -t blog:develop -f blog/Dockerfile .  # Force rebuild

echo "AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID:+SET}"               # Check env vars
echo "AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY:+SET}"
echo "GAPI: ${GAPI:+SET}"
```

## Troubleshooting

### Common Docker Issues

```bash
# Environment variables not set
echo "AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID:+SET}"  # Should show 'SET'

# Build cache issues
docker builder prune -af  # Clear build cache
docker build --no-cache ...  # Force rebuild

# Port conflicts
lsof -i :80  # Check what's using port 80
docker stop $(docker ps -aq) && docker compose down  # Stop all containers

# Test failures
docker compose -f blog/docker-compose.yml run --rm blog-test  # Run tests
```

## Contributing

1. Fork and create a feature branch
2. Make changes and add tests
3. Run test suite: `go test ./... && yarn test`
4. Submit a Pull Request

See [Contributing Guide](CONTRIBUTING.md) for details. Security issues: email [etzelm@live.com](mailto:etzelm@live.com).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

**Mitchell Etzel** • [etzelm@live.com](mailto:etzelm@live.com) • [@etzelm](https://github.com/etzelm) • [mitchelletzel.com](https://mitchelletzel.com)

 Issues & 💡 Features: Use GitHub Issues • 💬 Questions: [Contact form](https://mitchelletzel.com/contact)

---

**⭐ Star this project if you find it helpful!**
