package handlers

import (
	"net/http"
	"os"
	"strconv"

	"file-uploader/models"

	"github.com/gorilla/mux"
)

// StaticHandler handles static file serving
type StaticHandler struct {
	fileModel *models.FileModel
}

// NewStaticHandler creates a new StaticHandler
func NewStaticHandler(fileModel *models.FileModel) *StaticHandler {
	return &StaticHandler{
		fileModel: fileModel,
	}
}

// ServeFile serves uploaded files with authentication
func (h *StaticHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	// Get file ID from URL
	vars := mux.Vars(r)
	fileIDStr := vars["fileId"]

	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get file metadata from database
	fileMetadata, err := h.fileModel.GetByID(fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Check if user owns the file
	if fileMetadata.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Check if file exists on disk
	if _, err := os.Stat(fileMetadata.FilePath); os.IsNotExist(err) {
		http.Error(w, "File not found on disk", http.StatusNotFound)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", fileMetadata.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+fileMetadata.Filename+"\"")

	// Serve the file
	http.ServeFile(w, r, fileMetadata.FilePath)
}

// ServePublicFile serves files without authentication (optional endpoint)
func (h *StaticHandler) ServePublicFile(w http.ResponseWriter, r *http.Request) {
	// Get file ID from URL
	vars := mux.Vars(r)
	fileIDStr := vars["fileId"]

	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Get file metadata from database
	fileMetadata, err := h.fileModel.GetByID(fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Check if file exists on disk
	if _, err := os.Stat(fileMetadata.FilePath); os.IsNotExist(err) {
		http.Error(w, "File not found on disk", http.StatusNotFound)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", fileMetadata.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+fileMetadata.Filename+"\"")

	// Serve the file
	http.ServeFile(w, r, fileMetadata.FilePath)
}
