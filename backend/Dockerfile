# Backend Dockerfile
FROM golang:1.20-alpine AS builder

# Install FFmpeg, git and build dependencies
RUN apk add --no-cache ffmpeg git curl ca-certificates build-base

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
# Copy go.sum if it exists
COPY go.sum* ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o terminal-devtool-backend .

# Create a minimal production image
FROM alpine:latest

# Install FFmpeg and runtime dependencies
RUN apk add --no-cache ffmpeg ca-certificates curl

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/terminal-devtool-backend .

# Create directories for media files
RUN mkdir -p /app/media /app/test_media && \
    chmod +x /app/terminal-devtool-backend

# Expose the server port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=5s \
  CMD curl -f http://localhost:8080/health || exit 1

# Command to run the application
CMD ["./terminal-devtool-backend"]
