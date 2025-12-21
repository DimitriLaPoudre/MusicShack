package models

type User struct {
	ID        uint        `gorm:"primaryKey"`
	Username  string      `gorm:"not null;uniqueIndex"`
	Password  string      `gorm:"not null"`
	Sessions  UserSession `gorm:"foreignKey:UserId"`
	Follows   Follow      `gorm:"foreignKey:UserId"`
	Instances ApiInstance `gorm:"foreignKey:UserId"`
}

type UserRequest struct {
	Username string
	Password string
}
