package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"size:255;not null;index" json:"-"`
	UserAgent *string    `gorm:"size:500" json:"user_agent"`
	IPAddress *string    `gorm:"size:45" json:"ip_address"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (UserRefreshToken) TableName() string {
	return "user_refresh_tokens"
}
