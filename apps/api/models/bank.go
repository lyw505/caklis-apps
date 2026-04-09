package models

import "time"

type Bank struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Code      string    `gorm:"size:10;not null;uniqueIndex" json:"code"`
	LogoURL   *string   `gorm:"size:500" json:"logo_url"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (Bank) TableName() string {
	return "banks"
}
