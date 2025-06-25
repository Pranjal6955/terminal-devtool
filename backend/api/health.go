package api

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

// HealthCheckResponse represents the response structure for the health check endpoint
type HealthCheckResponse struct {
	Status          string            `json:"status"`
	Version         string            `json:"version"`
	FFmpegAvailable bool              `json:"ffmpeg_available"`
	FFmpegVersion   string            `json:"ffmpeg_version,omitempty"`
	Components      map[string]string `json:"components,omitempty"`
}

// HealthCheck verifies the availability of required tools and components
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Status:          "OK",
		Version:         "0.1.0",
		FFmpegAvailable: false,
		Components:      make(map[string]string),
	}

	// Check if FFmpeg is available
	ffmpegCmd := exec.Command("ffmpeg", "-version")
	ffmpegOutput, err := ffmpegCmd.Output()
	if err == nil {
		response.FFmpegAvailable = true

		// Extract version information
		version := extractFFmpegVersion(string(ffmpegOutput))
		response.FFmpegVersion = version

		// Check for specific codecs/components
		response.Components["libx264"] = checkFFmpegComponent("libx264")
		response.Components["libvpx"] = checkFFmpegComponent("libvpx")
		response.Components["libopus"] = checkFFmpegComponent("libopus")
	} else {
		response.Status = "Warning"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// extractFFmpegVersion extracts version from ffmpeg -version output
func extractFFmpegVersion(output string) string {
	// Just return the first line as the version
	for i, c := range output {
		if c == '\n' {
			return output[:i]
		}
	}
	return output
}

// checkFFmpegComponent checks if a specific component is available in FFmpeg
func checkFFmpegComponent(component string) string {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-encoders")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	if contains(string(output), component) {
		return "available"
	}
	return "not available"
}

// contains checks if a string contains another substring
func contains(s, substring string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substring))
}
