package models

type RequestFollow struct {
	Api string `json:"api"`
	Id  string `json:"id"`
}

type Follow struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserId   uint   `gorm:"not null;uniqueIndex:idx_follow" json:"userId"`
	Api      string `gorm:"not null;uniqueIndex:idx_follow" json:"api"`
	ArtistId string `gorm:"not null;uniqueIndex:idx_follow" json:"artistId"`
}

type FollowItem struct {
	Id     uint       `json:"id"`
	Api    string     `json:"api"`
	Artist ArtistData `json:"artist"`
}
