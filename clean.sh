#!/bin/bash

# Define the specific directory name
TARGET_DIR="my-app"

# Check if the directory exists
if [ ! -d "$TARGET_DIR" ]; then
    echo "Directory '$TARGET_DIR' does not exist."
    exit 1
fi

# Confirm before deleting
read -p "Are you sure you want to delete the '$TARGET_DIR' directory? (y/N): " confirm
if [[ "$confirm" =~ ^[Yy]$ ]]; then
    rm -rf "$TARGET_DIR"
    echo "Directory '$TARGET_DIR' removed."
else
    echo "Operation canceled."
fi