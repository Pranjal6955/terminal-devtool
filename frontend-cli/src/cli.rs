use clap::Parser;

/// Terminal DevTool CLI
#[derive(Parser, Debug)]
#[command(author, version, about = "Modern FFmpeg CLI Tool", long_about = None)]
pub struct CliArgs {
    /// Input file
    #[arg(short, long)]
    pub input: String,

    /// Output file
    #[arg(short, long)]
    pub output: Option<String>,

    /// Desired resolution (e.g., 1280x720)
    #[arg(long)]
    pub resolution: Option<String>,

    /// Bitrate for compression (e.g., 1000k)
    #[arg(long)]
    pub bitrate: Option<String>,

    /// Output format (e.g., webm, mp4)
    #[arg(long)]
    pub format: Option<String>,
}
