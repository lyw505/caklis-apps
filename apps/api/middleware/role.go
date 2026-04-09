package middleware

import (
	"github.com/aul-pkl/cakli/backend/utils"
	"github.com/gofiber/fiber/v3"
)

// RequireRole checks if admin has one of the allowed roles
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role := GetAdminRole(c)

		// Check if role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return utils.ForbiddenResponse(c, "Anda tidak memiliki akses ke resource ini")
	}
}

// RequireMasterOrOperation allows both master_admin and operation_admin
func RequireMasterOrOperation() fiber.Handler {
	return RequireRole("master_admin", "operation_admin")
}

// RequireMasterOnly allows only master_admin
func RequireMasterOnly() fiber.Handler {
	return RequireRole("master_admin")
}

// RequireOperationOnly allows only operation_admin
func RequireOperationOnly() fiber.Handler {
	return RequireRole("operation_admin")
}

// RequireReportingOnly allows only reporting_admin
func RequireReportingOnly() fiber.Handler {
	return RequireRole("reporting_admin")
}
