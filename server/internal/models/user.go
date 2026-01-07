package models

type User struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Username    string      `gorm:"not null;uniqueIndex" json:"username"`
	Password    string      `gorm:"not null" json:"password"`
	BestQuality bool        `gorm:"default:true" json:"bestQuality"`
	Sessions    UserSession `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"sessions"`
	Follows     Follow      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"follows"`
	Instances   ApiInstance `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"instances"`
}

type RequestUserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestUser struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	BestQuality bool   `json:"bestQuality"`
}

type ResponseUser struct {
	Username    string `json:"username"`
	BestQuality bool   `json:"bestQuality"`
}
