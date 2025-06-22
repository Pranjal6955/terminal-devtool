# 🎛️ Terminal DevTool

A high-performance terminal-based media processing tool that lets you convert, compress, compare, and generate media files (images/videos) using custom flags — just like `FFmpeg`, but with a modern CLI interface and clean architecture.

## 📝 Summary

Terminal DevTool combines the power of Go and Rust to create a modern media processing utility that leverages FFmpeg. The project features:

- **Client-Server Architecture**: A Go backend handles heavy media processing, while a Rust CLI provides the user interface
- **Dual Mode Operation**: Work online with the backend or offline with local processing
- **Media Comparison**: Compare media files to analyze optimization results
- **Modern CLI Experience**: Colored output, pretty-printed tables, and intuitive commands
- **Easy to Extend**: Architecture designed for adding new features and media processing capabilities

---

## ⚙️ Features

- ✅ Supports processing of images and videos
- 🎞️ Convert between formats (e.g., `.mp4` to `.webm`)
- 📏 Resize media with resolution flags
- 💨 Compress videos using bitrate control
- 📊 Compare media files to analyze optimization results
- 🛠️ Built with performance-first technologies (Rust + Go)
- 🧩 Modern client-server architecture

---

## 🧱 Architecture

```
terminal-devtool/
├── backend/         # Go-based backend for heavy processing logic
│   ├── api/         # HTTP API handlers
│   └── media/       # Media processing utilities
└── frontend-cli/    # Rust-based CLI built using Clap
    └── src/         # Source code with client-server architecture
```

- **Backend:** Written in Go for easy integration with `ffmpeg` and scalable backend processing.
- **Frontend CLI:** Built in Rust with `clap` for a modern, interactive terminal experience.
- **Client-Server Model:** The CLI can operate locally or communicate with the backend for advanced features.

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.18+** - For the backend server
  - Check with: `go version`
  - Install from: [https://golang.org/doc/install](https://golang.org/doc/install)
- **Rust 1.68+** - For the frontend CLI
  - Check with: `rustc --version`
  - Install from: [https://www.rust-lang.org/tools/install](https://www.rust-lang.org/tools/install)
- **FFmpeg** - For media processing
  - Check with: `ffmpeg -version`
  - Install from: [https://ffmpeg.org/download.html](https://ffmpeg.org/download.html)
- **Docker & Docker Compose** (Optional) - For containerized deployment
  - Check with: `docker --version` and `docker-compose --version`
  - Install from: [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)

### 1. Clone the Repository and Quick Setup

```bash
# Clone the repository
git clone https://github.com/Promptzy/terminal-devtool.git
cd terminal-devtool

# Run the setup script
chmod +x setup.sh  # Make the script executable if needed
./setup.sh
```

The setup script will:
- Check for required dependencies
- Build the backend Go application
- Build the frontend Rust CLI
- Create necessary directories
- Provide instructions to get started

### 2. Running the Backend (Go)

The backend provides HTTP endpoints for media processing, comparison, and information retrieval.

#### Option A: Run directly

```bash
# Navigate to the backend directory
cd backend

# Build and run (if not done by setup.sh)
go build -o terminal-devtool-backend
./terminal-devtool-backend

# Or run with go run
go run main.go
```

#### Option B: Run with Docker

```bash
# Build and start the backend server in a container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container when done
docker-compose down
```

The backend server will be available at `http://localhost:8080`. You can test it with:

```bash
# Check if the server is running
curl http://localhost:8080/health
```

### 3. Running the Frontend CLI (Rust)

The frontend CLI provides a user-friendly interface to interact with the backend or process media locally.

```bash
# Navigate to the frontend directory
cd frontend-cli

# Build the CLI (if not done by setup.sh)
cargo build

# Display help information
cargo run -- --help

# Or use the compiled binary directly
./target/debug/frontend-cli --help
```

## 🧪 Example Usage

### Process Media Files

```bash
# Process with backend (default)
cargo run -- process -i input.mp4 -o output.mp4 --resolution 1280x720 --bitrate 1000k --format mp4

# Process locally without backend
cargo run -- --local process -i input.mp4 -o output.mp4 --resolution 1280x720
```

### Compare Media Files

```bash
# Compare original and processed files
cargo run -- compare -o original.mp4 -p processed.mp4
```

### Get Media File Info

```bash
# Get detailed information about a media file
cargo run -- info -f video.mp4
```

## 🛠️ Command Reference

### Global Options

| Option | Description | Example |
|------|-------------|---------|
| `--backend-url` | Backend server URL | `--backend-url http://localhost:9000` |
| `--local` | Run locally without backend | `--local` |

### Process Command

| Option | Description | Example |
|------|-------------|---------|
| `-i, --input` | Input file path | `-i input.mp4` |
| `-o, --output` | Output file path (optional) | `-o output.mp4` |
| `--resolution` | Resize the media | `--resolution 1920x1080` |
| `--bitrate` | Set output bitrate | `--bitrate 800k` |
| `--format` | Convert to specific format | `--format webm` |

### Compare Command

| Option | Description | Example |
|------|-------------|---------|
| `-o, --original` | Original file path | `-o original.mp4` |
| `-p, --processed` | Processed file path | `-p processed.mp4` |

### Info Command

| Option | Description | Example |
|------|-------------|---------|
| `-f, --file` | Media file path | `-f video.mp4` |

## 🧰 Built With

- Rust (CLI frontend)
  - Clap for command-line parsing
  - Reqwest for HTTP communication
  - Serde for JSON serialization
  - Colored for terminal output
  - PrettyTable for formatted data display
- Go (Backend)
  - Standard library HTTP server
  - FFmpeg integration
- FFmpeg (Core media processing)

## 🔄 How It Works

1. The frontend CLI parses commands and options
2. If using the backend:
   - The CLI sends HTTP requests to the backend server
   - The backend processes the media using FFmpeg
   - Results are sent back to the CLI for display
3. If in local mode:
   - The CLI directly executes FFmpeg commands
   - Basic functionality works but advanced features like comparison require the backend

## 🚀 Advanced Features

### Media Comparison

The comparison feature analyzes two media files and shows:
- Format differences
- Resolution changes
- Bitrate differences
- Size reduction percentage

### Terminal Output

The CLI uses colored output and formatted tables to display:
- Processing status and progress
- Media file information
- Comparison results

## 🧪 Testing the Application

To quickly test all features of Terminal DevTool, run the included test script:

```bash
# Make the script executable (if needed)
chmod +x test_examples.sh

# Run the test script
./test_examples.sh
```

The test script will:
1. Generate a test video using FFmpeg
2. Start the backend server
3. Run various commands to demonstrate the CLI capabilities
4. Process media files with different settings
5. Compare the results
6. Display media information

## ❓ Troubleshooting

### Backend Issues

1. **Backend won't start**:
   - Check if port 8080 is already in use: `lsof -i :8080`
   - Try changing the port: `PORT=8081 ./terminal-devtool-backend`

2. **Processing fails**:
   - Verify FFmpeg is installed and available in PATH
   - Check file permissions of input/output directories
   - Look for detailed errors in the terminal output

### Frontend Issues

1. **CLI compilation errors**:
   - Update Rust: `rustup update`
   - Clean and rebuild: `cargo clean && cargo build`

2. **Connection to backend fails**:
   - Verify the backend is running: `curl http://localhost:8080/health`
   - Check if the `--backend-url` parameter is correct
   - Try using `--local` mode to bypass the backend

### Docker Issues

1. **Docker container won't start**:
   - Check Docker logs: `docker-compose logs`
   - Verify Docker has sufficient permissions
   - Ensure port 8080 is not already in use

2. **Volume mounting issues**:
   - Check paths in docker-compose.yml
   - Ensure the media directory exists and has proper permissions

## 📌 Future Enhancements

- [ ] Add batch processing support
- [ ] Implement WebSocket for real-time progress updates
- [ ] Add support for image optimization
- [ ] Create pre-compiled binaries for multiple platforms
- [ ] Add GPU acceleration options
- [ ] Create config files for preset configurations

## 🧑‍💻 Development

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend-cli
cargo test
```

### Building Release Binaries

```bash
# Frontend CLI
cd frontend-cli
cargo build --release
```

The compiled binary will be available at `frontend-cli/target/release/frontend-cli`. 
5. Add tests for core functionality
6. Implement file format conversions
7. Add compression and resizing features

## 📄 License

This project is licensed under the MIT License.

## 🤝 Contributing

Feel free to fork and submit a PR. Contributions are welcome!
