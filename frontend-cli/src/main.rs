mod cli;
mod processor;
mod client;
mod models;

use std::process;
use clap::Parser;
use cli::CliArgs;
use colored::*;

fn main() {
    // Print beautiful welcome header
    println!("\n🎬\x1b[1;34m Project A — High-Performance Media Toolkit\x1b[0m");
    println!("\x1b[1;34m  ____            _           _     _    \n |  _ \\ ___  _ __| |__   ___ | |__ (_)___\n | |_) / _ \\| '__| '_ \\ / _ \\| '_ \\| / __|\n |  __/ (_) | |  | |_) | (_) | |_) | \\__ \\\n |_|   \\___/|_|  |_.__/ \\___/|_.__/|_|___/\x1b[0m");
    println!("✨ Fast, friendly, and feature-rich! ✨\n");

    let args = CliArgs::parse();

    // For Convert subcommand, just print dry-run message (no error handling needed)
    match &args.command {
        cli::Commands::Convert { .. } => {
            let _ = processor::process_media(args);
        },
        _ => {
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
    }
}
