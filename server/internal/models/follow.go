package models

type RequestFollow struct {
	Api string
	Id  string
}

type Follow struct {
	ID          uint   `gorm:"primaryKey"`
	UserId      uint   `gorm:"not null"`
	Api         string `gorm:"not null"`
	ArtistId    string `gorm:"not null"`
	LastFetchId string
}
