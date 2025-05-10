#!/bin/bash

# Script to clean up files and directories ignored by Git

# Check if this is a Git repository
if [ ! -d .git ]; then
  echo "Error: This is not a Git repository. Please run this script from the root of your project."
  exit 1
fi

echo "Cleaning ignored files and directories..."

# The command to clean ignored files:
# -X: Remove only files ignored by Git.
# -d: Remove directories in addition to files.
# -f: Force the deletion (required by Git to prevent accidental data loss).
#
# If you want to see what would be deleted without actually deleting anything,
# you can replace '-f' with '-n' (dry-run).
# Example: git clean -Xdn
git clean -Xdf

echo "Cleanup complete."