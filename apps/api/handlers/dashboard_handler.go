package handlers

import (
	"github.com/aul-pkl/cakli/backend/config"
	"github.com/aul-pkl/cakli/backend/models"
	"github.com/aul-pkl/cakli/backend/utils"
	"github.com/gofiber/fiber/v3"
)

// GetDashboardStats returns dashboard statistics
func GetDashboardStats(c fiber.Ctx) error {
	var stats struct {
		Drivers struct {
			Total    int64 `json:"total"`
			Verified int64 `json:"verified"`
			Pending  int64 `json:"pending"`
			Rejected int64 `json:"rejected"`
		} `json:"drivers"`
		Users struct {
			Total  int64 `json:"total"`
			Active int64 `json:"active"`
		} `json:"users"`
	}

	// Count drivers
	config.DB.Model(&models.Driver{}).Count(&stats.Drivers.Total)
	config.DB.Model(&models.Driver{}).Where("verification_status = ?", "accepted").Count(&stats.Drivers.Verified)
	config.DB.Model(&models.Driver{}).Where("verification_status = ?", "pending").Count(&stats.Drivers.Pending)
	config.DB.Model(&models.Driver{}).Where("verification_status = ?", "rejected").Count(&stats.Drivers.Rejected)

	// Count users
	config.DB.Model(&models.User{}).Count(&stats.Users.Total)
	config.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.Users.Active)

	return utils.SuccessResponse(c, "Statistik dashboard berhasil diambil", stats)
}
