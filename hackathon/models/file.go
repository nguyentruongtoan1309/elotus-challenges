package models

import (
	"database/sql"
	"time"
)

// FileMetadata represents uploaded file metadata
type FileMetadata struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	FilePath    string    `json:"file_path"`
	UserAgent   string    `json:"user_agent"`
	RemoteAddr  string    `json:"remote_addr"`
	CreatedAt   time.Time `json:"created_at"`
}

// FileModel handles file metadata database operations
type FileModel struct {
	DB *sql.DB
}

// NewFileModel creates a new FileModel
func NewFileModel(db *sql.DB) *FileModel {
	return &FileModel{DB: db}
}

// CreateTable creates the files table if it doesn't exist
func (m *FileModel) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		filename TEXT NOT NULL,
		content_type TEXT NOT NULL,
		size INTEGER NOT NULL,
		file_path TEXT NOT NULL,
		user_agent TEXT,
		remote_addr TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id)
	)`
	_, err := m.DB.Exec(query)
	return err
}

// Create stores file metadata in the database
func (m *FileModel) Create(metadata *FileMetadata) (*FileMetadata, error) {
	query := `
	INSERT INTO files (user_id, filename, content_type, size, file_path, user_agent, remote_addr)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(query,
		metadata.UserID,
		metadata.Filename,
		metadata.ContentType,
		metadata.Size,
		metadata.FilePath,
		metadata.UserAgent,
		metadata.RemoteAddr,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return m.GetByID(int(id))
}

// GetByID retrieves file metadata by ID
func (m *FileModel) GetByID(id int) (*FileMetadata, error) {
	metadata := &FileMetadata{}
	query := `
	SELECT id, user_id, filename, content_type, size, file_path, user_agent, remote_addr, created_at
	FROM files WHERE id = ?`

	err := m.DB.QueryRow(query, id).Scan(
		&metadata.ID,
		&metadata.UserID,
		&metadata.Filename,
		&metadata.ContentType,
		&metadata.Size,
		&metadata.FilePath,
		&metadata.UserAgent,
		&metadata.RemoteAddr,
		&metadata.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
