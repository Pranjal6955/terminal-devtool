# 🎛️ Terminal DevTool

A high-performance terminal-based media processing tool that lets you convert, compress, and generate media files (images/videos) using custom flags — just like `FFmpeg`, but with a modern CLI interface and clean architecture.

---

## ⚙️ Features

- ✅ Supports batch processing of images and videos
- 🎞️ Convert between formats (e.g., `.mp4` to `.webm`)
- 📏 Resize media with resolution flags
- 💨 Compress videos using bitrate control
- 🛠️ Built with performance-first technologies (Rust + Go)
- 🧩 Easily extendable architecture

---

## 🧱 Architecture

```
terminal-devtool/
├── backend/ # Go-based backend for heavy processing logic
├── frontend-cli/ # Rust-based CLI built using Clap
└── README.md
```

- **Backend:** Written in Go for easy integration with `ffmpeg` and scalable backend processing.
- **Frontend CLI:** Built in Rust with `clap` for a modern, interactive terminal experience.

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/Promptzy/terminal-devtool.git
cd terminal-devtool
```

### 2. Backend (Go)
📦 Install Dependencies
Make sure Go is installed (`go version`)

```bash
cd backend
go mod tidy
go run main.go
```

### 3. Frontend CLI (Rust)
📦 Install Dependencies
Make sure Rust is installed (`rustc --version`)

```bash
cd frontend-cli
cargo run -- --help
```

## 🧪 Example Usage

```bash
cargo run -- -i input.mp4 -o output.webm --resolution 1280x720 --bitrate 1000k --format webm
```

## 🛠️ Flags Reference

| Flag | Description | Example |
|------|-------------|---------|
| `-i`nput | Input file path | `-i input.mp4` |
| `-o`utput | Output file path (optional) | `-o output.mp4` |
| `--resolution` | Resize the media | `--resolution 1920x1080` |
| `--bitrate` | Set output bitrate | `--bitrate 800k` |
| `--format` | Convert to specific format | `--format webm` |

## 🧰 Built With

- Rust
- Go
- FFmpeg
- Clap (Rust CLI Framework)

## 📌 To-Do

- [ ] Add ffmpeg execution in CLI
- [ ] Connect frontend to backend via IPC/gRPC (optional)
- [ ] Add support for image conversion
- [ ] Create release binaries

## 🗺️ Implementation Roadmap

### 1. Backend (Go)

Create the following files:
- `/backend/main.go` - Entry point for the Go backend
- `/backend/processor/video.go` - Video processing logic using FFmpeg
- `/backend/processor/image.go` - Image processing logic
- `/backend/api/server.go` - API endpoints for CLI communication
- `/backend/go.mod` - Go module definition

### 2. Frontend CLI (Rust)

Create the following files:
- `/frontend-cli/Cargo.toml` - Rust project configuration
- `/frontend-cli/src/main.rs` - CLI entry point
- `/frontend-cli/src/cli.rs` - Command-line argument parsing with Clap
- `/frontend-cli/src/processor.rs` - Media processing interface
- `/frontend-cli/src/utils.rs` - Helper functions

### 3. Configuration

- `/config.toml` - Application configuration (paths, defaults)

### 4. Development Steps

1. Set up Go module and basic FFmpeg wrapper
2. Implement Rust CLI using Clap
3. Create basic media processing functions
4. Connect CLI to backend processing 
5. Add tests for core functionality
6. Implement file format conversions
7. Add compression and resizing features

## 📄 License

This project is licensed under the MIT License.

## 🤝 Contributing

Feel free to fork and submit a PR. Contributions are welcome!
