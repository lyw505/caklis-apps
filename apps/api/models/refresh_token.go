package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AdminID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"admin_id"`
	TokenHash string     `gorm:"size:255;not null;index" json:"-"`
	UserAgent *string    `gorm:"size:500" json:"user_agent"`
	IPAddress *string    `gorm:"size:45" json:"ip_address"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	
	Admin Admin `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE" json:"admin,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
