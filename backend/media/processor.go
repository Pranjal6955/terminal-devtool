package media

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// MediaInfo represents metadata about a media file
type MediaInfo struct {
	Filename   string `json:"filename"`
	Format     string `json:"format"`
	Duration   string `json:"duration"`
	Resolution string `json:"resolution"`
	Bitrate    string `json:"bitrate"`
	Size       int64  `json:"size"`
	Codec      string `json:"codec"`
	FrameRate  string `json:"frame_rate"`
}

// ProcessRequest represents a request to process media
type ProcessRequest struct {
	Input      string `json:"input"`
	Output     string `json:"output,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	Bitrate    string `json:"bitrate,omitempty"`
	Format     string `json:"format,omitempty"`
	Codec      string `json:"codec,omitempty"`
	FrameRate  string `json:"frame_rate,omitempty"`
	CRF        string `json:"crf,omitempty"`    // Constant Rate Factor for quality-based compression
	Preset     string `json:"preset,omitempty"` // Encoding preset (ultrafast, fast, medium, slow, etc.)
	DryRun     bool   `json:"dry_run,omitempty"` // If true, return command string without executing
}

// ProcessProgress represents the progress of a media processing operation
type ProcessProgress struct {
	Stage    string  `json:"stage"`
	Progress float64 `json:"progress"`
	ETA      string  `json:"eta"`
	Speed    string  `json:"speed"`
}

// CompareResult represents the result of a media comparison
type CompareResult struct {
	Original          MediaInfo `json:"original"`
	Processed         MediaInfo `json:"processed"`
	SizeDiff          float64   `json:"size_diff_percent"`
	ResolutionChanged bool      `json:"resolution_changed"`
	BitrateReduction  float64   `json:"bitrate_reduction_percent"`
	FormatChanged     bool      `json:"format_changed"`
	CodecChanged      bool      `json:"codec_changed"`
}

// GetMediaInfo retrieves information about a media file using ffprobe
func GetMediaInfo(filepath string) (MediaInfo, error) {
	info := MediaInfo{
		Filename: filepath,
	}

	// Get file size
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return info, fmt.Errorf("failed to get file info: %w", err)
	}
	info.Size = fileInfo.Size()

	// Run ffprobe to get detailed media information in JSON format
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filepath)

	output, err := cmd.Output()
	if err != nil {
		return info, fmt.Errorf("ffprobe failed: %w", err)
	}

	// Parse JSON output
	var ffprobeOutput struct {
		Streams []struct {
			CodecName  string `json:"codec_name"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			BitRate    string `json:"bit_rate"`
			RFrameRate string `json:"r_frame_rate"`
		} `json:"streams"`
		Format struct {
			Filename   string `json:"filename"`
			FormatName string `json:"format_name"`
			Duration   string `json:"duration"`
			BitRate    string `json:"bit_rate"`
		} `json:"format"`
	}

	if err := json.Unmarshal(output, &ffprobeOutput); err != nil {
		return info, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	// Extract video stream information
	for _, stream := range ffprobeOutput.Streams {
		// Only look at video streams (skip audio)
		if stream.Width > 0 && stream.Height > 0 {
			info.Resolution = fmt.Sprintf("%dx%d", stream.Width, stream.Height)
			if stream.BitRate != "" {
				info.Bitrate = stream.BitRate
			}
			info.Codec = stream.CodecName

			// Parse framerate (often in the format "24/1")
			if stream.RFrameRate != "" {
				parts := strings.Split(stream.RFrameRate, "/")
				if len(parts) == 2 {
					numerator, denominator := parts[0], parts[1]
					if denominator == "1" {
						info.FrameRate = numerator + " fps"
					} else {
						// Calculate the framerate if denominator is not 1
						num, err1 := strconv.ParseFloat(numerator, 64)
						den, err2 := strconv.ParseFloat(denominator, 64)
						if err1 == nil && err2 == nil && den > 0 {
							fps := num / den
							info.FrameRate = fmt.Sprintf("%.2f fps", fps)
						} else {
							info.FrameRate = stream.RFrameRate
						}
					}
				} else {
					info.FrameRate = stream.RFrameRate
				}
			}
			break
		}
	}

	// Extract format information
	if ffprobeOutput.Format.FormatName != "" {
		info.Format = ffprobeOutput.Format.FormatName
	}

	if ffprobeOutput.Format.Duration != "" {
		info.Duration = ffprobeOutput.Format.Duration + "s"
	}

	// Use format bitrate as fallback if stream bitrate is not available
	if info.Bitrate == "" && ffprobeOutput.Format.BitRate != "" {
		info.Bitrate = ffprobeOutput.Format.BitRate
	}

	return info, nil
}

// ProcessMedia processes a media file based on the request parameters
func ProcessMedia(req ProcessRequest) (string, error) {
	// Set default output if not provided
	output := req.Output
	if output == "" {
		ext := ".mp4"
		if req.Format != "" {
			ext = "." + req.Format
		}
		output = "processed_" + filepath.Base(req.Input)
		output = strings.TrimSuffix(output, filepath.Ext(output)) + ext
	}

	// Skip getting media info if it's a dry run
	var duration time.Duration
	if !req.DryRun {
		// First, get the input file duration
		inputInfo, err := GetMediaInfo(req.Input)
		if err != nil {
			return "", fmt.Errorf("failed to get input file info: %w", err)
		}

		// Parse duration string to time.Duration
		if inputInfo.Duration != "" {
			durationStr := strings.TrimSuffix(inputInfo.Duration, "s")
			if durationFloat, err := strconv.ParseFloat(durationStr, 64); err == nil {
				duration = time.Duration(durationFloat * float64(time.Second))
			}
		}
	}

	// Build ffmpeg command with global options
	args := []string{
		"-hide_banner",  // Hide FFmpeg banner info
		"-y",            // Overwrite output files without asking
		"-i", req.Input, // Input file
		"-progress", "pipe:1", // Output progress to stdout
	}

	// Add video-specific options
	if req.Resolution != "" {
		args = append(args, "-s", req.Resolution)
	}

	if req.Bitrate != "" {
		args = append(args, "-b:v", req.Bitrate)
	}

	// Add codec if specified
	if req.Codec != "" {
		args = append(args, "-c:v", req.Codec)
	} else if req.Format == "webm" {
		// Use VP9 for WebM if codec not specified
		args = append(args, "-c:v", "libvpx-vp9")
	} else if req.Format == "mp4" {
		// Use H.264 for MP4 if codec not specified
		args = append(args, "-c:v", "libx264")
	}

	// Add framerate if specified
	if req.FrameRate != "" {
		args = append(args, "-r", req.FrameRate)
	}

	// Add constant rate factor (quality) if specified
	if req.CRF != "" {
		args = append(args, "-crf", req.CRF)
	}

	// Add encoding preset if specified
	if req.Preset != "" {
		args = append(args, "-preset", req.Preset)
	}

	// Handle audio based on format
	if req.Format == "gif" {
		// Remove audio for GIF
		args = append(args, "-an")
	} else {
		// Maintain audio quality or use AAC for most formats
		if req.Format == "webm" {
			args = append(args, "-c:a", "libopus")
		} else {
			args = append(args, "-c:a", "aac")
		}
	}

	// Add output format if specified and not in filename
	if req.Format != "" && filepath.Ext(output) == "" {
		output = output + "." + req.Format
	}

	// Add output filename as the last argument
	args = append(args, output)

	// Build the command string
	cmdString := fmt.Sprintf("ffmpeg %s", strings.Join(args, " "))

	// Log the command we're about to execute
	fmt.Printf("Executing: %s\n", cmdString)

	// If it's a dry run, just return the command string
	if req.DryRun {
		fmt.Printf("[Dry Run] %s\n", cmdString)
		return cmdString, nil
	}

	// Create a directory for the output file if it doesn't exist
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Execute ffmpeg with more detailed error and progress handling
	cmd := exec.Command("ffmpeg", args...)

	// Capture stderr to parse progress
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Process output to get progress
	go func() {
		scanner := bufio.NewScanner(stderr)
		var lastProgress *FFmpegProgress

		for scanner.Scan() {
			line := scanner.Text()

			// Parse progress information
			progress := ParseProgress(line, duration)

			// Only log meaningful progress updates
			if progress != nil && (lastProgress == nil ||
				progress.Percentage > lastProgress.Percentage+1.0 || // Report every 1% change
				progress.Frame > lastProgress.Frame+100) { // Or every 100 frames

				formattedProgress := FormatProgress(progress)
				fmt.Printf("Progress: %s\n", formattedProgress)
				lastProgress = progress
			}
		}
	}()

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("ffmpeg processing failed: %w", err)
	}

	fmt.Printf("Processing complete: %s\n", output)
	return output, nil
}

// CompareMedia compares original and processed media files
func CompareMedia(original, processed string) (CompareResult, error) {
	result := CompareResult{}

	// Get info for both files
	originalInfo, err := GetMediaInfo(original)
	if err != nil {
		return result, fmt.Errorf("failed to get original media info: %w", err)
	}

	processedInfo, err := GetMediaInfo(processed)
	if err != nil {
		return result, fmt.Errorf("failed to get processed media info: %w", err)
	}

	// Calculate size difference percentage
	originalSize := float64(originalInfo.Size)
	processedSize := float64(processedInfo.Size)

	sizeDiff := ((originalSize - processedSize) / originalSize) * 100

	// Determine if resolution changed
	result.ResolutionChanged = originalInfo.Resolution != processedInfo.Resolution

	// Determine if format changed
	result.FormatChanged = originalInfo.Format != processedInfo.Format

	// Determine if codec changed
	result.CodecChanged = originalInfo.Codec != processedInfo.Codec

	// Calculate bitrate reduction if available
	if originalInfo.Bitrate != "" && processedInfo.Bitrate != "" {
		// Parse bitrate values (removing units like 'kb/s')
		originalBitrateStr := strings.TrimRight(originalInfo.Bitrate, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ /")
		processedBitrateStr := strings.TrimRight(processedInfo.Bitrate, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ /")

		originalBitrateVal, err1 := strconv.ParseFloat(originalBitrateStr, 64)
		processedBitrateVal, err2 := strconv.ParseFloat(processedBitrateStr, 64)

		if err1 == nil && err2 == nil && originalBitrateVal > 0 {
			result.BitrateReduction = ((originalBitrateVal - processedBitrateVal) / originalBitrateVal) * 100
		}
	}

	result.Original = originalInfo
	result.Processed = processedInfo
	result.SizeDiff = sizeDiff

	return result, nil
}
