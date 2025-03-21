#!/bin/bash

# Use the first argument as the directory or default to current directory
directory="${1:-.}"

# Recursively find files starting with "rmn_" and not containing "proxy" in their names.
find "$directory" -type f -name "rmn_*" ! -name "*proxy*" | while IFS= read -r file
do
    # Check if the file does NOT contain "TRACE" (case insensitive)
    if ! grep -iq "TRACE" "$file"; then
        echo "$file"
    fi
done
