mod cli;
mod processor;
mod client;
mod models;

use std::process;
use clap::Parser;
use cli::CliArgs;
use colored::*;

fn main() {
    let args = CliArgs::parse();

    // Call the processor and handle any errors
    if let Err(e) = processor::process_media(args) {
        eprintln!("{} {}", "❌".bright_red(), format!("Error: {}", e).bright_red());
        
        // Print cause chain
        let mut source = e.source();
        while let Some(cause) = source {
            eprintln!("   {}", format!("Caused by: {}", cause).red());
            source = cause.source();
        }
        
        process::exit(1);
    }
}
