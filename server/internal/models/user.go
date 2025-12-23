package models

type User struct {
	ID        uint        `gorm:"primaryKey"`
	Username  string      `gorm:"not null;uniqueIndex"`
	Password  string      `gorm:"not null"`
	Sessions  UserSession `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Follows   Follow      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Instances ApiInstance `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

type RequestUser struct {
	Username string
	Password string
}
