package handlers

import (
	"github.com/aul-pkl/cakli/backend/services"
	"github.com/aul-pkl/cakli/backend/utils"
	"github.com/gofiber/fiber/v3"
)

type GenerateUploadURLRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Folder      string `json:"folder"`
}

// GenerateUploadURL generates presigned PUT URL for file upload
func GenerateUploadURL(c fiber.Ctx) error {
	var req GenerateUploadURLRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	// Validate input
	details := []utils.ErrorDetail{}
	if req.Filename == "" {
		details = append(details, utils.ErrorDetail{Field: "filename", Message: "Filename wajib diisi"})
	}
	if req.ContentType == "" {
		details = append(details, utils.ErrorDetail{Field: "content_type", Message: "Content type wajib diisi"})
	} else if !services.ValidateContentType(req.ContentType) {
		details = append(details, utils.ErrorDetail{Field: "content_type", Message: "Hanya mendukung image/jpeg dan image/png"})
	}
	if req.Folder == "" {
		details = append(details, utils.ErrorDetail{Field: "folder", Message: "Folder wajib diisi"})
	} else if !services.ValidateFolder(req.Folder) {
		details = append(details, utils.ErrorDetail{
			Field:   "folder",
			Message: "Folder tidak valid. Pilihan: drivers/photo-profile, drivers/photo-ktp, drivers/photo-face, users/photo-profile",
		})
	}

	if len(details) > 0 {
		return utils.ValidationErrorResponse(c, details)
	}

	// Generate presigned URL
	result, err := services.GeneratePresignedUploadURL(req.Filename, req.ContentType, req.Folder)
	if err != nil {
		return utils.InternalErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Presigned URL berhasil dibuat", result)
}

// GenerateViewURL generates presigned GET URL for viewing file
func GenerateViewURL(c fiber.Ctx) error {
	objectKey := c.Query("key")

	if objectKey == "" {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "key", Message: "Object key wajib diisi"},
		})
	}

	// Generate presigned URL
	result, err := services.GeneratePresignedViewURL(objectKey)
	if err != nil {
		if err.Error() == "file not found" {
			return utils.NotFoundResponse(c, "File tidak ditemukan")
		}
		return utils.InternalErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Presigned view URL berhasil dibuat", result)
}
