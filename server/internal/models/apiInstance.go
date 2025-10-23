package models

type ApiInstance struct {
	ID  uint   `gorm:"primaryKey"`
	Api string `gorm:"not null"`
	Url string `gorm:"not null;unique"`
}

type ApiInstanceRequest struct {
	Api string
	Url string
}
