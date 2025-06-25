package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/Promptzy/terminal-devtool/backend/media"
)

// Handler processes HTTP requests for the media API
type Handler struct {
	BaseDir string
}

// NewHandler creates a new API handler
func NewHandler(baseDir string) *Handler {
	return &Handler{
		BaseDir: baseDir,
	}
}

// ProcessMedia handles requests to process media files
func (h *Handler) ProcessMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req media.ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Resolve paths relative to base directory if not absolute
	if !filepath.IsAbs(req.Input) {
		req.Input = filepath.Join(h.BaseDir, req.Input)
	}

	if req.Output != "" && !filepath.IsAbs(req.Output) {
		req.Output = filepath.Join(h.BaseDir, req.Output)
	}

	// Process the media file
	outputPath, err := media.ProcessMedia(req)
	if err != nil {
		http.Error(w, "Processing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result
	response := map[string]string{
		"output": outputPath,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CompareMedia handles requests to compare original and processed media files
func (h *Handler) CompareMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Original  string `json:"original"`
		Processed string `json:"processed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Resolve paths relative to base directory if not absolute
	if !filepath.IsAbs(req.Original) {
		req.Original = filepath.Join(h.BaseDir, req.Original)
	}

	if !filepath.IsAbs(req.Processed) {
		req.Processed = filepath.Join(h.BaseDir, req.Processed)
	}

	// Compare the media files
	result, err := media.CompareMedia(req.Original, req.Processed)
	if err != nil {
		http.Error(w, "Comparison failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// CompressMedia handles requests to compress video files with a specific bitrate
func (h *Handler) CompressMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Input   string `json:"input"`
		Output  string `json:"output,omitempty"`
		Bitrate string `json:"bitrate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Input == "" {
		http.Error(w, "Input path is required", http.StatusBadRequest)
		return
	}

	if req.Bitrate == "" {
		http.Error(w, "Bitrate is required", http.StatusBadRequest)
		return
	}

	// Resolve paths relative to base directory if not absolute
	if !filepath.IsAbs(req.Input) {
		req.Input = filepath.Join(h.BaseDir, req.Input)
	}

	if req.Output != "" && !filepath.IsAbs(req.Output) {
		req.Output = filepath.Join(h.BaseDir, req.Output)
	}

	// Compress the media file
	err := media.CompressMedia(req.Input, req.Output, req.Bitrate)
	if err != nil {
		http.Error(w, "Compression failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result
	response := map[string]string{
		"output":  req.Output,
		"status":  "success",
		"message": "Video compressed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMediaInfo handles requests to get media file information
func (h *Handler) GetMediaInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}

	// Resolve path relative to base directory if not absolute
	filePath := path
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(h.BaseDir, filePath)
	}

	// Get media info
	info, err := media.GetMediaInfo(filePath)
	if err != nil {
		http.Error(w, "Failed to get media info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
