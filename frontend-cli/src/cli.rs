use clap::{Parser, Subcommand};

/// Terminal DevTool CLI
#[derive(Parser, Debug)]
#[command(
    author,
    version,
    about = "🎬 Project A — High-Performance Media Toolkit",
    long_about = "\
  ____            _           _     _    \n |  _ \\ ___  _ __| |__   ___ | |__ (_)___\n | |_) / _ \\| '__| '_ \\ / _ \\| '_ \\| / __|\n |  __/ (_) | |  | |_) | (_) | |_) | \\__ \\\n |_|   \\___/|_|  |_.__/ \\___/|_.__/|_|___/\n\nA modern, beautiful, and high-performance media CLI toolkit.\n\n✨ Fast, friendly, and feature-rich! ✨\n",
    next_line_help = true
)]
pub struct CliArgs {
    /// Backend server URL (default: http://localhost:8080)
    #[arg(long, env = "DEVTOOL_BACKEND_URL")]
    pub backend_url: Option<String>,

    /// Run locally without backend
    #[arg(long)]
    pub local: bool,

    #[command(subcommand)]
    pub command: Commands,
}

#[derive(Subcommand, Debug)]
pub enum Commands {
    /// 🎬 Process media file
    #[command(
        about = "Process a media file (transcode, resize, etc.)",
        long_about = "Process a media file with options for transcoding, resizing, bitrate, and format.\n\nExample:\n  terminal-devtool process --input input.mp4 --output output.webm --bitrate 1000k --format webm --dry-run\n"
    )]
    Process {
        /// 📥 Input file
        #[arg(short, long)]
        input: String,
        /// 📤 Output file
        #[arg(short, long)]
        output: Option<String>,
        /// 🖼️ Desired resolution (e.g., 1280x720)
        #[arg(long)]
        resolution: Option<String>,
        /// 🎚️ Bitrate for compression (e.g., 1000k)
        #[arg(long)]
        bitrate: Option<String>,
        /// 🗂️ Output format (e.g., webm, mp4)
        #[arg(long)]
        format: Option<String>,
        /// 🧪 Print the command that would be run without executing it
        #[arg(long)]
        dry_run: bool,
    },
    
    /// 🆚 Compare two media files
    #[command(
        about = "Compare two media files (size, quality, etc.)",
        long_about = "Compare two media files and display differences in size, format, resolution, and more.\n\nExample:\n  terminal-devtool compare --original original.mp4 --processed processed.webm\n"
    )]
    Compare {
        /// 🗂️ Original file
        #[arg(short, long)]
        original: String,
        /// 🗂️ Processed file
        #[arg(short, long)]
        processed: String,
    },
    
    /// ℹ️ Get media file info
    #[command(
        about = "Get info about a media file",
        long_about = "Display detailed information about a media file, including format, duration, resolution, bitrate, and size.\n\nExample:\n  terminal-devtool info --file input.mp4\n"
    )]
    Info {
        /// 🗂️ Media file path
        #[arg(short, long)]
        file: String,
    },

    /// 🔄 Convert media file format
    #[command(
        about = "Convert a media file to another format",
        long_about = "Convert a media file to a different format (e.g., mp4 to webm, jpg to png).\n\nExample:\n  terminal-devtool convert --input input.mp4 --format webm\n"
    )]
    Convert {
        /// 📥 Input file
        #[arg(long)]
        input: String,

        /// 📤 Output file (optional)
        #[arg(long)]
        output: Option<String>,

        /// 🗂️ Target format (e.g., webm, jpg)
        #[arg(long)]
        format: String,
    },

    /// 📦 Compress a video file
    #[command(
        about = "Compress a video file to a target bitrate",
        long_about = "Compress a video file to a specified bitrate.\n\nExample:\n  terminal-devtool compress --input input.mp4 --bitrate 500k\n"
    )]
    Compress {
        /// 📥 Input file
        #[arg(long)]
        input: String,

        /// 📤 Output file (optional)
        #[arg(long)]
        output: Option<String>,

        /// 🎚️ Target bitrate (e.g., 1000k, 2M)
        #[arg(long)]
        bitrate: String,
    },
}
