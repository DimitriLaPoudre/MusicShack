package models

type RequestApiInstance struct {
	Url string `json:"url"`
}

type ApiInstance struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserId uint   `gorm:"not null;uniqueIndex:idx_instance" json:"userId"`
	Api    string `gorm:"not null" json:"api"`
	Url    string `gorm:"not null;uniqueIndex:idx_instance" json:"url"`
	Ping   int64  `gorm:"" json:"ping"`
}

type ApiInstanceItem struct {
	Id   uint   `json:"id"`
	Api  string `json:"api"`
	Url  string `json:"url"`
	Ping int64  `json:"ping"`
}
