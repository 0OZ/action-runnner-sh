#!/bin/bash

# Script to build and deploy the runner-manager Go application
# This ensures the binary is available for cron jobs

RUNNER_DIR="/home/christian/actions-runner"
BINARY_NAME="runner-manager"

echo "Building runner-manager..."

# Change to the project directory
cd "$RUNNER_DIR" || exit 1

# Build the Go application
if make build; then
    echo "✅ Build successful!"
    
    # Copy the binary to the runner directory for cron access
    if cp "build/$BINARY_NAME" "$RUNNER_DIR/$BINARY_NAME"; then
        echo "✅ Binary deployed to $RUNNER_DIR/$BINARY_NAME"
        chmod +x "$RUNNER_DIR/$BINARY_NAME"
        
        # Test the binary
        if ./"$BINARY_NAME" --help >/dev/null 2>&1; then
            echo "✅ Binary is working correctly!"
        else
            echo "⚠️  Warning: Binary may not be working correctly"
        fi
    else
        echo "❌ Failed to copy binary to runner directory"
        exit 1
    fi
else
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "Runner-manager is ready for cron jobs!"
echo "You can now run ./crontab-setup.sh <github_token> to set up automated updates"
