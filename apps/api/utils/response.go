package utils

import "github.com/gofiber/fiber/v3"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
	Meta    *MetaData   `json:"meta,omitempty"`
}

type ErrorData struct {
	Code    string        `json:"code"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type MetaData struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// SuccessResponse sends success response
func SuccessResponse(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse sends created response
func CreatedResponse(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// PaginatedResponse sends paginated response
func PaginatedResponse(c fiber.Ctx, message string, data interface{}, meta MetaData) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

// ErrorResponse sends error response
func ErrorResponse(c fiber.Ctx, status int, message string, code string, details []ErrorDetail) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Message: message,
		Error: &ErrorData{
			Code:    code,
			Details: details,
		},
	})
}

// ValidationErrorResponse sends validation error response
func ValidationErrorResponse(c fiber.Ctx, details []ErrorDetail) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "Input tidak valid", "VALIDATION_ERROR", details)
}

// UnauthorizedResponse sends unauthorized response
func UnauthorizedResponse(c fiber.Ctx, message string, code string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message, code, nil)
}

// ForbiddenResponse sends forbidden response
func ForbiddenResponse(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message, "FORBIDDEN", nil)
}

// NotFoundResponse sends not found response
func NotFoundResponse(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message, "NOT_FOUND", nil)
}

// ConflictResponse sends conflict response
func ConflictResponse(c fiber.Ctx, details []ErrorDetail) error {
	return ErrorResponse(c, fiber.StatusConflict, "Data sudah terdaftar", "CONFLICT", details)
}

// InternalErrorResponse sends internal server error response
func InternalErrorResponse(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message, "INTERNAL_ERROR", nil)
}
