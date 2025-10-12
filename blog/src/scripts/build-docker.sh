#!/bin/bash

# Build script for blog-in-golang Docker images
# This script builds both the production and development Docker images
# that are commonly used in the blog-in-golang project.

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    printf "${1}${2}${NC}\n"
}

# Function to check if required environment variables are set
check_env_vars() {
    local missing_vars=()
    
    if [[ -z "${AWS_ACCESS_KEY_ID}" ]]; then
        missing_vars+=("AWS_ACCESS_KEY_ID")
    fi
    
    if [[ -z "${AWS_SECRET_ACCESS_KEY}" ]]; then
        missing_vars+=("AWS_SECRET_ACCESS_KEY")
    fi
    
    if [[ -z "${GAPI}" ]]; then
        missing_vars+=("GAPI")
    fi
    
    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        print_color $RED "❌ Missing required environment variables:"
        for var in "${missing_vars[@]}"; do
            print_color $RED "   - $var"
        done
        print_color $YELLOW "💡 Please set these variables before running the script:"
        print_color $YELLOW "   export AWS_ACCESS_KEY_ID='your_key'"
        print_color $YELLOW "   export AWS_SECRET_ACCESS_KEY='your_secret'"
        print_color $YELLOW "   export GAPI='your_google_client_id.apps.googleusercontent.com'"
        exit 1
    fi
}

# Function to build production image
build_production() {
    print_color $BLUE "🔨 Building production Docker image (blog:develop)..."
    
    cd /Users/etzelm/Documents/GitHub/etzelm/blog-in-golang
    
    local cache_args=""
    if [[ "$no_cache" == true ]]; then
        cache_args="--no-cache"
        print_color $YELLOW "   Using --no-cache flag"
    fi
    
    docker build \
        --build-arg AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
        --build-arg AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
        --build-arg GAPI="${GAPI}" \
        -t blog:develop \
        -f blog/Dockerfile \
        $cache_args \
        .
    
    if [[ $? -eq 0 ]]; then
        print_color $GREEN "✅ Production image built successfully: blog:develop"
    else
        print_color $RED "❌ Failed to build production image"
        exit 1
    fi
}

# Function to build development/test image
build_development() {
    print_color $BLUE "🔨 Building development Docker image (blog:develop-test)..."
    
    cd /Users/etzelm/Documents/GitHub/etzelm/blog-in-golang
    
    local cache_args=""
    if [[ "$no_cache" == true ]]; then
        cache_args="--no-cache"
        print_color $YELLOW "   Using --no-cache flag"
    fi
    
    docker build \
        --build-arg AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
        --build-arg AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
        --build-arg GAPI="${GAPI}" \
        --target development \
        -t blog:develop-test \
        -f blog/Dockerfile \
        $cache_args \
        .
    
    if [[ $? -eq 0 ]]; then
        print_color $GREEN "✅ Development image built successfully: blog:develop-test"
    else
        print_color $RED "❌ Failed to build development image"
        exit 1
    fi
}

# Function to show usage
usage() {
    print_color $BLUE "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -p, --production     Build only production image (blog:develop)"
    echo "  -d, --development    Build only development image (blog:develop-test)"
    echo "  -a, --all            Build both images (default)"
    echo "  --no-cache           Build without using Docker cache"
    echo "  -h, --help           Show this help message"
    echo ""
    echo "Environment variables required:"
    echo "  AWS_ACCESS_KEY_ID      Your AWS access key"
    echo "  AWS_SECRET_ACCESS_KEY  Your AWS secret key"
    echo "  GAPI                   Your Google API client ID"
    echo ""
    echo "Examples:"
    echo "  $0                     # Build both images with cache"
    echo "  $0 --production        # Build only production image"
    echo "  $0 --development       # Build only development image"
    echo "  $0 --no-cache          # Build both images without cache"
    echo "  $0 --production --no-cache  # Build production image without cache"
}

# Main script logic
main() {
    local build_production=false
    local build_development=false
    local no_cache=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -p|--production)
                build_production=true
                shift
                ;;
            -d|--development)
                build_development=true
                shift
                ;;
            -a|--all)
                build_production=true
                build_development=true
                shift
                ;;
            --no-cache)
                no_cache=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                print_color $RED "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    # If no specific options provided, build both images
    if [[ "$build_production" == false && "$build_development" == false ]]; then
        build_production=true
        build_development=true
    fi
    
    # Check environment variables
    check_env_vars
    
    print_color $GREEN "🚀 Starting Docker image build process..."
    print_color $BLUE "📋 Environment variables:"
    print_color $BLUE "   AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID:+SET}"
    print_color $BLUE "   AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY:+SET}"
    print_color $BLUE "   GAPI: ${GAPI:+SET}"
    echo ""
    
    # Build images based on options
    if [[ "$build_production" == true ]]; then
        build_production
        echo ""
    fi
    
    if [[ "$build_development" == true ]]; then
        build_development
        echo ""
    fi
    
    print_color $GREEN "🎉 Build process completed successfully!"
    print_color $BLUE "📦 Built images:"
    
    if [[ "$build_production" == true ]]; then
        print_color $BLUE "   - blog:develop (production)"
    fi
    
    if [[ "$build_development" == true ]]; then
        print_color $BLUE "   - blog:develop-test (development)"
    fi
    
    echo ""
    print_color $YELLOW "💡 Next steps:"
    print_color $YELLOW "   Run with Docker Compose: docker compose --file blog/docker-compose.yml up -d"
    print_color $YELLOW "   Run tests: docker compose -f blog/docker-compose.yml run --rm blog-test"
}

# Run the main function with all arguments
main "$@"