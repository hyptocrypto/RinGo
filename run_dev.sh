#!/bin/bash

# Function to handle SIGINT
cleanup() {
    echo "Caught SIGINT. Cleaning up and exiting."
    pkill RinGo
    exit 1
}

# Register the cleanup function to be called on SIGINT
trap cleanup INT

# Kill any running server process
pkill RinGo

# Build and start the server
go run .
# Watch for changes in .go files
fswatch . | while read f; do 
    extension="${f##*.}"
    if [ "$extension" = "go" ]; then
        echo "Changes detected. Restarting server."
        pkill RinGo
        go run .
    fi
done
