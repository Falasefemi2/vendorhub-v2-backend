package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// SaveFile saves a file from multipart.FileHeader and returns the public URL
	SaveFile(ctx context.Context, file *multipart.FileHeader) (url string, err error)
	// DeleteFile removes a file
	DeleteFile(ctx context.Context, filename string) error
	// GetURL returns the full URL path for serving the file
	GetURL(filename string) string
}

// SupabaseStorage implements Storage interface using Supabase Storage
type SupabaseStorage struct {
	client      *supabase.Client
	bucket      string // Supabase bucket name
	supabaseURL string // Supabase project URL
	maxFileSize int64  // Max file size in bytes (default: 10MB)
}

// NewSupabaseStorage creates a new Supabase storage instance
func NewSupabaseStorage(url, key, bucket string) (*SupabaseStorage, error) {
	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Supabase client: %w", err)
	}

	return &SupabaseStorage{
		client:      client,
		bucket:      bucket,
		supabaseURL: strings.TrimSuffix(url, "/"),
		maxFileSize: 10 * 1024 * 1024, // 10MB default
	}, nil
}

// SaveFile saves a file to Supabase storage and returns the public URL
func (ss *SupabaseStorage) SaveFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is nil")
	}

	// Validate file size
	if file.Size > ss.maxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", ss.maxFileSize)
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

	// Upload to Supabase Storage
	_, err = ss.client.Storage.UploadFile(ss.bucket, filename, src)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to Supabase: %w", err)
	}

	// Log upload success (useful for debugging in production)
	fmt.Printf("info: uploaded file to Supabase bucket=%s filename=%s\n", ss.bucket, filename)

	// Return the FULL PUBLIC URL (not just the filename)
	publicURL := ss.GetURL(filename)
	fmt.Printf("info: public URL=%s\n", publicURL)
	return publicURL, nil
}

// DeleteFile removes a file from Supabase storage
func (ss *SupabaseStorage) DeleteFile(ctx context.Context, filename string) error {
	// Extract just the filename if full URL is passed
	if strings.Contains(filename, "/") {
		parts := strings.Split(filename, "/")
		filename = parts[len(parts)-1]
	}

	// Prevent directory traversal attacks
	if strings.Contains(filename, "..") {
		return fmt.Errorf("invalid filename")
	}

	// Delete from Supabase (RemoveFile accepts a slice of filenames)
	_, err := ss.client.Storage.RemoveFile(ss.bucket, []string{filename})
	if err != nil {
		// Log error but don't fail if file doesn't exist
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("file not found")
		}
		return fmt.Errorf("failed to delete file from Supabase: %w", err)
	}

	return nil
}

// GetURL returns the public URL for accessing the file
func (ss *SupabaseStorage) GetURL(filename string) string {
	// Extract just the filename if full URL is passed
	if strings.Contains(filename, "/") {
		parts := strings.Split(filename, "/")
		filename = parts[len(parts)-1]
	}

	// Supabase Storage uses this URL format for public files
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		ss.supabaseURL,
		ss.bucket,
		filename)
}

// generateUniqueFilename creates a unique filename with timestamp and UUID
func generateUniqueFilename(ext string) string {
	timestamp := time.Now().Unix()
	id := uuid.New().String()[:8] // Use first 8 chars of UUID
	return fmt.Sprintf("%d_%s%s", timestamp, id, ext)
}

// SetMaxFileSize sets the maximum allowed file size
func (ss *SupabaseStorage) SetMaxFileSize(size int64) {
	if size > 0 {
		ss.maxFileSize = size
	}
}

