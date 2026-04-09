package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Driver struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name                string         `gorm:"size:255;not null" json:"name"`
	Email               string         `gorm:"size:255;not null;uniqueIndex:idx_drivers_email,where:deleted_at IS NULL" json:"email"`
	Phone               string         `gorm:"size:20;not null;uniqueIndex:idx_drivers_phone,where:deleted_at IS NULL" json:"phone"`
	PasswordHash        string         `gorm:"size:255;not null" json:"-"`
	NIK                 string         `gorm:"size:16;not null;uniqueIndex:idx_drivers_nik,where:deleted_at IS NULL" json:"nik"`
	BirthPlace          string         `gorm:"size:100;not null" json:"birth_place"`
	BirthDate           time.Time      `gorm:"type:date;not null" json:"birth_date"`
	BankID              *uint          `gorm:"index" json:"bank_id"`
	BankAccountNumber   *string        `gorm:"size:30" json:"bank_account_number"`
	PhotoProfileKey     *string        `gorm:"size:500" json:"photo_profile_key"`
	PhotoKTPKey         *string        `gorm:"size:500" json:"photo_ktp_key"`
	PhotoFaceKey        *string        `gorm:"size:500" json:"photo_face_key"`
	VerificationStatus  string         `gorm:"size:20;not null;default:pending;index" json:"verification_status"`
	IsActive            bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	Bank *Bank `gorm:"foreignKey:BankID;constraint:OnDelete:SET NULL" json:"bank,omitempty"`
}

func (Driver) TableName() string {
	return "drivers"
}
