package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	Email           string         `gorm:"size:255;not null;uniqueIndex:idx_users_email,where:deleted_at IS NULL" json:"email"`
	Phone           string         `gorm:"size:20;not null;uniqueIndex:idx_users_phone,where:deleted_at IS NULL" json:"phone"`
	PasswordHash    string         `gorm:"size:255;not null" json:"-"`
	PhotoProfileKey *string        `gorm:"size:500" json:"photo_profile_key"`
	IsActive        bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}
