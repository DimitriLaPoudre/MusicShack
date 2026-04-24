package models

type RequestFollow struct {
	Provider string `json:"provider"`
	Id       string `json:"id"`
}

type Follow struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	UserId           uint   `gorm:"not null;uniqueIndex:idx_follow" json:"userId"`
	Provider         string `gorm:"not null;uniqueIndex:idx_follow" json:"provider"`
	ArtistId         string `gorm:"not null;uniqueIndex:idx_follow" json:"artistId"`
	ArtistName       string `gorm:"not null" json:"artistName"`
	ArtistPictureUrl string `gorm:"not null" json:"artistPictureUrl"`
}

type FollowItem struct {
	Id               uint   `json:"id"`
	Provider         string `json:"provider"`
	ArtistId         string `json:"artistId"`
	ArtistName       string `json:"artistName"`
	ArtistPictureUrl string `json:"artistPictureUrl"`
}
