package models

type RequestApiInstance struct {
	Api string `json:"api"`
	Url string `json:"url"`
}

type ApiInstance struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserId uint   `gorm:"not null;uniqueIndex:idx_instance" json:"userId"`
	Api    string `gorm:"not null" json:"api"`
	Url    string `gorm:"not null;uniqueIndex:idx_instance" json:"url"`
}

type ApiInstanceItem struct {
	Id  uint   `json:"id"`
	Api string `json:"api"`
	Url string `json:"url"`
}
