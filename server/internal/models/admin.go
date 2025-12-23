package models

import "time"

type Admin struct {
	ID        bool      `gorm:"primaryKey;default:true" json:"id"`
	Password  string    `gorm:"not null" json:"password"`
	Token     string    `gorm:"size:64;uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
}

type RequestAdmin struct {
	Password string `json:"password"`
}

type RequestAdminPassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
