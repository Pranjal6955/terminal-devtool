use std::{process::Command, path::Path};
use anyhow::{Result, Context};
use colored::*;
use prettytable::{Table, row};

use crate::{cli::{CliArgs, Commands}, client::ApiClient, models::{ProcessRequest, MediaInfo, CompareResult, CompressOptions}};

pub fn process_media(args: CliArgs) -> Result<()> {
    // Create API client if not in local mode
    let client = if !args.local {
        Some(ApiClient::new(args.backend_url))
    } else {
        None
    };
    
    // Check backend health if using API
    if let Some(client) = &client {
        if let Err(e) = client.check_health() {
            eprintln!("{} {}", "⚠️".yellow(), format!("Backend connection failed: {}", e).yellow());
            eprintln!("{} {}", "ℹ️".blue(), "Falling back to local mode.".blue());
        }
    }
    
    match args.command {
        Commands::Process { 
            input,
            output,
            resolution,
            bitrate,
            format,
            dry_run: _
        } => {
            // TODO: Integrate backend API for processing media here
            // Example: client.process_media(...)
            // Remove local ffmpeg and dry run logic
        },
        Commands::Compare { 
            original, 
            processed 
        } => {
            // TODO: Integrate backend API for comparing media here
            // Example: client.compare_media(...)
        },
        Commands::Info { 
            file 
        } => {
            // TODO: Integrate backend API for getting media info here
            // Example: client.get_media_info(...)
        },
        Commands::Convert {
            input,
            output,
            format
        } => {
            // TODO: Integrate backend API for converting media format here
            // Example: client.convert_media(...)
        },
        Commands::Compress {
            input,
            output,
            bitrate
        } => {
            // TODO: Integrate backend API for compressing media here
            // Example: client.compress_media(...)
        }
    }
    
    Ok(())
}

fn compress_command(opts: CompressOptions) -> Result<()> {
    if !opts.bitrate.ends_with('k') && !opts.bitrate.ends_with('M') {
        eprintln!("{} {}", "[x]".bright_red(), "Invalid bitrate format. Must end with 'k' or 'M'.".bright_red());
        return Err(anyhow::anyhow!("Invalid bitrate format. Must end with 'k' or 'M'."));
    }

    if let Some(output_path) = opts.output {
        println!("{} {}", "[✓]".bright_green(), format!("Compressing {} to {} with bitrate {}...", opts.input, output_path, opts.bitrate).bright_green());
    } else {
        println!("{} {}", "[✓]".bright_green(), format!("Compressing {} with bitrate {}...", opts.input, opts.bitrate).bright_green());
    }

    Ok(())
}

fn process_command(
    client: Option<ApiClient>,
    input: String,
    output: Option<String>,
    resolution: Option<String>,
    bitrate: Option<String>,
    format: Option<String>,
    dry_run: bool
) -> Result<()> {
    if dry_run {
        println!("{} {}", "[!]".yellow(), "Dry run mode - displaying command without executing...".yellow());
    } else {
        println!("{} {}", "[✓]".bright_green(), "Starting media processing...".bright_green());
    }
    
    // Use API client if available
    if let Some(client) = client {
        let request = ProcessRequest {
            input: input.clone(),
            output: output.clone(),
            resolution: resolution.clone(),
            bitrate: bitrate.clone(),
            format: format.clone(),
            dry_run,
        };
        
        println!("{} {}", "🌐".cyan(), "Sending request to backend server...".cyan());
        
        match client.process_media(request) {
            Ok(response) => {
                println!("{} {}", "✅".green(), format!("Processing complete. Output: {}", response.output).green());
                return Ok(());
            },
            Err(e) => {
                eprintln!("{} {}", "⚠️".yellow(), format!("Backend processing failed: {}", e).yellow());
                eprintln!("{} {}", "ℹ️".blue(), "Falling back to local processing...".blue());
            }
        }
    }
    
    // Fall back to local processing
    let output = output.unwrap_or_else(|| {
        let output_name = format!("processed_{}", Path::new(&input).file_name().unwrap().to_string_lossy());
        if let Some(fmt) = &format {
            return format!("{}.{}", output_name, fmt);
        }
        output_name.to_string()
    });
    
    let mut ffmpeg_args = vec!["-i", &input];
    
    if let Some(res) = &resolution {
        ffmpeg_args.push("-s");
        ffmpeg_args.push(res);
    }
    
    if let Some(br) = &bitrate {
        ffmpeg_args.push("-b:v");
        ffmpeg_args.push(br);
    }
    
    ffmpeg_args.push(&output);
    
    // Build command string
    let cmd_string = format!("ffmpeg {}", ffmpeg_args.join(" "));
    
    if dry_run {
        // Just print the command that would be run
        println!("{} {}", "[!]".yellow(), format!("[Dry Run] {}", cmd_string).yellow());
        return Ok(());
    }
    
    println!("{} {}", "🛠️".yellow(), format!("Running ffmpeg locally with args: {:?}", ffmpeg_args).yellow());
    
    let status = Command::new("ffmpeg")
        .args(&ffmpeg_args)
        .status()
        .context("Failed to execute ffmpeg")?;
        
    if status.success() {
        println!("{} {}", "[✓]".bright_green(), format!("Processing complete. Output: {}", output).bright_green());
    } else {
        eprintln!("{} {}", "[x]".bright_red(), "Processing failed".bright_red());
        return Err(anyhow::anyhow!("ffmpeg command failed"));
    }
    
    Ok(())
}

fn compare_command(client: Option<ApiClient>, original: &str, processed: &str) -> Result<()> {
    println!("{} {}", "[!]".bright_blue(), "Comparing media files...".bright_blue());
    let result = if let Some(client) = client {
        println!("{} {}", "[!]".yellow(), "Sending comparison request to backend server...".yellow());
        client.compare_media(original, processed)?
    } else {
        eprintln!("{} {}", "[x]".bright_red(), "Local comparison requires backend server".bright_red());
        return Err(anyhow::anyhow!("Local comparison requires backend server"));
    };
    print_comparison_result(&result);
    Ok(())
}

fn info_command(client: Option<ApiClient>, file_path: &str) -> Result<()> {
    println!("{} {}", "[!]".bright_blue(), format!("Getting info for {}...", file_path).bright_blue());
    let info = if let Some(client) = client {
        println!("{} {}", "[!]".yellow(), "Sending info request to backend server...".yellow());
        client.get_media_info(file_path)?
    } else {
        eprintln!("{} {}", "[x]".bright_red(), "Local media info requires backend server".bright_red());
        return Err(anyhow::anyhow!("Local media info requires backend server"));
    };
    print_media_info(&info);
    Ok(())
}

fn print_media_info(info: &MediaInfo) {
    let mut table = Table::new();
    
    table.add_row(row!["Property", "Value"]);
    table.add_row(row!["Filename", &info.filename]);
    table.add_row(row!["Format", &info.format]);
    table.add_row(row!["Duration", &info.duration]);
    table.add_row(row!["Resolution", &info.resolution]);
    table.add_row(row!["Bitrate", &info.bitrate]);
    table.add_row(row!["Size (bytes)", &info.size.to_string()]);
    
    println!();
    table.printstd();
    println!();
}

fn print_comparison_result(result: &CompareResult) {
    let mut table = Table::new();
    table.add_row(row!["Property", "Original", "Processed"]);
    table.add_row(row!["Filename", &result.original.filename, &result.processed.filename]);
    table.add_row(row!["Format", &result.original.format, &result.processed.format]);
    table.add_row(row!["Resolution", &result.original.resolution, &result.processed.resolution]);
    table.add_row(row!["Bitrate", &result.original.bitrate, &result.processed.bitrate]);
    table.add_row(row!["Size (bytes)", 
                       &result.original.size.to_string(), 
                       &result.processed.size.to_string()]);
    println!();
    table.printstd();
    println!();
    // Print size difference
    let diff_text = format!("Size reduction: {:.2}%", result.size_diff_percent.abs());
    if result.size_diff_percent > 0.0 {
        println!("{} {}", "[✓]".bright_green(), diff_text.bright_green());
    } else if result.size_diff_percent < 0.0 {
        println!("{} {}", "[x]".bright_red(), diff_text.bright_red());
    } else {
        println!("{} {}", "[!]".yellow(), diff_text.yellow());
    }
    println!();
}

