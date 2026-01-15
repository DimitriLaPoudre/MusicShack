package models

type Song struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Isrc    string `gorm:"uniqueIndex"`
	Path    string `gorm:"not null"`
	Quality string `gorm:""`
	UserId  uint   `gorm:"not null;uniqueIndex" json:"userId"`
}
