package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"file-uploader/models"
)

// UploadHandler handles file upload operations
type UploadHandler struct {
	fileModel *models.FileModel
}

// NewUploadHandler creates a new UploadHandler
func NewUploadHandler(fileModel *models.FileModel) *UploadHandler {
	return &UploadHandler{
		fileModel: fileModel,
	}
}

// UploadResponse represents the upload response
type UploadResponse struct {
	Message  string               `json:"message"`
	FileID   int                  `json:"file_id"`
	Metadata *models.FileMetadata `json:"metadata"`
}

// Upload handles file upload with validation
func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "User not authenticated"})
		return
	}

	// Parse multipart form with configurable max memory
	maxFileSize := getMaxFileSize()
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to parse multipart form"})
		return
	}

	// Get the file from form data
	file, fileHeader, err := r.FormFile("data")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No file provided or invalid file field name"})
		return
	}
	defer file.Close()

	if fileHeader.Size > maxFileSize {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: fmt.Sprintf("File size exceeds %d bytes limit", maxFileSize),
		})
		return
	}

	// Check content type is an image
	contentType := fileHeader.Header.Get("Content-Type")
	if !isImageContentType(contentType) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "File must be an image (JPEG, PNG, GIF, WebP, BMP, TIFF)"})
		return
	}

	uploadDir := getUploadDir()
	tempFileName := fmt.Sprintf("upload_%d_%d_%s", userID, time.Now().Unix(), fileHeader.Filename)
	tempFilePath := filepath.Join(uploadDir, tempFileName)

	// Create the temporary file
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create temporary file"})
		return
	}
	defer tempFile.Close()

	// Copy uploaded file content to temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		// Clean up the temporary file if copy fails
		os.Remove(tempFilePath)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to save file"})
		return
	}

	// Prepare file metadata
	metadata := &models.FileMetadata{
		UserID:      userID,
		Filename:    fileHeader.Filename,
		ContentType: contentType,
		Size:        fileHeader.Size,
		FilePath:    tempFilePath,
		UserAgent:   r.Header.Get("User-Agent"),
		RemoteAddr:  getClientIP(r),
	}

	// Save metadata to database
	savedMetadata, err := h.fileModel.Create(metadata)
	if err != nil {
		// Clean up the temporary file if database save fails
		os.Remove(tempFilePath)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to save file metadata"})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UploadResponse{
		Message:  "File uploaded successfully",
		FileID:   savedMetadata.ID,
		Metadata: savedMetadata,
	})
}

// isImageContentType checks if the content type is a valid image type
func isImageContentType(contentType string) bool {
	validImageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
		"image/tiff",
		"image/tif",
		"image/svg+xml",
	}

	contentType = strings.ToLower(contentType)
	for _, validType := range validImageTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// getMaxFileSize gets the max file size from environment variable
func getMaxFileSize() int64 {
	maxSizeStr := os.Getenv("MAX_UPLOAD_SIZE")
	if maxSizeStr == "" {
		return 8 << 20 // Default 8MB
	}

	maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64)
	if err != nil {
		return 8 << 20 // Default on error
	}

	return maxSize
}

// getUploadDir gets the upload directory from environment variable
func getUploadDir() string {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		return "/tmp" // Default fallback
	}
	return uploadDir
}
