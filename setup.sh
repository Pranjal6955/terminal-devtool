#!/bin/bash
# Quick start script for Terminal DevTool

echo "🎛️ Terminal DevTool Quick Start"
echo "=============================="
echo

# Check if Go is installed with proper version
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install it first."
    echo "Visit https://golang.org/doc/install for instructions."
    exit 1
else
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo "✅ Go version $GO_VERSION detected"
fi

# Check if Rust is installed
if ! command -v rustc &> /dev/null; then
    echo "❌ Rust is not installed. Please install it first."
    echo "Visit https://www.rust-lang.org/tools/install for instructions."
    exit 1
else
    RUST_VERSION=$(rustc --version | awk '{print $2}')
    echo "✅ Rust version $RUST_VERSION detected"
fi

# Check if FFmpeg is installed
if ! command -v ffmpeg &> /dev/null; then
    echo "❌ FFmpeg is not installed. Please install it first."
    echo "Visit https://ffmpeg.org/download.html for instructions."
    exit 1
else
    FFMPEG_VERSION=$(ffmpeg -version | head -n1 | awk '{print $3}')
    echo "✅ FFmpeg version $FFMPEG_VERSION detected"
fi

# Check for Docker (optional)
if command -v docker &> /dev/null && command -v docker-compose &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    COMPOSE_VERSION=$(docker-compose --version | awk '{print $3}' | sed 's/,//')
    echo "✅ Docker version $DOCKER_VERSION and Docker Compose version $COMPOSE_VERSION detected"
else
    echo "⚠️ Docker not detected - Docker-based deployment won't be available"
    echo "   Install Docker from https://docs.docker.com/get-docker/ if needed"
fi
echo

# Create necessary directories
echo "📁 Creating project directories..."
mkdir -p media test_media
echo

echo "🔧 Setting up Terminal DevTool..."

# Build backend
echo "🔄 Building Go backend..."
cd backend
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ Go dependencies installation failed"
    exit 1
fi

go build -o terminal-devtool-backend
if [ $? -ne 0 ]; then
    echo "❌ Backend build failed"
    exit 1
fi
echo "✅ Backend built successfully"
echo

# Build frontend
echo "🔄 Building Rust frontend..."
cd ../frontend-cli
cargo build
if [ $? -ne 0 ]; then
    echo "❌ Frontend build failed"
    exit 1
fi
echo "✅ Frontend built successfully"
echo

# Back to root directory
cd ..

echo "🚀 Terminal DevTool setup complete!"
echo
echo "To run the backend server:"
echo "  cd backend && ./terminal-devtool-backend"
echo
echo "To use the frontend CLI:"
echo "  cd frontend-cli && cargo run -- --help"
echo
echo "Try some example commands:"
echo "  # Process a video file"
echo "  cd frontend-cli && cargo run -- process -i \"path/to/video.mp4\" -o \"output.mp4\" --resolution 1280x720"
echo
echo "  # Compare original and processed files"
echo "  cd frontend-cli && cargo run -- compare -o \"path/to/original.mp4\" -p \"path/to/processed.mp4\""
echo
echo "  # Get media file info"
echo "  cd frontend-cli && cargo run -- info -f \"path/to/video.mp4\""
echo
echo "Or run the test examples (generates a test video):"
echo "  ./test_examples.sh"
echo
echo "For more information, see the README.md file."
