package media

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// CompressMedia compresses a video file using a user-defined bitrate
func CompressMedia(inputPath, outputPath, bitrate string) error {
	// Validate bitrate format
	if !isValidBitrate(bitrate) {
		return fmt.Errorf("invalid bitrate format '%s': must end with 'k' or 'M'", bitrate)
	}

	// If output path is not provided, generate one based on input
	if outputPath == "" {
		dir := filepath.Dir(inputPath)
		filename := filepath.Base(inputPath)
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		outputPath = filepath.Join(dir, fmt.Sprintf("%s_compressed%s", name, ext))
	}

	// Create directory for output file if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build the FFmpeg command
	args := []string{
		"-i", inputPath,
		"-b:v", bitrate,
		"-c:v", "libx264", // Use H.264 codec for compression
		"-preset", "medium", // Default preset for compression efficiency
		"-c:a", "copy", // Copy audio stream without re-encoding
		outputPath,
	}

	// Execute the FFmpeg command
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("compression failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Successfully compressed video to %s with bitrate %s\n", outputPath, bitrate)
	return nil
}

// isValidBitrate checks if the bitrate has the correct format
func isValidBitrate(bitrate string) bool {
	// Validate that bitrate ends with 'k' or 'M'
	regex := regexp.MustCompile(`^\d+[kM]$`)
	return regex.MatchString(bitrate)
}
