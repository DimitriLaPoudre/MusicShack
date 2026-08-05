package dto

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

type AlbumItem struct {
	ID                     uint64 `json:"id"`
	Title                  string `json:"title"`
	Duration               int    `json:"duration"`
	StreamReady            bool   `json:"streamReady"`
	PayToStream            bool   `json:"payToStream"`
	AdSupportedStreamReady bool   `json:"adSupportedStreamReady"`
	DjReady                bool   `json:"djReady"`
	StemReady              bool   `json:"stemReady"`
	StreamStartDate        string `json:"streamStartDate"`
	AllowStreaming         bool   `json:"allowStreaming"`
	PremiumStreamingOnly   bool   `json:"premiumStreamingOnly"`

	NumberOfTracks  int `json:"numberOfTracks"`
	NumberOfVideos  int `json:"numberOfVideos"`
	NumberOfVolumes int `json:"numberOfVolumes"`

	ReleaseDate string `json:"releaseDate"`
	Copyright   string `json:"copyright"`
	Type        string `json:"type"`

	// Version *string `json:"version"`

	URL          string `json:"url"`
	Cover        string `json:"cover"`
	VibrantColor string `json:"vibrantColor"`

	// VideoCover *string `json:"videoCover"`

	Explicit     bool     `json:"explicit"`
	UPC          string   `json:"upc"`
	Popularity   int      `json:"popularity"`
	AudioQuality string   `json:"audioQuality"`
	AudioModes   []string `json:"audioModes"`

	MediaMetadata struct {
		Tags []string `json:"tags"`
	} `json:"mediaMetadata"`

	Upload bool `json:"upload"`
	AI     bool `json:"ai"`

	Artist MiniArtistItem `json:"artist"` // struct in /album but not /search/?al

	Artists []MiniArtistItem `json:"artists"`
}

func (a AlbumItem) ToAlbumInfo(size int) model.AlbumInfo {
	releaseDate, err := time.Parse(ReleaseDateLayout, a.ReleaseDate)
	if err != nil {
		slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse album %s releaseDate %s", a.Title, a.ReleaseDate), slog.String("err", err.Error()))
	}

	audioQuality := LOW
	switch a.AudioQuality {
	case "HIGH":
		audioQuality = HIGH
	case "LOSSLESS":
		audioQuality = LOSSLESS
	}
	for _, quality := range a.MediaMetadata.Tags {
		switch quality {
		case "HIRES_LOSSLESS":
			audioQuality = HIRES
		case "LOSSLESS", "DOLBY_ATMOS":
			if audioQuality != HIRES {
				audioQuality = LOSSLESS
			}
		}
	}

	artists := []model.ArtistInfo{}
	for _, artist := range a.Artists {
		artists = append(artists, model.ArtistInfo{
			ID:   strconv.FormatUint(uint64(artist.ID), 10),
			Name: artist.Name,
		})
	}

	return model.AlbumInfo{
		ID:            strconv.FormatUint(uint64(a.ID), 10),
		Title:         a.Title,
		Duration:      a.Duration,
		ReleaseDate:   releaseDate,
		NumberTracks:  a.NumberOfTracks,
		NumberVolumes: a.NumberOfVolumes,
		CoverURL:      hifi_utils.GetImageURL(a.Cover, size),
		AudioQuality:  audioQuality,
		Popularity:    a.Popularity,
		Explicit:      a.Explicit,
		Artists:       artists,
	}
}
