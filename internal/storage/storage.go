package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// SaveFile saves a file from multipart.FileHeader and returns the filename/URL
	SaveFile(ctx context.Context, file *multipart.FileHeader) (filename string, err error)
	// DeleteFile removes a file
	DeleteFile(ctx context.Context, filename string) error
	// GetURL returns the full URL path for serving the file
	GetURL(filename string) string
}

// LocalStorage implements Storage interface using local filesystem
type LocalStorage struct {
	uploadDir   string // Directory to store uploads
	baseURL     string // Base URL for serving files (e.g., http://localhost:8080/uploads)
	maxFileSize int64  // Max file size in bytes (default: 10MB)
}

// NewLocalStorage creates a new local filesystem storage
func NewLocalStorage(uploadDir string, baseURL string) (*LocalStorage, error) {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &LocalStorage{
		uploadDir:   uploadDir,
		baseURL:     strings.TrimSuffix(baseURL, "/"), // Remove trailing slash
		maxFileSize: 10 * 1024 * 1024,                 // 10MB default
	}, nil
}

// SaveFile saves a file from multipart upload
func (ls *LocalStorage) SaveFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is nil")
	}

	// Validate file size
	if file.Size > ls.maxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", ls.maxFileSize)
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return "", fmt.Errorf("file type %s not allowed. Allowed types: jpg, jpeg, png, gif, webp", ext)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	filename := generateUniqueFilename(ext)
	filepath := filepath.Join(ls.uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		// Clean up on error
		os.Remove(filepath)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filename, nil
}

// DeleteFile removes a file from storage
func (ls *LocalStorage) DeleteFile(ctx context.Context, filename string) error {
	// Prevent directory traversal attacks
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("invalid filename")
	}

	filepath := filepath.Join(ls.uploadDir, filename)
	if err := os.Remove(filepath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found")
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetURL returns the full URL for accessing the file
func (ls *LocalStorage) GetURL(filename string) string {
	return fmt.Sprintf("%s/%s", ls.baseURL, filename)
}

// generateUniqueFilename creates a unique filename with timestamp and UUID
func generateUniqueFilename(ext string) string {
	timestamp := time.Now().Unix()
	id := uuid.New().String()[:8] // Use first 8 chars of UUID
	return fmt.Sprintf("%d_%s%s", timestamp, id, ext)
}

// SetMaxFileSize sets the maximum allowed file size
func (ls *LocalStorage) SetMaxFileSize(size int64) {
	if size > 0 {
		ls.maxFileSize = size
	}
}
