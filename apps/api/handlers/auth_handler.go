package handlers

import (
	"time"

	"github.com/aul-pkl/cakli/backend/config"
	"github.com/aul-pkl/cakli/backend/middleware"
	"github.com/aul-pkl/cakli/backend/models"
	"github.com/aul-pkl/cakli/backend/utils"
	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	ExpiresIn   int         `json:"expires_in"`
	TokenType   string      `json:"token_type"`
	Admin       AdminDetail `json:"admin"`
}

type AdminDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Login authenticates admin and returns access token
func Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
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
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	// Check if admin is active
	if !admin.IsActive {
		return utils.ForbiddenResponse(c, "Akun anda tidak aktif, hubungi administrator")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, admin.PasswordHash) {
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

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
		AdminID:   admin.ID,
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
		Admin: AdminDetail{
			ID:    admin.ID.String(),
			Name:  admin.Name,
			Email: admin.Email,
			Role:  admin.Role,
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

	// Get admin
	var admin models.Admin
	if err := config.DB.First(&admin, claims.UserID).Error; err != nil {
		return utils.UnauthorizedResponse(c, "Admin tidak ditemukan", "UNAUTHORIZED")
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(admin.ID, admin.Role)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat access token")
	}

	// Generate new refresh token (rotation)
	newRefreshToken, err := utils.GenerateRefreshToken(admin.ID)
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
		AdminID:   admin.ID,
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

type UserLoginResponse struct {
	AccessToken string      `json:"access_token"`
	ExpiresIn   int         `json:"expires_in"`
	TokenType   string      `json:"token_type"`
	User        UserDetail  `json:"user"`
}

type UserDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserLogin authenticates mobile user and returns access token
func UserLogin(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.ValidationErrorResponse(c, []utils.ErrorDetail{
			{Field: "body", Message: "Invalid JSON format"},
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
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
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	// Check if user is active
	if !user.IsActive {
		return utils.ForbiddenResponse(c, "Akun anda tidak aktif, hubungi administrator")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return utils.UnauthorizedResponse(c, "Email atau password salah", "INVALID_CREDENTIALS")
	}

	// Generate access token (Role is "user" for mobile users)
	accessToken, err := utils.GenerateAccessToken(user.ID, "user")
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat access token")
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "Gagal membuat refresh token")
	}

	// Save refresh token to database
	tokenHash := utils.HashToken(refreshToken)
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()
	expiryDuration, _ := time.ParseDuration("168h") // 7 days

	refreshTokenModel := models.UserRefreshToken{
		UserID:    user.ID,
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
		Name:     "user_refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/user/auth",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60,
	})

	// Return response
	return utils.SuccessResponse(c, "Login berhasil", UserLoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
		TokenType:   "Bearer",
		User: UserDetail{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	})
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

// GetCurrentUser returns current authenticated mobile user
func GetCurrentUser(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.NotFoundResponse(c, "User tidak ditemukan")
	}

	return utils.SuccessResponse(c, "Data user berhasil diambil", fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}
