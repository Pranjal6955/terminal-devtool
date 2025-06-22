package media

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MediaInfo represents metadata about a media file
type MediaInfo struct {
	Filename   string `json:"filename"`
	Format     string `json:"format"`
	Duration   string `json:"duration"`
	Resolution string `json:"resolution"`
	Bitrate    string `json:"bitrate"`
	Size       int64  `json:"size"`
}

// ProcessRequest represents a request to process media
type ProcessRequest struct {
	Input      string `json:"input"`
	Output     string `json:"output,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	Bitrate    string `json:"bitrate,omitempty"`
	Format     string `json:"format,omitempty"`
}

// CompareResult represents the result of a media comparison
type CompareResult struct {
	Original MediaInfo `json:"original"`
	Processed MediaInfo `json:"processed"`
	SizeDiff float64 `json:"size_diff_percent"`
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
	
	// Run ffprobe to get media information
	cmd := exec.Command("ffprobe", 
		"-v", "error", 
		"-select_streams", "v:0", 
		"-show_entries", "stream=width,height,bit_rate:format=duration,format_name", 
		"-of", "csv=p=0", 
		filepath)
	
	output, err := cmd.Output()
	if err != nil {
		return info, fmt.Errorf("ffprobe failed: %w", err)
	}
	
	// Parse the output
	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(parts) >= 5 {
		width := parts[0]
		height := parts[1]
		info.Resolution = fmt.Sprintf("%sx%s", width, height)
		info.Bitrate = parts[2]
		info.Duration = parts[3]
		info.Format = parts[4]
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
	
	// Build ffmpeg command
	args := []string{"-i", req.Input}
	
	if req.Resolution != "" {
		args = append(args, "-s", req.Resolution)
	}
	
	if req.Bitrate != "" {
		args = append(args, "-b:v", req.Bitrate)
	}
	
	// Add output format if specified
	if req.Format != "" && filepath.Ext(output) == "" {
		output = output + "." + req.Format
	}
	
	args = append(args, output)
	
	// Execute ffmpeg
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg processing failed: %w", err)
	}
	
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
	
	result.Original = originalInfo
	result.Processed = processedInfo
	result.SizeDiff = sizeDiff
	
	return result, nil
}
