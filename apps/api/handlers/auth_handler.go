package handlers

import (
	"log"
	"time"

	"github.com/aul-pkl/cakli/backend/config"
	"github.com/aul-pkl/cakli/backend/middleware"
	"github.com/aul-pkl/cakli/backend/models"
	"github.com/aul-pkl/cakli/backend/services"
	"github.com/aul-pkl/cakli/backend/utils"
	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string        `json:"access_token"`
	ExpiresIn   int           `json:"expires_in"`
	TokenType   string        `json:"token_type"`
	Admin       *AdminDetail  `json:"admin,omitempty"`
	User        *UserDetail   `json:"user,omitempty"`
	Driver      *DriverDetail `json:"driver,omitempty"`
}

type AdminDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DriverDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Login authenticates admin and returns access token
func Login(c fiber.Ctx) error {
	log.Println("📥 Login request received")
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		log.Printf("❌ Bind error: %v", err)
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	log.Printf("👤 Attempting login for email: %s", req.Email)

	// Validate input
	if req.Email == "" || req.Password == "" {
		log.Println("❌ Missing email or password")
		details := []utils.ErrorDetail{}
		if req.Email == "" {
			details = append(details, utils.ErrorDetail{Field: "email", Message: "Email wajib diisi"})
		}
		if req.Password == "" {
			details = append(details, utils.ErrorDetail{Field: "password", Message: "Password wajib diisi"})
		}
		return utils.ValidationErrorResponse(c, details)
	}

	// Find admin by email
	var admin models.Admin
	if err := config.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		log.Printf("❌ Admin not found or DB error: %v", err)
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	// Check if admin is active
	if !admin.IsActive {
		log.Println("❌ Admin account is inactive")
		return utils.ForbiddenResponse(c, "Akun anda tidak aktif, hubungi administrator")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, admin.PasswordHash) {
		log.Println("❌ Invalid password")
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	log.Println("✅ Credentials verified, generating tokens...")

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(admin.ID, admin.Role)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat access token")
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(admin.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat refresh token")
	}

	// Save refresh token to database
	tokenHash := utils.HashToken(refreshToken)
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()
	expiryDuration, _ := time.ParseDuration("168h") // 7 days

	refreshTokenModel := models.RefreshToken{
		AdminID:   &admin.ID,
		TokenHash: tokenHash,
		UserAgent: &userAgent,
		IPAddress: &ipAddress,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	if err := config.DB.Create(&refreshTokenModel).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal menyimpan refresh token")
	}

	// Set refresh token as httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	// Return response
	return utils.SuccessResponse(c, "Login berhasil", LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   900, // 15 minutes
		TokenType:   "Bearer",
		Admin: &AdminDetail{
			ID:    admin.ID.String(),
			Name:  admin.Name,
			Email: admin.Email,
			Role:  admin.Role,
		},
	})
}

// LoginUser authenticates user and returns access token
func LoginUser(c fiber.Ctx) error {
	log.Println("📥 User Login request received")
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		log.Printf("❌ User Bind error: %v", err)
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	log.Printf("👤 Attempting user login for email: %s", req.Email)

	// Validate input
	if req.Email == "" || req.Password == "" {
		log.Println("❌ User login: Missing email or password")
		details := []utils.ErrorDetail{}
		if req.Email == "" {
			details = append(details, utils.ErrorDetail{Field: "email", Message: "Email wajib diisi"})
		}
		if req.Password == "" {
			details = append(details, utils.ErrorDetail{Field: "password", Message: "Password wajib diisi"})
		}
		return utils.ValidationErrorResponse(c, details)
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		log.Printf("❌ User not found or DB error: %v", err)
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	// Check if user is active
	if !user.IsActive {
		log.Println("❌ User account is inactive")
		return utils.ForbiddenResponse(c, "Akun anda tidak aktif, hubungi administrator")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		log.Println("❌ User login: Invalid password")
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	log.Println("✅ User credentials verified, generating tokens...")

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(user.ID, "user")
	if err != nil {
		log.Printf("❌ User login: Failed to generate access token: %v", err)
		return utils.InternalErrorResponse(c, "Gagal membuat access token")
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("❌ User login: Failed to generate refresh token: %v", err)
		return utils.InternalErrorResponse(c, "Gagal membuat refresh token")
	}

	// Save refresh token to database
	tokenHash := utils.HashToken(refreshToken)
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()
	expiryDuration, _ := time.ParseDuration("168h")

	refreshTokenModel := models.RefreshToken{
		UserID:    &user.ID,
		TokenHash: tokenHash,
		UserAgent: &userAgent,
		IPAddress: &ipAddress,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	if err := config.DB.Create(&refreshTokenModel).Error; err != nil {
		log.Printf("❌ User login: Failed to save refresh token: %v", err)
		return utils.InternalErrorResponse(c, "Gagal menyimpan refresh token")
	}

	log.Println("✅ User login successful")
	// Set refresh token as httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60,
	})

	// Return response
	return utils.SuccessResponse(c, "Login berhasil", LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
		TokenType:   "Bearer",
		User: &UserDetail{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// LoginDriver authenticates driver and returns access token
func LoginDriver(c fiber.Ctx) error {
	log.Println("📥 Driver Login request received")
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		log.Printf("❌ Driver Bind error: %v", err)
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	log.Printf("👤 Attempting driver login for email: %s", req.Email)

	// Validate input
	if req.Email == "" || req.Password == "" {
		log.Println("❌ Driver login: Missing email or password")
		details := []utils.ErrorDetail{}
		if req.Email == "" {
			details = append(details, utils.ErrorDetail{Field: "email", Message: "Email wajib diisi"})
		}
		if req.Password == "" {
			details = append(details, utils.ErrorDetail{Field: "password", Message: "Password wajib diisi"})
		}
		return utils.ValidationErrorResponse(c, details)
	}

	// Find driver by email
	var driver models.Driver
	if err := config.DB.Where("email = ?", req.Email).First(&driver).Error; err != nil {
		log.Printf("❌ Driver not found or DB error: %v", err)
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	// Check if driver is active
	if !driver.IsActive {
		log.Println("❌ Driver account is inactive")
		return utils.ForbiddenResponse(c, "Akun anda tidak aktif, hubungi administrator")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, driver.PasswordHash) {
		log.Println("❌ Driver login: Invalid password")
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	log.Println("✅ Driver credentials verified, generating tokens...")

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(driver.ID, "driver")
	if err != nil {
		log.Printf("❌ Driver login: Failed to generate access token: %v", err)
		return utils.InternalErrorResponse(c, "Gagal membuat access token")
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(driver.ID)
	if err != nil {
		log.Printf("❌ Driver login: Failed to generate refresh token: %v", err)
		return utils.InternalErrorResponse(c, "Gagal membuat refresh token")
	}

	// Save refresh token to database
	tokenHash := utils.HashToken(refreshToken)
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()
	expiryDuration, _ := time.ParseDuration("168h")

	refreshTokenModel := models.RefreshToken{
		DriverID:  &driver.ID,
		TokenHash: tokenHash,
		UserAgent: &userAgent,
		IPAddress: &ipAddress,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	if err := config.DB.Create(&refreshTokenModel).Error; err != nil {
		log.Printf("❌ Driver login: Failed to save refresh token: %v", err)
		return utils.InternalErrorResponse(c, "Gagal menyimpan refresh token")
	}

	log.Println("✅ Driver login successful")
	// Set refresh token as httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60,
	})

	// Return response
	return utils.SuccessResponse(c, "Login berhasil", LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
		TokenType:   "Bearer",
		Driver: &DriverDetail{
			ID:    driver.ID.String(),
			Name:  driver.Name,
			Email: driver.Email,
		},
	})
}

// Refresh generates new access token using refresh token
func Refresh(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return utils.UnauthorizedResponse(c, "Sesi tidak valid, silakan login kembali", "REFRESH_TOKEN_INVALID")
	}

	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Sesi sudah berakhir, silakan login kembali", "REFRESH_TOKEN_EXPIRED")
	}

	// Check if token exists and not revoked
	tokenHash := utils.HashToken(refreshToken)
	var storedToken models.RefreshToken
	if err := config.DB.Where("token_hash = ? AND revoked_at IS NULL", tokenHash).First(&storedToken).Error; err != nil {
		return utils.UnauthorizedResponse(c, "Sesi tidak valid, silakan login kembali", "REFRESH_TOKEN_INVALID")
	}

	// Get entity and generate new access token
	var role string

	if storedToken.AdminID != nil {
		var admin models.Admin
		if err := config.DB.First(&admin, storedToken.AdminID).Error; err != nil {
			return utils.UnauthorizedResponse(c, "Admin tidak ditemukan", "UNAUTHORIZED")
		}
		role = admin.Role
	} else if storedToken.UserID != nil {
		var user models.User
		if err := config.DB.First(&user, storedToken.UserID).Error; err != nil {
			return utils.UnauthorizedResponse(c, "User tidak ditemukan", "UNAUTHORIZED")
		}
		role = "user"
	} else if storedToken.DriverID != nil {
		var driver models.Driver
		if err := config.DB.First(&driver, storedToken.DriverID).Error; err != nil {
			return utils.UnauthorizedResponse(c, "Driver tidak ditemukan", "UNAUTHORIZED")
		}
		role = "driver"
	} else {
		return utils.UnauthorizedResponse(c, "Sesi tidak valid", "UNAUTHORIZED")
	}

	accessToken, err := utils.GenerateAccessToken(claims.UserID, role)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat access token")
	}

	// Generate new refresh token (rotation)
	newRefreshToken, err := utils.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat refresh token")
	}

	// Revoke old refresh token
	now := time.Now()
	storedToken.RevokedAt = &now
	config.DB.Save(&storedToken)

	// Save new refresh token
	newTokenHash := utils.HashToken(newRefreshToken)
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()
	expiryDuration, _ := time.ParseDuration("168h")

	newRefreshTokenModel := models.RefreshToken{
		AdminID:   storedToken.AdminID,
		UserID:    storedToken.UserID,
		DriverID:  storedToken.DriverID,
		TokenHash: newTokenHash,
		UserAgent: &userAgent,
		IPAddress: &ipAddress,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	if err := config.DB.Create(&newRefreshTokenModel).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal menyimpan refresh token")
	}

	// Set new refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60,
	})

	return utils.SuccessResponse(c, "Token berhasil diperbarui", fiber.Map{
		"access_token": accessToken,
		"expires_in":   900,
		"token_type":   "Bearer",
	})
}

// Logout revokes refresh token
func Logout(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken != "" {
		tokenHash := utils.HashToken(refreshToken)
		var storedToken models.RefreshToken
		if err := config.DB.Where("token_hash = ?", tokenHash).First(&storedToken).Error; err == nil {
			now := time.Now()
			storedToken.RevokedAt = &now
			config.DB.Save(&storedToken)
		}
	}

	// Clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		MaxAge:   -1,
	})

	return utils.SuccessResponse(c, "Logout berhasil", nil)
}

// GetCurrentAdmin returns current authenticated admin
func GetCurrentAdmin(c fiber.Ctx) error {
	adminID := middleware.GetAdminID(c)

	var admin models.Admin
	if err := config.DB.First(&admin, adminID).Error; err != nil {
		return utils.NotFoundResponse(c, "Admin tidak ditemukan")
	}

	return utils.SuccessResponse(c, "Data admin berhasil diambil", fiber.Map{
		"id":         admin.ID,
		"name":       admin.Name,
		"email":      admin.Email,
		"role":       admin.Role,
		"is_active":  admin.IsActive,
		"created_at": admin.CreatedAt,
		"updated_at": admin.UpdatedAt,
	})
}

// GetCurrentUser returns current authenticated user
func GetCurrentUser(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.NotFoundResponse(c, "User tidak ditemukan")
	}

	var photoProfileUrl *string
	if user.PhotoProfileKey != nil {
		resp, err := services.GeneratePresignedViewURL(*user.PhotoProfileKey)
		if err == nil && resp != nil {
			photoProfileUrl = &resp.ViewURL
		}
	}

	return utils.SuccessResponse(c, "Data user berhasil diambil", fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"photo_profile_url": photoProfileUrl,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// GetCurrentDriver returns current authenticated driver
func GetCurrentDriver(c fiber.Ctx) error {
	driverID := middleware.GetUserID(c)

	var driver models.Driver
	if err := config.DB.First(&driver, driverID).Error; err != nil {
		return utils.NotFoundResponse(c, "Driver tidak ditemukan")
	}

	return utils.SuccessResponse(c, "Data driver berhasil diambil", fiber.Map{
		"id":         driver.ID,
		"name":       driver.Name,
		"email":      driver.Email,
		"phone":      driver.Phone,
		"is_active":  driver.IsActive,
		"created_at": driver.CreatedAt,
		"updated_at": driver.UpdatedAt,
	})
}
