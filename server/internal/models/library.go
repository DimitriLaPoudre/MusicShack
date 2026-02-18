package models

import (
	"mime/multipart"
	"time"
)

type Song struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	UserId uint      `gorm:"not null" json:"userId"`
	Path   string    `gorm:"not null" json:"path"`
	Isrc   string    `gorm:"index" json:"isrc"`
	MTime  time.Time `json:"mTime"`
}

type RequestEditSong struct {
	Cover        *multipart.FileHeader `form:"cover"`
	Title        *string               `form:"title"`
	ReleaseDate  *string               `form:"releaseDate"`
	TrackNumber  *uint                 `form:"trackNumber"`
	VolumeNumber *uint                 `form:"volumeNumber"`
	Explicit     *bool                 `form:"explicit"`
	Isrc         *string               `form:"isrc"`
	Album        *string               `form:"album"`
	AlbumArtists *[]string             `form:"albumArtists"`
	Artists      *[]string             `form:"artists"`
	AlbumGain    *float64              `form:"albumGain"`
	AlbumPeak    *float64              `form:"albumPeak"`
	TrackGain    *float64              `form:"trackGain"`
	TrackPeak    *float64              `form:"trackPeak"`
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

type RequestUploadSong struct {
	Cover        *multipart.FileHeader `form:"cover"`
	File         *multipart.FileHeader `form:"file"`
	Title        *string               `form:"title"`
	ReleaseDate  *string               `form:"releaseDate"`
	TrackNumber  *uint                 `form:"trackNumber"`
	VolumeNumber *uint                 `form:"volumeNumber"`
	Explicit     *bool                 `form:"explicit"`
	Isrc         *string               `form:"isrc"`
	Album        *string               `form:"album"`
	AlbumArtists *[]string             `form:"albumArtists"`
	Artists      *[]string             `form:"artists"`
	AlbumGain    *float64              `form:"albumGain"`
	AlbumPeak    *float64              `form:"albumPeak"`
	TrackGain    *float64              `form:"trackGain"`
	TrackPeak    *float64              `form:"trackPeak"`
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
