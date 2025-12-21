package models

import "time"

type UserSession struct {
	ID        uint      `gorm:"primaryKey"`
	UserId    uint      `gorm:"not null;index"`
	Token     string    `gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
