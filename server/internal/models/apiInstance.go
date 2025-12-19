package models

type RequestApiInstance struct {
	Api string
	Url string
}

type ApiInstance struct {
	ID     uint   `gorm:"primaryKey"`
	UserId uint   `gorm:"not null;uniqueIndex:idx_instance"`
	Api    string `gorm:"not null"`
	Url    string `gorm:"not null;uniqueIndex:idx_instance"`
}

type ApiInstanceItem struct {
	Id  uint
	Api string
	Url string
}
