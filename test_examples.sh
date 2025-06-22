#!/bin/bash
# Test script for Terminal DevTool
# This script demonstrates how to use the tool with sample commands

set -e  # Exit on error

echo "🧪 Terminal DevTool Test Script"
echo "==============================="
echo

# Check if FFmpeg is installed
if ! command -v ffmpeg &> /dev/null; then
    echo "❌ FFmpeg is not installed. Please install it first."
    echo "Visit https://ffmpeg.org/download.html for instructions."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install it first."
    echo "Visit https://golang.org/doc/install for instructions."
    exit 1
fi

# Check if Rust is installed
if ! command -v cargo &> /dev/null; then
    echo "❌ Rust is not installed. Please install it first."
    echo "Visit https://www.rust-lang.org/tools/install for instructions."
    exit 1
fi

# Create test directory if it doesn't exist
TEST_DIR="./test_media"
mkdir -p "$TEST_DIR"

# Generate a test video if it doesn't exist
TEST_VIDEO="$TEST_DIR/test_input.mp4"
if [ ! -f "$TEST_VIDEO" ]; then
    echo "🎬 Generating test video..."
    ffmpeg -f lavfi -i "testsrc=duration=5:size=1280x720:rate=30" -c:v libx264 -r 30 "$TEST_VIDEO" -y
    if [ $? -ne 0 ]; then
        echo "❌ Failed to generate test video. Please check FFmpeg installation."
        exit 1
    fi
fi

echo "✅ Test video created: $TEST_VIDEO"
echo

# Make sure the backend is built
echo "🔄 Building backend..."
cd backend
go build -o terminal-devtool-backend
if [ $? -ne 0 ]; then
    echo "❌ Failed to build backend"
    exit 1
fi

# Start the backend server in the background
echo "🚀 Starting backend server..."
./terminal-devtool-backend &
BACKEND_PID=$!
cd ..

# Wait for server to start and check it's running
echo "⏳ Waiting for server to start..."
sleep 3
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Backend server failed to start properly"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2> /dev/null || true
    fi
    exit 1
fi

echo "✅ Backend server started (PID: $BACKEND_PID)"
echo

# Build the frontend if it hasn't been built
echo "🔄 Building frontend..."
cd frontend-cli
cargo build
if [ $? -ne 0 ]; then
    echo "❌ Failed to build frontend"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2> /dev/null || true
    fi
    exit 1
fi

# Run the CLI commands
echo "📋 Running CLI examples..."
echo

# Error handling function
handle_error() {
    echo "❌ Command failed: $1"
    cd ../
    if [ ! -z "$BACKEND_PID" ]; then
        echo "🧹 Cleaning up server process..."
        kill $BACKEND_PID 2> /dev/null || true
    fi
    exit 1
}

# 1. Process video (using the backend)
echo "Example 1: Process video using backend"
cargo run -- process -i "../$TEST_VIDEO" -o "../$TEST_DIR/processed_backend.mp4" --resolution 640x360 --bitrate 500k || handle_error "Video processing with backend failed"
echo

# 2. Process video (locally)
echo "Example 2: Process video locally"
cargo run -- --local process -i "../$TEST_VIDEO" -o "../$TEST_DIR/processed_local.mp4" --resolution 320x240 || handle_error "Local video processing failed"
echo

# 3. Compare the results
echo "Example 3: Compare original and processed files"
cargo run -- compare -o "../$TEST_VIDEO" -p "../$TEST_DIR/processed_backend.mp4" || handle_error "Video comparison failed"
echo

# 4. Get media info
echo "Example 4: Get media info"
cargo run -- info -f "../$TEST_VIDEO" || handle_error "Media info retrieval failed"
cd ..
echo

# Clean up
echo "🧹 Cleaning up..."
if [ ! -z "$BACKEND_PID" ]; then
    kill $BACKEND_PID 2> /dev/null || true
    echo "✅ Backend server stopped"
fi

echo
echo "✅ All tests completed successfully!"
echo "The processed videos are in the $TEST_DIR directory."
echo
echo "📝 Summary:"
echo "- Original video: $TEST_VIDEO"
echo "- Backend processed: $TEST_DIR/processed_backend.mp4"
echo "- Locally processed: $TEST_DIR/processed_local.mp4"
