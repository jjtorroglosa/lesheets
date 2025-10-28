#!/bin/bash

file="$1"
if [ ! -f "$file" ]; then
    echo "File $file doesn't exist"
    exit 1
fi

cp "$file" "$file.tmp" && mv "$file.tmp" "$file"

echo "Updated $file"
