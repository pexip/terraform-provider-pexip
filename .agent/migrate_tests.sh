#!/bin/bash

# Script to migrate all test files from terraform-plugin-sdk/v2 to terraform-plugin-testing

echo "Migrating all test files from Plugin SDK v2 to Plugin Framework testing..."

# Find all test files that still import terraform-plugin-sdk/v2
test_files=$(grep -l "terraform-plugin-sdk/v2" internal/provider/*_test.go)

total_files=$(echo "$test_files" | wc -l)
current_file=0

echo "Found $total_files test files to migrate:"

for file in $test_files; do
    current_file=$((current_file + 1))
    echo "[$current_file/$total_files] Migrating $file..."
    
    # Create backup
    cp "$file" "$file.bak"
    
    # Replace import statements
    sed -i '' 's|github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource|github.com/hashicorp/terraform-plugin-testing/helper/resource|g' "$file"
    
    # Remove os import if it's only used for TF_ACC
    # First check if os is used for anything other than Setenv("TF_ACC", "1")
    if grep -q 'os\.' "$file"; then
        # Check if os is used for anything other than TF_ACC
        if ! grep -v 'os\.Setenv("TF_ACC"' "$file" | grep -q 'os\.'; then
            # Only used for TF_ACC, remove import
            sed -i '' '/^[[:space:]]*"os"$/d' "$file"
        fi
    else
        # Not used at all, remove import
        sed -i '' '/^[[:space:]]*"os"$/d' "$file"
    fi
    
    # Remove TF_ACC environment variable setting
    sed -i '' '/_ = os\.Setenv("TF_ACC", "1")/d' "$file"
    sed -i '' '/os\.Setenv("TF_ACC", "1")/d' "$file"
    
    # Clean up empty lines that might have been left
    sed -i '' '/^[[:space:]]*$/N;/^\n$/d' "$file"
    
    echo "  Migrated $file"
done

echo "Migration complete! Cleaning up backup files..."
rm -f internal/provider/*_test.go.bak

echo "Running go mod tidy to update dependencies..."
go mod tidy

echo "All test files have been migrated to Plugin Framework testing!"