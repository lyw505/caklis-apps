package middleware

import (
	"strings"

	"github.com/aul-pkl/cakli/backend/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT access token
func AuthMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	
	if authHeader == "" {
		return utils.UnauthorizedResponse(c, "Token tidak ditemukan", "UNAUTHORIZED")
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return utils.UnauthorizedResponse(c, "Format token tidak valid", "TOKEN_INVALID")
	}

	tokenString := parts[1]

	// Validate token
	claims, err := utils.ValidateAccessToken(tokenString)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return utils.UnauthorizedResponse(c, "Token sudah expired, silakan refresh", "TOKEN_EXPIRED")
		}
		return utils.UnauthorizedResponse(c, "Token tidak valid", "TOKEN_INVALID")
	}

	// Store user/admin ID and role in context
	c.Locals("user_id", claims.UserID)
	c.Locals("user_role", claims.Role)

	return c.Next()
}

// GetUserID retrieves user/admin ID from context
func GetUserID(c fiber.Ctx) uuid.UUID {
	val := c.Locals("user_id")
	if val == nil {
		return uuid.Nil
	}
	return val.(uuid.UUID)
}

// GetUserRole retrieves user/admin role from context
func GetUserRole(c fiber.Ctx) string {
	val := c.Locals("user_role")
	if val == nil {
		return ""
	}
	return val.(string)
}

// Legacy helpers for Admin
func GetAdminID(c fiber.Ctx) uuid.UUID {
	return GetUserID(c)
}

func GetAdminRole(c fiber.Ctx) string {
	return GetUserRole(c)
}
