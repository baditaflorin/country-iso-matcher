package gui

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

// Handler handles GUI requests
type Handler struct {
	logger *slog.Logger
	webDir string
}

// NewHandler creates a new GUI handler
func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
		webDir: "web",
	}
}

// ServeGUI serves the GUI static files
func (h *Handler) ServeGUI(w http.ResponseWriter, r *http.Request) {
	// Get the requested file path
	filePath := filepath.Join(h.webDir, r.URL.Path)

	// If path is a directory or empty, serve index.html
	if r.URL.Path == "/" || r.URL.Path == "" {
		filePath = filepath.Join(h.webDir, "index.html")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If file doesn't exist and it's not a specific file, serve index.html
		filePath = filepath.Join(h.webDir, "index.html")
	}

	http.ServeFile(w, r, filePath)
}
