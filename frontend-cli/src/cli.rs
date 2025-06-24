use clap::{Parser, Subcommand};

/// Terminal DevTool CLI
#[derive(Parser, Debug)]
#[command(author, version, about = "Modern FFmpeg CLI Tool", long_about = None)]
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
    /// Process media file
    Process {
        /// Input file
        #[arg(short, long)]
        input: String,

        /// Output file
        #[arg(short, long)]
        output: Option<String>,

        /// Desired resolution (e.g., 1280x720)
        #[arg(long)]
        resolution: Option<String>,

        /// Bitrate for compression (e.g., 1000k)
        #[arg(long)]
        bitrate: Option<String>,

        /// Output format (e.g., webm, mp4)
        #[arg(long)]
        format: Option<String>,

        /// Print the command that would be run without executing it
        #[arg(long)]
        dry_run: bool,
    },
    
    /// Compare two media files
    Compare {
        /// Original file
        #[arg(short, long)]
        original: String,
        
        /// Processed file
        #[arg(short, long)]
        processed: String,
    },
    
    /// Get media file info
    Info {
        /// Media file path
        #[arg(short, long)]
        file: String,
    },

    /// Convert media file format
    Convert {
        /// Input file
        #[arg(long)]
        input: String,

        /// Output file (optional)
        #[arg(long)]
        output: Option<String>,

        /// Target format (e.g., webm, jpg)
        #[arg(long)]
        format: String,
    },

    /// Compress a video file
    Compress {
        /// Input file
        #[arg(long)]
        input: String,

        /// Output file (optional)
        #[arg(long)]
        output: Option<String>,

        /// Target bitrate (e.g., 1000k, 2M)
        #[arg(long)]
        bitrate: String,
    },
}
