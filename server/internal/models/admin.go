package models

import "time"

type Admin struct {
	ID        bool      `gorm:"primaryKey;default:true"`
	Password  string    `gorm:"not null"`
	Token     string    `gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

type RequestAdmin struct {
	Password string
}

type RequestAdminPassword struct {
	OldPassword string
	NewPassword string
}
