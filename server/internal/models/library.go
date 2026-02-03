package models

type Song struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserId uint   `gorm:"not null;uniqueIndex:idx_song" json:"userId"`
	Isrc   string `gorm:"uniqueIndex:idx_song" json:"isrc"`
	Path   string `gorm:"not null" json:"path"`
}

type ResponseSong struct {
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	Duration     uint     `json:"duration"`
	ReleaseDate  string   `json:"releaseDate"`
	TrackNumber  uint     `json:"trackNumber"`
	VolumeNumber uint     `json:"volumeNumber"`
	Explicit     bool     `json:"explicit"`
	Isrc         string   `json:"isrc"`
	Album        string   `json:"album"`
	Artists      []string `json:"artists"`
}
