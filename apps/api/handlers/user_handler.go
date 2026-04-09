package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aul-pkl/cakli/backend/config"
	"github.com/aul-pkl/cakli/backend/models"
	"github.com/aul-pkl/cakli/backend/services"
	"github.com/aul-pkl/cakli/backend/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	PhotoProfileURL *string   `json:"photo_profile_url"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	Password        string  `json:"password"`
	PhotoProfileKey *string `json:"photo_profile_key"`
}

type UpdateUserRequest struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Phone           *string `json:"phone"`
	Password        *string `json:"password"`
	PhotoProfileKey *string `json:"photo_profile_key"`
}

// ListUsers returns paginated list of users
func ListUsers(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	query := config.DB.Model(&models.User{})

	// Search
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Sort
	orderClause := fmt.Sprintf("%s %s", sortBy, strings.ToUpper(sortOrder))
	query = query.Order(orderClause)

	// Paginate
	var users []models.User
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data user")
	}

	// Transform to response
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = transformUserToResponse(user)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return utils.PaginatedResponse(c, "Daftar user berhasil diambil", userResponses, utils.MetaData{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	})
}

// GetUser returns user detail by ID
func GetUser(c fiber.Ctx) error {
	id := c.Params("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		return utils.NotFoundResponse(c, "User tidak ditemukan")
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.NotFoundResponse(c, "User tidak ditemukan")
	}

	return utils.SuccessResponse(c, "Detail user berhasil diambil", transformUserToResponse(user))
}

// CreateUser creates new user
func CreateUser(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	// Validate required fields
	details := validateUserRequest(req.Name, req.Email, req.Phone, req.Password)
	if len(details) > 0 {
		return utils.ValidationErrorResponse(c, details)
	}

	// Check unique constraints
	if err := checkUserUnique(req.Email, req.Phone, uuid.Nil); err != nil {
		return err
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal meng-hash password")
	}

	// Create user
	user := models.User{
		Name:            req.Name,
		Email:           req.Email,
		Phone:           req.Phone,
		PasswordHash:    passwordHash,
		PhotoProfileKey: req.PhotoProfileKey,
		IsActive:        true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat user")
	}

	return utils.CreatedResponse(c, "User berhasil ditambahkan", transformUserToResponse(user))
}

// UpdateUser updates existing user
func UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		return utils.NotFoundResponse(c, "User tidak ditemukan")
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.NotFoundResponse(c, "User tidak ditemukan")
	}

	var req UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		// Check unique email
		if *req.Email != user.Email {
			if err := checkUserUnique(*req.Email, "", userID); err != nil {
				return err
			}
		}
		user.Email = *req.Email
	}
	if req.Phone != nil {
		// Check unique phone
		if *req.Phone != user.Phone {
			if err := checkUserUnique("", *req.Phone, userID); err != nil {
				return err
			}
		}
		user.Phone = *req.Phone
	}
	if req.Password != nil && *req.Password != "" {
		passwordHash, err := utils.HashPassword(*req.Password)
		if err != nil {
			return utils.InternalErrorResponse(c, "Gagal meng-hash password")
		}
		user.PasswordHash = passwordHash
	}
	if req.PhotoProfileKey != nil {
		user.PhotoProfileKey = req.PhotoProfileKey
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal memperbarui user")
	}

	return utils.SuccessResponse(c, "User berhasil diperbarui", transformUserToResponse(user))
}

// DeleteUser soft deletes user
func DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		return utils.NotFoundResponse(c, "User tidak ditemukan")
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.NotFoundResponse(c, "User tidak ditemukan")
	}

	user.IsActive = false
	if err := config.DB.Delete(&user).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal menghapus user")
	}

	return utils.SuccessResponse(c, "User berhasil dihapus", nil)
}

// Helper functions

func transformUserToResponse(user models.User) UserResponse {
	response := UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Generate presigned URL for photo profile
	if user.PhotoProfileKey != nil && *user.PhotoProfileKey != "" {
		url := services.GeneratePresignedViewURLForKey(*user.PhotoProfileKey)
		response.PhotoProfileURL = &url
	}

	return response
}

func validateUserRequest(name, email, phone, password string) []utils.ErrorDetail {
	details := []utils.ErrorDetail{}

	if name == "" {
		details = append(details, utils.ErrorDetail{Field: "name", Message: "Nama wajib diisi"})
	}
	if email == "" {
		details = append(details, utils.ErrorDetail{Field: "email", Message: "Email wajib diisi"})
	}
	if phone == "" {
		details = append(details, utils.ErrorDetail{Field: "phone", Message: "No. telepon wajib diisi"})
	}
	if password == "" {
		details = append(details, utils.ErrorDetail{Field: "password", Message: "Password wajib diisi"})
	} else if len(password) < 8 {
		details = append(details, utils.ErrorDetail{Field: "password", Message: "Password minimal 8 karakter"})
	}

	return details
}

func checkUserUnique(email, phone string, excludeID uuid.UUID) error {
	details := []utils.ErrorDetail{}

	if email != "" {
		var count int64
		query := config.DB.Model(&models.User{}).Where("email = ?", email)
		if excludeID != uuid.Nil {
			query = query.Where("id != ?", excludeID)
		}
		query.Count(&count)
		if count > 0 {
			details = append(details, utils.ErrorDetail{Field: "email", Message: "Email sudah digunakan oleh user lain"})
		}
	}

	if phone != "" {
		var count int64
		query := config.DB.Model(&models.User{}).Where("phone = ?", phone)
		if excludeID != uuid.Nil {
			query = query.Where("id != ?", excludeID)
		}
		query.Count(&count)
		if count > 0 {
			details = append(details, utils.ErrorDetail{Field: "phone", Message: "No. telepon sudah digunakan oleh user lain"})
		}
	}

	if len(details) > 0 {
		return &fiber.Error{
			Code:    409,
			Message: "Data sudah terdaftar",
		}
	}

	return nil
}
