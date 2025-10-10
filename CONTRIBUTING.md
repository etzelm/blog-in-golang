# Contributing to blog-in-golang

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Contact](#contact)

## Getting Started

### Prerequisites

Ensure you have the following installed:

- Go 1.24.6+
- Node.js 18+ and Yarn
- Docker and Docker Compose
- Git

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

   ```bash
   git clone https://github.com/YOUR_USERNAME/blog-in-golang.git
   cd blog-in-golang
   ```

3. Add the upstream remote:

   ```bash
   git remote add upstream https://github.com/etzelm/blog-in-golang.git
   ```

### Environment Setup

1. Copy the environment template:

   ```bash
   cp .env.example .env
   ```

2. Fill in your AWS credentials and API keys in `.env`
3. Install dependencies:

   ```bash
   # Backend dependencies
   cd blog && go mod download
   
   # Frontend dependencies  
   cd ../realtor && yarn install
   ```

## Development Workflow

### Branch Strategy

- `master`: Production-ready code (auto-deploys to live site)
- `develop`: Development branch (auto-deploys to test environment)
- `feature/*`: Feature development branches
- `bugfix/*`: Bug fix branches
- `hotfix/*`: Critical production fixes

### Creating a Feature

1. **Sync with upstream:**

   ```bash
   git checkout develop
   git pull upstream develop
   ```

2. **Create a feature branch:**

   ```bash
   git checkout -b feature/your-amazing-feature
   ```

3. **Make your changes** following our code standards

4. **Test your changes:**

   ```bash
   # Backend tests
   cd blog && go test ./...
   
   # Frontend tests
   cd realtor && yarn test
   ```

5. **Commit with descriptive messages:**

   ```bash
   git add .
   git commit -m "Add feature: brief description of what it does"
   ```

6. **Push to your fork:**

   ```bash
   git push origin feature/your-amazing-feature
   ```

## Code Standards

### Go Backend (`/blog`)

- Follow Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Maintain test coverage above 80%
- Use structured logging with logrus
- Handle errors appropriately (don't ignore them)

### React Frontend (`/realtor`)

- Use functional components with hooks
- Follow React best practices
- Use PropTypes for type checking
- Maintain test coverage above 90%
- Use meaningful component and variable names
- Keep components small and focused

### General Guidelines

- Keep commits atomic and well-described
- Write self-documenting code
- Update documentation for new features
- Follow existing code patterns and style
- Use meaningful commit messages (see [Conventional Commits](https://www.conventionalcommits.org/))

## Testing

### Required Testing

All contributions must include appropriate tests:

**Backend (Go):**

```bash
cd blog
go test -v ./...                    # Run all tests
go test -v -cover ./...             # With coverage
go test -v ./src/handlers           # Specific package
```

**Frontend (React):**

```bash
cd realtor
yarn test                           # Run all tests
yarn test --coverage                # With coverage
yarn test:watch                     # Watch mode
```

### Test Requirements

- **New Features**: Must include unit tests
- **Bug Fixes**: Must include regression tests
- **API Changes**: Must update integration tests
- **Minimum Coverage**: 80% backend, 90% frontend

## Pull Request Process

### Before Submitting

1. **Rebase on latest develop:**

   ```bash
   git fetch upstream
   git rebase upstream/develop
   ```

2. **Run full test suite:**

   ```bash
   # Backend
   cd blog && go test ./... && go mod tidy
   
   # Frontend
   cd realtor && yarn test && yarn build
   ```

3. **Check for lint issues:**

   ```bash
   # Go formatting
   cd blog && gofmt -s -w .
   
   # React linting
   cd realtor && yarn lint
   ```

### PR Template

When creating a PR, include:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

### Review Process

1. **Automated Checks**: GitHub Actions will run tests and builds
2. **Code Review**: Maintainer will review code quality and architecture
3. **Testing**: Changes will be tested in development environment
4. **Approval**: PR needs approval from a maintainer
5. **Merge**: Maintainer will merge when ready

## Issue Reporting

### Bug Reports

Use the bug report template and include:

- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, browser, versions)
- Screenshots/logs if applicable

### Feature Requests

Use the feature request template and include:

- Clear description of the feature
- Use case and motivation
- Possible implementation approach
- Alternatives considered

### Security Issues

**Do not create public issues for security vulnerabilities.**
Email security issues to [etzelm@live.com](mailto:etzelm@live.com).

## Development Environment

### Useful Commands

```bash
# Full development setup
make setup                          # If Makefile exists

# Backend development
cd blog
go run app.go                       # Start server
go test -v ./...                    # Test
go mod tidy                         # Clean dependencies

# Frontend development  
cd realtor
yarn dev                            # Start dev server
yarn build                          # Production build
yarn test:watch                     # Test in watch mode

# Docker development
docker compose -f blog/docker-compose.yml up --build
```

### Debugging

- **Backend**: Use Go debugger or add debug logs with logrus
- **Frontend**: Use browser dev tools and React Developer Tools
- **Integration**: Check both frontend and backend logs

## Contact

- **Mitchell Etzel**
- Email: [etzelm@live.com](mailto:etzelm@live.com)
- GitHub: [@etzelm](https://github.com/etzelm)

## Support

- Submit feedback via the contact form at `/contact`
- Report issues using GitHub issue templates
- Join discussions in GitHub Discussions (if enabled)

---

By contributing, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).
