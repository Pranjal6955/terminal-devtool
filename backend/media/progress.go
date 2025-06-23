package media

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// FFmpegProgress holds the parsed progress information from FFmpeg output
type FFmpegProgress struct {
	Frame      int
	FPS        float64
	TotalSize  int64
	Time       time.Duration
	Bitrate    string
	Speed      string
	Percentage float64
}

// Regular expressions for parsing FFmpeg output
var (
	frameRegex    = regexp.MustCompile(`frame=\s*(\d+)`)
	fpsRegex      = regexp.MustCompile(`fps=\s*(\d+\.?\d*)`)
	sizeRegex     = regexp.MustCompile(`size=\s*(\d+)kB`)
	timeRegex     = regexp.MustCompile(`time=\s*(\d{2}):(\d{2}):(\d{2})\.(\d{2})`)
	bitrateRegex  = regexp.MustCompile(`bitrate=\s*(\d+\.?\d*\w+/s)`)
	speedRegex    = regexp.MustCompile(`speed=\s*(\d+\.?\d*x)`)
	durationRegex = regexp.MustCompile(`Duration:\s*(\d{2}):(\d{2}):(\d{2})\.(\d{2})`)
)

// ParseProgress parses FFmpeg output line and extracts progress information
func ParseProgress(line string, totalDuration time.Duration) *FFmpegProgress {
	progress := &FFmpegProgress{}

	// Extract frame count
	if matches := frameRegex.FindStringSubmatch(line); len(matches) > 1 {
		if frame, err := strconv.Atoi(matches[1]); err == nil {
			progress.Frame = frame
		}
	}

	// Extract FPS
	if matches := fpsRegex.FindStringSubmatch(line); len(matches) > 1 {
		if fps, err := strconv.ParseFloat(matches[1], 64); err == nil {
			progress.FPS = fps
		}
	}

	// Extract size
	if matches := sizeRegex.FindStringSubmatch(line); len(matches) > 1 {
		if size, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			progress.TotalSize = size * 1024 // Convert kB to bytes
		}
	}

	// Extract time
	if matches := timeRegex.FindStringSubmatch(line); len(matches) > 4 {
		h, _ := strconv.Atoi(matches[1])
		m, _ := strconv.Atoi(matches[2])
		s, _ := strconv.Atoi(matches[3])
		ms, _ := strconv.Atoi(matches[4])

		progress.Time = time.Duration(h)*time.Hour +
			time.Duration(m)*time.Minute +
			time.Duration(s)*time.Second +
			time.Duration(ms*10)*time.Millisecond

		// Calculate percentage if we have total duration
		if totalDuration > 0 {
			progress.Percentage = float64(progress.Time) / float64(totalDuration) * 100
			if progress.Percentage > 100 {
				progress.Percentage = 100
			}
		}
	}

	// Extract bitrate
	if matches := bitrateRegex.FindStringSubmatch(line); len(matches) > 1 {
		progress.Bitrate = matches[1]
	}

	// Extract speed
	if matches := speedRegex.FindStringSubmatch(line); len(matches) > 1 {
		progress.Speed = matches[1]
	}

	return progress
}

// ParseDuration extracts video duration from FFmpeg output
func ParseDuration(output string) time.Duration {
	if matches := durationRegex.FindStringSubmatch(output); len(matches) > 4 {
		h, _ := strconv.Atoi(matches[1])
		m, _ := strconv.Atoi(matches[2])
		s, _ := strconv.Atoi(matches[3])
		ms, _ := strconv.Atoi(matches[4])

		return time.Duration(h)*time.Hour +
			time.Duration(m)*time.Minute +
			time.Duration(s)*time.Second +
			time.Duration(ms*10)*time.Millisecond
	}

	return 0
}

// FormatProgress formats the progress information for display
func FormatProgress(progress *FFmpegProgress) string {
	parts := []string{}

	if progress.Percentage > 0 {
		parts = append(parts, fmt.Sprintf("%.1f%%", progress.Percentage))
	}

	if progress.Time > 0 {
		parts = append(parts, formatDuration(progress.Time))
	}

	if progress.Frame > 0 {
		parts = append(parts, fmt.Sprintf("%d frames", progress.Frame))
	}

	if progress.FPS > 0 {
		parts = append(parts, fmt.Sprintf("%.2f fps", progress.FPS))
	}

	if progress.Speed != "" {
		parts = append(parts, progress.Speed)
	}

	if progress.Bitrate != "" {
		parts = append(parts, progress.Bitrate)
	}

	if len(parts) == 0 {
		return "Processing..."
	}

	return strings.Join(parts, " | ")
}

// formatDuration formats time.Duration to HH:MM:SS.xx
func formatDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond

	return fmt.Sprintf("%02d:%02d:%02d.%02d", h, m, s, ms/10)
}
