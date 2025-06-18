mod cli;
mod processor;

use clap::Parser;
use cli::CliArgs;

fn main() {
    let args = CliArgs::parse();
    println!("✅ CLI Args: {:#?}", args);

    // Call the processor
    processor::process_media(args);
}
