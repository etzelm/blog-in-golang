#!/bin/bash

# Test runner script for blog-in-golang
# Runs Go tests through Docker using the blog-test service

set -e  # Exit on any error

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    printf "${1}${2}${NC}\n"
}

# Function to show usage
usage() {
    print_color $BLUE "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --no-cache           Build Docker images without using cache"
    echo "  -h, --help           Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                   # Run tests with cached build"
    echo "  $0 --no-cache        # Run tests with clean build (no cache)"
}

# Parse command line arguments
no_cache_arg=""
while [[ $# -gt 0 ]]; do
    case $1 in
        --no-cache)
            no_cache_arg="--no-cache"
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

print_color $BLUE "🔨 Building Docker images first..."

# Change to project root
cd /Users/etzelm/Documents/GitHub/etzelm/blog-in-golang

# Build Docker images (only development needed for tests)
./blog/src/scripts/build-docker.sh --development $no_cache_arg 

print_color $BLUE "📊 Extracting Go coverage report from built image..."

# Extract and display the Go coverage report that was generated during build
docker run --rm blog:develop-test go tool cover -func=coverage.out

print_color $BLUE "📊 Extracting React coverage report from built image..."

# Extract and display the React coverage report that was saved during build
docker run --rm --entrypoint="" blog:develop-test sh -c "
echo '=== React Test Results ==='
if [ -f /app/react-test-results.txt ]; then 
    cat /app/react-test-results.txt
else 
    echo 'React test results not found'
fi
"

print_color $GREEN "✅ Coverage reports displayed!"