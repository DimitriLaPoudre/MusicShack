package models

type RequestFollow struct {
	Api string `json:"api"`
	Id  string `json:"id"`
}

type Follow struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	UserId           uint   `gorm:"not null;uniqueIndex:idx_follow" json:"userId"`
	Api              string `gorm:"not null;uniqueIndex:idx_follow" json:"api"`
	ArtistId         string `gorm:"not null;uniqueIndex:idx_follow" json:"artistId"`
	ArtistName       string `gorm:"not null" json:"artistName"`
	ArtistPictureUrl string `gorm:"not null" json:"artistPictureUrl"`
}

type FollowItem struct {
	Id               uint   `json:"id"`
	Api              string `json:"api"`
	ArtistId         string `json:"artistId"`
	ArtistName       string `json:"artistName"`
	ArtistPictureUrl string `json:"artistPictureUrl"`
}
