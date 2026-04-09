package handlers

import (
	"github.com/aul-pkl/cakli/backend/config"
	"github.com/aul-pkl/cakli/backend/models"
	"github.com/aul-pkl/cakli/backend/utils"
	"github.com/gofiber/fiber/v3"
)

// ListBanks returns all active banks
func ListBanks(c fiber.Ctx) error {
	var banks []models.Bank

	if err := config.DB.Where("is_active = ?", true).Order("name ASC").Find(&banks).Error; err != nil {
		return utils.InternalErrorResponse(c, "Gagal mengambil data bank")
	}

	return utils.SuccessResponse(c, "Daftar bank berhasil diambil", banks)
}
