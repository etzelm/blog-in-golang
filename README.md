# blog-in-golang

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg?style=flat&logo=go)
![React](https://img.shields.io/badge/React-19.1+-61DAFB.svg?style=flat&logo=react)
![Testing](https://img.shields.io/badge/Testing-Vitest-green.svg?style=flat&logo=vitest)
![Code Coverage](https://img.shields.io/badge/Code%20Coverage-%3E90%25-green.svg?style=flat&logo=vitest)
![AWS](https://img.shields.io/badge/AWS-ACM%20%7C%20CloudFront%20%7C%20DynamoDB%20%7C%20Route53%20%7C%20S3-FF9900.svg?style=flat&logo=amazonaws)
![GCP Compute Engine](https://img.shields.io/badge/GCP-Compute%20Engine-4285F4.svg?style=flat&logo=google-cloud)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

A full-stack web application featuring a blog and realtor listing service. Built with **Go** (Gin) for the blog/backend, **React** for the sample frontend application, and **AWS DynamoDB/S3** for storage. Supports user authentication, blog post management, and real estate CRUD/search.

## Deployment

- Published at [mitchelletzel.com](https://mitchelletzel.com)
- Sample React app at [mitchelletzel.com/realtor](https://mitchelletzel.com/realtor)
- Deployed with Docker using `docker compose`
- Uses GitHub Actions (`.github/workflows/NAS-workflow.yml`) to:
  - Validate PRs with simple build/test commands
  - Build docker containers on push to `develop`/`master`
  - Deploy to NAS/GCP containers conditionally based on branch
  - Invalidate CloudFront caches
  - Test deployed URLs listed in 
    - `.github/workflows/public-urls.txt`
    - `.github/workflows/local-urls.txt`

## Environment Variables

The application and deployment scripts rely on several environment variables. Key variables include:

- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`: For AWS API access (DynamoDB, S3, CloudFront).
- `GAPI`: Your Google API Client ID for the React frontend authentication.
- `ARTICLES`: The DynamoDB table name for blog articles (e.g., `Live-Articles`, `Test-Articles`).
- `NAS_*`, `GCP_*`: Secrets for deploying to NAS and GCP environments (defined in GitHub Secrets).
- `DISTRIBUTION_ID1`, `DISTRIBUTION_ID2`: CloudFront distribution IDs for cache invalidation.

Ensure these are configured appropriately in your local environment and/or in your CI/CD settings.

## Sample Commands

- **Local Development & Testing**:
  - Install frontend dependencies: `cd realtor && yarn install`
  - Test frontend (using Vitest): `yarn test`
  - Build frontend: `yarn build`
  - Install backend dependencies: `cd ../blog && go mod download`
  - Run backend server: `go run app.go` 
  - Run article update daemon: `cd ../daemon && go run app.go 1`

- **Docker Helper Commands**:
  - Stop all containers: `docker stop $(docker ps -aq)`
  - Remove all containers: `docker rm $(docker ps -aq)`
  - Remove all images: `docker rmi --force $(docker images -q)`
  - Build image: `docker build --build-arg GAPI=$GAPI -t blog:<branch-name> -f blog/Dockerfile .`
  - Run container: `docker run -d -p 80:8080 blog:<tag>` 
  - Start with compose: `docker compose --file blog/docker-compose.yml up --force-recreate -d`
  - Stop with compose: `docker compose down`
  - Clean up unused Docker resources: `docker system prune -a -f`

## Features

- **Blog**: Create and categorize posts using Go templates, stored in DynamoDB.
- **Realtor**: React-based SPA to manage listings with image uploads and multi-parameter searches.
- **Auth**: Google OAuth2 integration in React frontend as well as a custom Go auth API implementation.
- **Performance**: Gin middleware for Gzip compression and caching (in-memory and HTTP headers).
- **Security**: Middleware to block malicious request paths. Automatic HTTPS via CertMagic in prod envs.
