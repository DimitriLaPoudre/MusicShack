package models

import "time"

type Song struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UserId    uint   `gorm:"not null;uniqueIndex:idx_song" json:"userId"`
	Path      string `gorm:"not null" json:"path"`
	Isrc      string `gorm:"uniqueIndex:idx_song" json:"isrc"`
	UpdatedAt time.Time
}

type MetadataInfo struct {
	Title        string   `json:"title"`
	ReleaseDate  string   `json:"releaseDate"`
	TrackNumber  string   `json:"trackNumber"`
	VolumeNumber string   `json:"volumeNumber"`
	Explicit     string   `json:"explicit"`
	Isrc         string   `json:"isrc"`
	Album        string   `json:"album"`
	AlbumArtists []string `json:"albumArtists"`
	Artists      []string `json:"artists"`
	AlbumGain    string   `json:"albumGain"`
	AlbumPeak    string   `json:"albumPeak"`
	TrackGain    string   `json:"trackGain"`
	TrackPeak    string   `json:"trackPeak"`
}

type UploadSong struct {
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

type ResponseSong struct {
	ID           uint     `json:"id"`
	Duration     uint     `json:"duration"`
	Title        string   `json:"title"`
	ReleaseDate  string   `json:"releaseDate"`
	TrackNumber  uint     `json:"trackNumber"`
	VolumeNumber uint     `json:"volumeNumber"`
	Explicit     bool     `json:"explicit"`
	Isrc         string   `json:"isrc"`
	Album        string   `json:"album"`
	AlbumArtists []string `json:"albumArtists"`
	Artists      []string `json:"artists"`
	AlbumGain    float64  `json:"albumGain"`
	AlbumPeak    float64  `json:"albumPeak"`
	TrackGain    float64  `json:"trackGain"`
	TrackPeak    float64  `json:"trackPeak"`
}

type ResponseLibrary struct {
	Total  int            `json:"total"`
	Count  int            `json:"count"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Items  []ResponseSong `json:"items"`
}
