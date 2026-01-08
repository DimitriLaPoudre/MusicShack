package models

type RequestInstance struct {
	Url string `json:"url"`
}

type Instance struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserId   uint   `gorm:"not null;uniqueIndex:idx_instance" json:"userId"`
	Api      string `gorm:"not null" json:"api"`
	Provider string `gorm:"not null" json:"provider"`
	Url      string `gorm:"not null;uniqueIndex:idx_instance" json:"url"`
}

type InstanceItem struct {
	Id       uint   `json:"id"`
	Api      string `json:"api"`
	Provider string `json:"provider"`
	Url      string `json:"url"`
	Ping     int64  `json:"ping"`
}
