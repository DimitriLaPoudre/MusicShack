package models

import (
	"mime/multipart"
	"time"
)

type Song struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	UserId uint      `gorm:"not null;uniqueIndex" json:"userId"`
	Path   string    `gorm:"not null" json:"path"`
	Isrc   string    `gorm:"index" json:"isrc"`
	MTime  time.Time `json:"mTime"`
}

type RequestUploadSong struct {
	Cover        *multipart.FileHeader `form:"cover"`
	File         *multipart.FileHeader `form:"file"`
	Title        *string               `form:"title"`
	ReleaseDate  *string               `form:"releaseDate"`
	TrackNumber  *uint                 `form:"trackNumber"`
	VolumeNumber *uint                 `form:"volumeNumber"`
	Explicit     *string               `form:"explicit"`
	Isrc         *string               `form:"isrc"`
	Album        *string               `form:"album"`
	AlbumArtists *[]string             `form:"albumArtists"`
	Artists      *[]string             `form:"artists"`
	AlbumGain    *float64              `form:"albumGain"`
	AlbumPeak    *float64              `form:"albumPeak"`
	TrackGain    *float64              `form:"trackGain"`
	TrackPeak    *float64              `form:"trackPeak"`
	ExtraTags    map[string][]string   `form:"extraTags"`
}

type RequestEditSong struct {
	Cover        *multipart.FileHeader `form:"cover"`
	Title        *string               `form:"title"`
	ReleaseDate  *string               `form:"releaseDate"`
	TrackNumber  *uint                 `form:"trackNumber"`
	VolumeNumber *uint                 `form:"volumeNumber"`
	Explicit     *string               `form:"explicit"`
	Isrc         *string               `form:"isrc"`
	Album        *string               `form:"album"`
	AlbumArtists *[]string             `form:"albumArtists"`
	Artists      *[]string             `form:"artists"`
	AlbumGain    *float64              `form:"albumGain"`
	AlbumPeak    *float64              `form:"albumPeak"`
	TrackGain    *float64              `form:"trackGain"`
	TrackPeak    *float64              `form:"trackPeak"`
	ExtraTags    map[string][]string   `form:"extraTags"`
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
