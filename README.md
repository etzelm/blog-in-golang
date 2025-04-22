# blog-in-golang

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg?style=flat&logo=go)
![React](https://img.shields.io/badge/React-19.1+-61DAFB.svg?style=flat&logo=react)
![AWS](https://img.shields.io/badge/AWS-CloudFront%20%7C%20DynamoDB%20%7C%20S3-FF9900.svg?style=flat&logo=amazonaws)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

A full-stack web application featuring a blog and realtor listing service. Built with **Go** (Gin) for the blog/backend, **React** for the sample frontend application, and **AWS DynamoDB/S3** for storage. Supports user authentication, blog post management, and real estate CRUD/search.

## Deployment

- Published at [mitchelletzel.com](https://mitchelletzel.com)
- Sample React app at [mitchelletzel.com/realtor](https://mitchelletzel.com/realtor)
- Deployed with Docker using `docker compose`
- Uses GitHub Actions (`NAS-workflow.yml`) to:
  - Validate PRs with simple build/test commands
  - Build docker containers on push to `develop`/`master`
  - Deploy to NAS/GCP containers conditionally
  - Invalidate CloudFront caches
  - Test URLs in `public-urls.txt`/`local-urls.txt`

## Sample Commands

- **Local Test Commands**:
  - Install frontend deps: `cd realtor && yarn install`
  - Test frontend: `yarn test`
  - Build frontend: `yarn build`
  - Install backend deps: `cd ../ && go mod download`
  - Run backend: `go run app.go`
  - Run article daemon: `cd daemon && go run app.go 1`

- **Docker Helper Commands**:
  - Stop containers: `docker stop $(docker ps -aq)`
  - Remove containers: `docker rm $(docker ps -aq)`
  - Remove images: `docker rmi --force $(docker images -q)`
  - Build image: `docker build -t blog:develop .`
  - Run container: `docker run -d -p 8080:3000 blog`
  - Start with compose: `docker compose up --force-recreate -d`
  - Stop with compose: `docker compose down`
  - Clean up: `docker system prune -a`

## Features

- **Blog**: Create and categorize posts, stored in DynamoDB.
- **Realtor**: Manage listings with image uploads and multi-parameter searches.
- **Auth**: Google OAuth2 or custom go auth for secure sign-in.
- **Performance**: Caching, Gzip compression, and security middleware.

## Contributing

- Fork the repo
- Create a branch: `git checkout -b feature/your-feature`
- Commit changes: `git commit -m "Add feature"`
- Push to branch: `git push origin feature/your-feature`
- Open a pull request

## Contact

- **Mitchell Etzel**
- Email: [etzelm@live.com](mailto:etzelm@live.com)
- GitHub: [etzelm](https://github.com/etzelm)

## Support

- Submit feedback via the contact form at `/contact`
- Report issues on GitHub