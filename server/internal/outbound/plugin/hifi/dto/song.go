package dto

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

type SongItem struct {
	ID                     uint64  `json:"id"`
	Title                  string  `json:"title"`
	Duration               int     `json:"duration"`
	ReplayGain             float64 `json:"replayGain"`
	Peak                   float64 `json:"peak"`
	AllowStreaming         bool    `json:"allowStreaming"`
	StreamReady            bool    `json:"streamReady"`
	PayToStream            bool    `json:"payToStream"`
	AdSupportedStreamReady bool    `json:"adSupportedStreamReady"`
	DjReady                bool    `json:"djReady"`
	StemReady              bool    `json:"stemReady"`
	StreamStartDate        string  `json:"streamStartDate"`
	PremiumStreamingOnly   bool    `json:"premiumStreamingOnly"`
	TrackNumber            int     `json:"trackNumber"`
	VolumeNumber           int     `json:"volumeNumber"`

	// Version *string `json:"version"`

	Popularity int    `json:"popularity"`
	Copyright  string `json:"copyright"`

	// BPM      *int   `json:"bpm"`
	// Key      *string `json:"key"`
	// KeyScale *string `json:"keyScale"`

	URL      string `json:"url"`
	ISRC     string `json:"isrc"`
	Editable bool   `json:"editable"`
	Explicit bool   `json:"explicit"`

	AudioQuality string   `json:"audioQuality"`
	AudioModes   []string `json:"audioModes"`

	MediaMetadata struct {
		Tags []string `json:"tags"`
	} `json:"mediaMetadata"`

	Upload      bool    `json:"upload"`
	AccessType  *string `json:"accessType"`
	Spotlighted bool    `json:"spotlighted"`
	AI          bool    `json:"ai"`

	Artist  MiniArtistItem   `json:"artist"`
	Artists []MiniArtistItem `json:"artists"`

	Album struct {
		ID           uint64 `json:"id"`
		Title        string `json:"title"`
		Cover        string `json:"cover"`
		VibrantColor string `json:"vibrantColor"`

		// VideoCover *string `json:"videoCover"`
	} `json:"album"`

	Mixes struct {
		TrackMix string `json:"TRACK_MIX"`
	} `json:"mixes"`
}

type SongItemTyped struct {
	Item SongItem `json:"item"`
	Type string   `json:"type"`
}

func (s SongItem) ToSongInfo(size int) model.SongInfo {
	releaseDate, err := time.Parse(StreamStartDateLayout, s.StreamStartDate)
	if err != nil {
		slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse song %s releaseDate %s", s.Title, s.StreamStartDate), slog.String("err", err.Error()))
	}

	audioQuality := LOW
	switch s.AudioQuality {
	case "HIGH":
		audioQuality = HIGH
	case "LOSSLESS":
		audioQuality = LOSSLESS
	}
	for _, quality := range s.MediaMetadata.Tags {
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
	for _, artist := range s.Artists {
		artists = append(artists, model.ArtistInfo{
			ID:   strconv.FormatUint(uint64(artist.ID), 10),
			Name: artist.Name,
		})
	}

	return model.SongInfo{
		ID:              strconv.FormatUint(uint64(s.ID), 10),
		Title:           s.Title,
		Duration:        s.Duration,
		Peak:            s.Peak,
		AlbumReplayGain: s.ReplayGain,
		ReleaseDate:     releaseDate,
		TrackNumber:     s.TrackNumber,
		VolumeNumber:    s.VolumeNumber,
		AudioQuality:    audioQuality,
		Explicit:        s.Explicit,
		Popularity:      s.Popularity,
		ISRC:            s.ISRC,
		Artists:         artists,
		Album: model.AlbumInfo{
			ID:       strconv.FormatUint(uint64(s.Album.ID), 10),
			Title:    s.Album.Title,
			CoverURL: hifi_utils.GetImageURL(s.Album.Cover, size),
		},
	}
}
