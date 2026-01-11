package models

type User struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Username  string      `gorm:"not null;uniqueIndex" json:"username"`
	Password  string      `gorm:"not null" json:"password"`
	HiRes     bool        `gorm:"default:false" json:"hiRes"`
	Sessions  UserSession `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"sessions"`
	Follows   Follow      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"follows"`
	Instances Instance    `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"instances"`
}

type RequestUserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	HiRes    bool   `json:"hiRes"`
}

type ResponseUser struct {
	Username string `json:"username"`
	HiRes    bool   `json:"hiRes"`
}
