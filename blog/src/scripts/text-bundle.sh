#!/bin/bash

# Script to combine multiple source code files into a single text file
# for the realtor React app and blog Go app

# Output file
OUTPUT_FILE="combined_code_snapshot.txt"
TIMESTAMP=$(date +"%Y-%m-%d %H:%M:%S")

# Initialize the output file with a header
echo "Combined Code Snapshot - Generated on $TIMESTAMP" > "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Function to append a file to the output with delimiters
append_file() {
    local file_path="$1"
    if [ -f "$file_path" ]; then
        echo "// Start of $file_path //" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
        cat "$file_path" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
        echo "// End of $file_path //" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    else
        echo "Warning: File $file_path not found" >&2
    fi
}

# List of files from realtor React app (adjusted for realtor/src directory)
REACT_FILES=(
    "realtor/src/components/__tests__/index.test.jsx"
    "realtor/src/components/__tests__/App.test.jsx"
    "realtor/src/components/__tests__/Home.test.jsx"
    "realtor/src/components/__tests__/Listing.test.jsx"
    "realtor/src/components/__tests__/Main.test.jsx"
    "realtor/src/components/__tests__/MyListing.test.jsx"
    "realtor/src/components/__tests__/MyListings.test.jsx"
    "realtor/src/components/__tests__/NavBar.test.jsx"
    "realtor/src/components/__tests__/Search.test.jsx"
    "realtor/src/components/__tests__/Tile.test.jsx"
    "realtor/src/components/MyListing.jsx"
    "realtor/src/components/Listing.jsx"
    "realtor/src/components/TileDeck.jsx"
    "realtor/src/components/NavBar.jsx"
    "realtor/src/components/Tile.jsx"
    "realtor/src/components/Home.jsx"
    "realtor/src/components/Main.jsx"
    "realtor/src/components/MyListings.jsx"
    "realtor/src/components/Search.jsx"
    "realtor/src/index.jsx"
    "realtor/src/App.jsx"
    "realtor/src/setupTests.jsx"
    "realtor/package.json"
    "realtor/vitest.config.js"
    "realtor/test-data/index.js"
    "realtor/test-data/listing1.json"
    "realtor/test-data/listing2.json"
    "realtor/test-data/listing3.json"
    "realtor/test-data/listing4.json"
    "realtor/test-data/listing5.json"
)

# List of files from blog Go app
GO_FILES=(
    ".github/docker-compose/github-runner/Dockerfile"
    ".github/workflows/NAS-workflow.yml"
    "blog/docker-compose.yml"
    "blog/Dockerfile"
    "blog/README.md"
    "blog/go.mod"
    "blog/app.go"
    "blog/app_test.go"
    "blog/src/models/realtor.models.go"
    "blog/src/models/realtor.models_test.go"
    "blog/src/models/blog.models.go"
    "blog/src/models/blog.models_test.go"
    "blog/src/models/auth.models.go"
    "blog/src/handlers/blog.handlers.go"
    "blog/src/handlers/blog.handlers_test.go"
    "blog/src/handlers/auth.handlers.go"
    "blog/src/handlers/auth.handlers_test.go"
    "blog/src/handlers/realtor.handlers.go"
    "daemon/articles/awsEMR/awsEMR.go"
    "daemon/articles/awsEMR/awsEMR.html"
    "daemon/articles/infraCode/infraCode.go"
    "daemon/articles/infraCode/infraCode.html"
    "daemon/articles/graphStore/graphStore.go"
    "daemon/articles/graphStore/graphStore.html"
    "daemon/articles/reactRealtor/reactRealtor.go"
    "daemon/articles/reactRealtor/reactRealtor.html"
    "daemon/articles/googleSRE/googleSRE.go"
    "daemon/articles/googleSRE/googleSRE.html"
    "daemon/app.go"
    "daemon/README.md"
    "README.md"
)

# Append React app files
echo "Appending React app files..." >&2
for file in "${REACT_FILES[@]}"; do
    append_file "$file"
done

# Append Go app files
echo "Appending Go app files..." >&2
for file in "${GO_FILES[@]}"; do
    append_file "$file"
done

echo "Code snapshot generated at $OUTPUT_FILE" >&2