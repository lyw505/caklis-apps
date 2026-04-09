package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Email        string         `gorm:"size:255;not null;uniqueIndex:idx_admins_email,where:deleted_at IS NULL" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Role         string         `gorm:"size:50;not null;default:master_admin" json:"role"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Admin) TableName() string {
	return "admins"
}
