package models

type User struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Username  string      `gorm:"not null;uniqueIndex" json:"username"`
	Password  string      `gorm:"not null" json:"password"`
	Sessions  UserSession `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"sessions"`
	Follows   Follow      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"follows"`
	Instances ApiInstance `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"instances"`
}

type RequestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
