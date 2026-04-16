package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AdminID   *uuid.UUID `gorm:"type:uuid;index" json:"admin_id,omitempty"`
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	DriverID  *uuid.UUID `gorm:"type:uuid;index" json:"driver_id,omitempty"`
	TokenHash string     `gorm:"size:255;not null;index" json:"-"`
	UserAgent *string    `gorm:"size:500" json:"user_agent"`
	IPAddress *string    `gorm:"size:45" json:"ip_address"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	
	Admin  *Admin  `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE" json:"admin,omitempty"`
	User   *User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Driver *Driver `gorm:"foreignKey:DriverID;constraint:OnDelete:CASCADE" json:"driver,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
