package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

type UserRequest struct {
	Username string
	Password string
}
