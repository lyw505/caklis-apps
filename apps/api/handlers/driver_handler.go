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

type DriverResponse struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Email               string     `json:"email"`
	Phone               string     `json:"phone"`
	NIK                 string     `json:"nik"`
	BirthPlace          string     `json:"birth_place"`
	BirthDate           string     `json:"birth_date"`
	Bank                *BankInfo  `json:"bank,omitempty"`
	BankAccountNumber   *string    `json:"bank_account_number"`
	PhotoProfileURL     *string    `json:"photo_profile_url"`
	PhotoKTPURL         *string    `json:"photo_ktp_url"`
	PhotoFaceURL        *string    `json:"photo_face_url"`
	VerificationStatus  string     `json:"verification_status"`
	IsActive            bool       `json:"is_active"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type BankInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type CreateDriverRequest struct {
	Name                string  `json:"name"`
	Email               string  `json:"email"`
	Phone               string  `json:"phone"`
	Password            string  `json:"password"`
	NIK                 string  `json:"nik"`
	BirthPlace          string  `json:"birth_place"`
	BirthDate           string  `json:"birth_date"`
	BankID              *uint   `json:"bank_id"`
	BankAccountNumber   *string `json:"bank_account_number"`
	PhotoProfileKey     *string `json:"photo_profile_key"`
	PhotoKTPKey         *string `json:"photo_ktp_key"`
	PhotoFaceKey        *string `json:"photo_face_key"`
	VerificationStatus  string  `json:"verification_status"`
}

type UpdateDriverRequest struct {
	Name                *string `json:"name"`
	Email               *string `json:"email"`
	Phone               *string `json:"phone"`
	Password            *string `json:"password"`
	NIK                 *string `json:"nik"`
	BirthPlace          *string `json:"birth_place"`
	BirthDate           *string `json:"birth_date"`
	BankID              *uint   `json:"bank_id"`
	BankAccountNumber   *string `json:"bank_account_number"`
	PhotoProfileKey     *string `json:"photo_profile_key"`
	PhotoKTPKey         *string `json:"photo_ktp_key"`
	PhotoFaceKey        *string `json:"photo_face_key"`
	VerificationStatus  *string `json:"verification_status"`
}

type UpdateVerificationStatusRequest struct {
	VerificationStatus string `json:"verification_status"`
}

// ListDrivers returns paginated list of drivers
func ListDrivers(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	verificationStatus := c.Query("verification_status", "")
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

	query := config.DB.Model(&models.Driver{}).Preload("Bank")

	// Search
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR nik ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Filter by verification status
	if verificationStatus != "" {
		query = query.Where("verification_status = ?", verificationStatus)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Sort
	orderClause := fmt.Sprintf("%s %s", sortBy, strings.ToUpper(sortOrder))
	query = query.Order(orderClause)

	// Paginate
	var drivers []models.Driver
	if err := query.Offset(offset).Limit(limit).Find(&drivers).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data driver")
	}

	// Transform to response
	driverResponses := make([]DriverResponse, len(drivers))
	for i, driver := range drivers {
		driverResponses[i] = transformDriverToResponse(driver)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return utils.PaginatedResponse(c, "Daftar driver berhasil diambil", driverResponses, utils.MetaData{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	})
}

// GetDriver returns driver detail by ID
func GetDriver(c fiber.Ctx) error {
	id := c.Params("id")

	driverID, err := uuid.Parse(id)
	if err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	var driver models.Driver
	if err := config.DB.Preload("Bank").First(&driver, driverID).Error; err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	return utils.SuccessResponse(c, "Detail driver berhasil diambil", transformDriverToResponse(driver))
}

// CreateDriver creates new driver
func CreateDriver(c fiber.Ctx) error {
	var req CreateDriverRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	// Validate required fields
	details := validateDriverRequest(req.Name, req.Email, req.Phone, req.Password, req.NIK, req.BirthPlace, req.BirthDate)
	if len(details) > 0 {
		return utils.ValidationErrorResponse(c, details)
	}

	// Check unique constraints
	if err := checkDriverUnique(req.Email, req.Phone, req.NIK, uuid.Nil); err != nil {
		return err
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "birth_date", Message: "Format tanggal lahir tidak valid (gunakan YYYY-MM-DD)"},
		})
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal meng-hash password")
	}

	// Create driver
	driver := models.Driver{
		Name:                req.Name,
		Email:               req.Email,
		Phone:               req.Phone,
		PasswordHash:        passwordHash,
		NIK:                 req.NIK,
		BirthPlace:          req.BirthPlace,
		BirthDate:           birthDate,
		BankID:              req.BankID,
		BankAccountNumber:   req.BankAccountNumber,
		PhotoProfileKey:     req.PhotoProfileKey,
		PhotoKTPKey:         req.PhotoKTPKey,
		PhotoFaceKey:        req.PhotoFaceKey,
		VerificationStatus:  req.VerificationStatus,
		IsActive:            true,
	}

	if driver.VerificationStatus == "" {
		driver.VerificationStatus = "pending"
	}

	if err := config.DB.Create(&driver).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat driver")
	}

	// Load bank relation
	config.DB.Preload("Bank").First(&driver, driver.ID)

	return utils.CreatedResponse(c, "Driver berhasil ditambahkan", transformDriverToResponse(driver))
}

// UpdateDriver updates existing driver
func UpdateDriver(c fiber.Ctx) error {
	id := c.Params("id")

	driverID, err := uuid.Parse(id)
	if err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	var driver models.Driver
	if err := config.DB.First(&driver, driverID).Error; err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	var req UpdateDriverRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	// Update fields if provided
	if req.Name != nil {
		driver.Name = *req.Name
	}
	if req.Email != nil {
		// Check unique email
		if *req.Email != driver.Email {
			if err := checkDriverUnique(*req.Email, "", "", driverID); err != nil {
				return err
			}
		}
		driver.Email = *req.Email
	}
	if req.Phone != nil {
		// Check unique phone
		if *req.Phone != driver.Phone {
			if err := checkDriverUnique("", *req.Phone, "", driverID); err != nil {
				return err
			}
		}
		driver.Phone = *req.Phone
	}
	if req.Password != nil && *req.Password != "" {
		passwordHash, err := utils.HashPassword(*req.Password)
		if err != nil {
			return utils.InternalErrorResponse(c, "Gagal meng-hash password")
		}
		driver.PasswordHash = passwordHash
	}
	if req.NIK != nil {
		// Check unique NIK
		if *req.NIK != driver.NIK {
			if err := checkDriverUnique("", "", *req.NIK, driverID); err != nil {
				return err
			}
		}
		driver.NIK = *req.NIK
	}
	if req.BirthPlace != nil {
		driver.BirthPlace = *req.BirthPlace
	}
	if req.BirthDate != nil {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
				{Field: "birth_date", Message: "Format tanggal lahir tidak valid (gunakan YYYY-MM-DD)"},
			})
		}
		driver.BirthDate = birthDate
	}
	if req.BankID != nil {
		driver.BankID = req.BankID
	}
	if req.BankAccountNumber != nil {
		driver.BankAccountNumber = req.BankAccountNumber
	}
	// if req.PhotoProfileKey != nil {
	// 	driver.PhotoProfileKey = req.PhotoProfileKey
	// }

	if req.PhotoProfileKey != nil {
	if strings.Contains(*req.PhotoProfileKey, "http") {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "photo_profile_key", Message: "Gunakan object key, bukan URL"},
		})
	}
	driver.PhotoProfileKey = req.PhotoProfileKey
}

	if req.PhotoKTPKey != nil {
		driver.PhotoKTPKey = req.PhotoKTPKey
	}
	if req.PhotoFaceKey != nil {
		driver.PhotoFaceKey = req.PhotoFaceKey
	}
	if req.VerificationStatus != nil {
		driver.VerificationStatus = *req.VerificationStatus
	}

	if err := config.DB.Save(&driver).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal memperbarui driver")
	}

	// Load bank relation
	config.DB.Preload("Bank").First(&driver, driver.ID)

	return utils.SuccessResponse(c, "Driver berhasil diperbarui", transformDriverToResponse(driver))
}

// DeleteDriver soft deletes driver
func DeleteDriver(c fiber.Ctx) error {
	id := c.Params("id")

	driverID, err := uuid.Parse(id)
	if err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	var driver models.Driver
	if err := config.DB.First(&driver, driverID).Error; err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	driver.IsActive = true
	if err := config.DB.Unscoped().Delete(&driver).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal menghapus driver")
	}

	// 🔥 hapus file dulu
	if driver.PhotoProfileKey != nil && *driver.PhotoProfileKey != "" {
		services.DeleteFile(*driver.PhotoProfileKey)
	}
	if driver.PhotoKTPKey != nil && *driver.PhotoKTPKey != "" {
		services.DeleteFile(*driver.PhotoKTPKey)
	}
	if driver.PhotoFaceKey != nil && *driver.PhotoFaceKey != "" {
		services.DeleteFile(*driver.PhotoFaceKey)
	}

	return utils.SuccessResponse(c, "Driver berhasil dihapus", nil)
}

// UpdateVerificationStatus updates driver verification status
func UpdateVerificationStatus(c fiber.Ctx) error {
	id := c.Params("id")

	driverID, err := uuid.Parse(id)
	if err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	var driver models.Driver
	if err := config.DB.First(&driver, driverID).Error; err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	var req UpdateVerificationStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	// Validate status
	validStatuses := map[string]bool{"pending": true, "accepted": true, "rejected": true}
	if !validStatuses[req.VerificationStatus] {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "verification_status", Message: "Status harus salah satu dari: pending, accepted, rejected"},
		})
	}

	driver.VerificationStatus = req.VerificationStatus
	if err := config.DB.Save(&driver).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal memperbarui status verifikasi")
	}

	return utils.SuccessResponse(c, "Status verifikasi driver berhasil diperbarui", fiber.Map{
		"id":                  driver.ID,
		"name":                driver.Name,
		"verification_status": driver.VerificationStatus,
		"updated_at":          driver.UpdatedAt,
	})
}

// Helper functions

func transformDriverToResponse(driver models.Driver) DriverResponse {
	response := DriverResponse{
		ID:                 driver.ID.String(),
		Name:               driver.Name,
		Email:              driver.Email,
		Phone:              driver.Phone,
		NIK:                driver.NIK,
		BirthPlace:         driver.BirthPlace,
		BirthDate:          driver.BirthDate.Format("2006-01-02"),
		BankAccountNumber:  driver.BankAccountNumber,
		VerificationStatus: driver.VerificationStatus,
		IsActive:           driver.IsActive,
		CreatedAt:          driver.CreatedAt,
		UpdatedAt:          driver.UpdatedAt,
	}

	// Add bank info if exists
	if driver.Bank != nil {
		response.Bank = &BankInfo{
			ID:   driver.Bank.ID,
			Name: driver.Bank.Name,
			Code: driver.Bank.Code,
		}
	}

	// Generate presigned URLs for photos
	if driver.PhotoProfileKey != nil && *driver.PhotoProfileKey != "" {
		url := services.GeneratePresignedViewURLForKey(*driver.PhotoProfileKey)
		response.PhotoProfileURL = &url
	}
	if driver.PhotoKTPKey != nil && *driver.PhotoKTPKey != "" {
		url := services.GeneratePresignedViewURLForKey(*driver.PhotoKTPKey)
		response.PhotoKTPURL = &url
	}
	if driver.PhotoFaceKey != nil && *driver.PhotoFaceKey != "" {
		url := services.GeneratePresignedViewURLForKey(*driver.PhotoFaceKey)
		response.PhotoFaceURL = &url
	}

	return response
}

func validateDriverRequest(name, email, phone, password, nik, birthPlace, birthDate string) []utils.ErrorDetail {
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
	if nik == "" {
		details = append(details, utils.ErrorDetail{Field: "nik", Message: "NIK wajib diisi"})
	} else if len(nik) != 16 {
		details = append(details, utils.ErrorDetail{Field: "nik", Message: "NIK harus tepat 16 digit"})
	}
	if birthPlace == "" {
		details = append(details, utils.ErrorDetail{Field: "birth_place", Message: "Tempat lahir wajib diisi"})
	}
	if birthDate == "" {
		details = append(details, utils.ErrorDetail{Field: "birth_date", Message: "Tanggal lahir wajib diisi"})
	}

	return details
}

func checkDriverUnique(email, phone, nik string, excludeID uuid.UUID) error {
	details := []utils.ErrorDetail{}

	if email != "" {
		var count int64
		query := config.DB.Model(&models.Driver{}).Where("email = ?", email)
		if excludeID != uuid.Nil {
			query = query.Where("id != ?", excludeID)
		}
		query.Count(&count)
		if count > 0 {
			details = append(details, utils.ErrorDetail{Field: "email", Message: "Email sudah digunakan oleh driver lain"})
		}
	}

	if phone != "" {
		var count int64
		query := config.DB.Model(&models.Driver{}).Where("phone = ?", phone)
		if excludeID != uuid.Nil {
			query = query.Where("id != ?", excludeID)
		}
		query.Count(&count)
		if count > 0 {
			details = append(details, utils.ErrorDetail{Field: "phone", Message: "No. telepon sudah digunakan oleh driver lain"})
		}
	}

	if nik != "" {
		var count int64
		query := config.DB.Model(&models.Driver{}).Where("nik = ?", nik)
		if excludeID != uuid.Nil {
			query = query.Where("id != ?", excludeID)
		}
		query.Count(&count)
		if count > 0 {
			details = append(details, utils.ErrorDetail{Field: "nik", Message: "NIK sudah digunakan oleh driver lain"})
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
