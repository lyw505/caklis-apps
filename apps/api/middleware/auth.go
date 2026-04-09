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

	// Store admin ID and role in context
	c.Locals("admin_id", claims.AdminID)
	c.Locals("admin_role", claims.Role)

	return c.Next()
}

// GetAdminID retrieves admin ID from context
func GetAdminID(c fiber.Ctx) uuid.UUID {
	return c.Locals("admin_id").(uuid.UUID)
}

// GetAdminRole retrieves admin role from context
func GetAdminRole(c fiber.Ctx) string {
	return c.Locals("admin_role").(string)
}
