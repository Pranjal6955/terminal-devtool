use clap::Parser;

/// Terminal DevTool CLI
#[derive(Parser)]
#[command(name = "mediatool")]
#[command(version = "1.0")]
#[command(about = "Convert, compress, and generate media files", long_about = None)]
struct Cli {
    #[arg(short, long)]
    input: String,

    #[arg(short, long)]
    output: Option<String>,

    #[arg(long)]
    resolution: Option<String>,

    #[arg(long)]
    bitrate: Option<String>,

    #[arg(long)]
    format: Option<String>,
}

fn main() {
    let args = Cli::parse();
    println!("Processing: {}", args.input);
    // Integrate ffmpeg command execution here
}
