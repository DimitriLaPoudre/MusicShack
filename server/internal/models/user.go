package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

type UpdateUserRequest struct {
	Username string
	Password string
}
