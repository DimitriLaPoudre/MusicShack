package models

type RequestFollow struct {
	Api string
	Id  string
}

type Follow struct {
	ID          uint   `gorm:"primaryKey"`
	UserId      uint   `gorm:"not null;uniqueIndex:idx_follow"`
	Api         string `gorm:"not null;uniqueIndex:idx_follow"`
	ArtistId    string `gorm:"not null;uniqueIndex:idx_follow"`
	LastFetchId string
}

type FollowItem struct {
	Id     uint
	Artist ArtistData
}
