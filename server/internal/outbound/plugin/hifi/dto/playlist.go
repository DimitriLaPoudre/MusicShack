package dto

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

type PlaylistInfo struct {
	UUID           string `json:"uuid"`
	Title          string `json:"title"`
	NumberOfTracks int    `json:"numberOfTracks"`
	NumberOfVideos int    `json:"numberOfVideos"`
	Creator        struct {
		ID int `json:"id"`
	} `json:"creator"`
	Description     string           `json:"description"`
	Duration        int              `json:"duration"`
	LastUpdated     string           `json:"lastUpdated"`
	Created         string           `json:"created"`
	Type            string           `json:"type"`
	PublicPlaylist  bool             `json:"publicPlaylist"`
	URL             string           `json:"url"`
	Image           string           `json:"image"`
	Popularity      int              `json:"popularity"`
	SquareImage     string           `json:"squareImage"`
	PromotedArtists []MiniArtistItem `json:"promotedArtists"`
	LastItemAddedAt string           `json:"lastItemAddedAt"`
}

func (p PlaylistInfo) ToPlaylistInfo(size int) model.PlaylistInfo {
	lastUpdated, err := time.Parse(StreamStartDateLayout, p.LastUpdated)
	if err != nil {
		slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse playlist last update date: %s", p.LastUpdated), slog.String("err", err.Error()))
	}

	return model.PlaylistInfo{
		ID:             p.UUID,
		Title:          p.Title,
		Description:    p.Description,
		Duration:       p.Duration,
		NumberOfTracks: p.NumberOfTracks,
		CoverURL:       hifi_utils.GetImageURL(p.SquareImage, size),
		LastUpdated:    lastUpdated,
	}
}
