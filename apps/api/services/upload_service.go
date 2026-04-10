package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aul-pkl/cakli/backend/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

var validFolders = map[string]bool{
	"drivers/photo-profile": true,
	"drivers/photo-ktp":     true,
	"drivers/photo-face":    true,
	"users/photo-profile":   true,
}

var validContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

type PresignedUploadResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	ExpiresIn int    `json:"expires_in"`
}

type PresignedViewResponse struct {
	ViewURL   string `json:"view_url"`
	ExpiresIn int    `json:"expires_in"`
}

// GeneratePresignedUploadURL generates presigned PUT URL for file upload
func GeneratePresignedUploadURL(filename, contentType, folder string) (*PresignedUploadResponse, error) {
	// Validate content type
	if !validContentTypes[contentType] {
		return nil, fmt.Errorf("invalid content type: only image/jpeg and image/png are supported")
	}

	// Validate folder
	if !validFolders[folder] {
		return nil, fmt.Errorf("invalid folder: must be one of drivers/photo-profile, drivers/photo-ktp, drivers/photo-face, users/photo-profile")
	}

	// Generate unique object key
	ext := filepath.Ext(filename)
	objectKey := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Generate presigned PUT URL (15 minutes expiry)
	ctx := context.Background()
	bucketName := "cakli"
	expiry := 15 * time.Minute

	presignedURL, err := config.GetMinioClient().PresignedPutObject(ctx, bucketName, objectKey, expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	// Replace internal Docker hostname with localhost for browser access
	urlString := presignedURL.String()
	urlString = strings.Replace(urlString, "minio:9000", "localhost:9000", 1)

	return &PresignedUploadResponse{
		UploadURL: urlString,
		ObjectKey: objectKey,
		ExpiresIn: 900, // 15 minutes in seconds
	}, nil
}

// GeneratePresignedViewURL generates presigned GET URL for viewing file
func GeneratePresignedViewURL(objectKey string) (*PresignedViewResponse, error) {
	if objectKey == "" {
		return nil, fmt.Errorf("object key is required")
	}

	ctx := context.Background()
	bucketName := "cakli"

	// Check if object exists
	_, err := config.GetMinioClient().StatObject(ctx, bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("file not found")
	}

	// Generate presigned GET URL (1 hour expiry)
	expiry := 1 * time.Hour
	presignedURL, err := config.GetMinioClient().PresignedGetObject(ctx, bucketName, objectKey, expiry, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	// Replace internal Docker hostname with localhost for browser access
	urlString := presignedURL.String()
	urlString = strings.Replace(urlString, "minio:9000", "localhost:9000", 1)

	return &PresignedViewResponse{
		ViewURL:   urlString,
		ExpiresIn: 3600, // 1 hour in seconds
	}, nil
}

// GeneratePresignedViewURLForKey generates presigned GET URL without checking existence
func GeneratePresignedViewURLForKey(objectKey string) string {
	if objectKey == "" {
		return ""
	}

	ctx := context.Background()
	bucketName := "cakli"
	expiry := 1 * time.Hour

	presignedURL, err := config.GetMinioClient().PresignedGetObject(ctx, bucketName, objectKey, expiry, nil)
	if err != nil {
		return ""
	}

	// Replace internal Docker hostname with localhost for browser access
	urlString := presignedURL.String()
	urlString = strings.Replace(urlString, "minio:9000", "localhost:9000", 1)

	return urlString
}

// ValidateContentType checks if content type is valid
func ValidateContentType(contentType string) bool {
	return validContentTypes[contentType]
}

// ValidateFolder checks if folder is valid
func ValidateFolder(folder string) bool {
	return validFolders[folder]
}

// GetFileExtension returns file extension from filename
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}
