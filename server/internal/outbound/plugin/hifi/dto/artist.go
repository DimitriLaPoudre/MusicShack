package dto

import (
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

type ArtistItem struct {
	ID                         uint64   `json:"id"`
	Name                       string   `json:"name"`
	ArtistTypes                []string `json:"artistTypes"`
	URL                        string   `json:"url"`
	Picture                    string   `json:"picture"`
	SelectedAlbumCoverFallback *string  `json:"selectedAlbumCoverFallback"`
	Popularity                 int      `json:"popularity"`

	ArtistRoles []struct {
		CategoryID int    `json:"categoryId"`
		Category   string `json:"category"`
	} `json:"artistRoles"`

	Mixes struct {
		ArtistMix string `json:"ARTIST_MIX"`
	} `json:"mixes"`

	// Handle *string `json:"handle"`
	// UserID *uint64 `json:"userId"`

	Spotlighted bool `json:"spotlighted"`
}

func (a ArtistItem) ToArtistInfo(size int) model.ArtistInfo {
	pictureURL := a.Picture
	if pictureURL == "" && a.SelectedAlbumCoverFallback != nil {
		pictureURL = *a.SelectedAlbumCoverFallback
	}
	pictureURL = hifi_utils.GetImageURL(pictureURL, size)

	return model.ArtistInfo{
		ID:         strconv.FormatUint(uint64(a.ID), 10),
		Name:       a.Name,
		PictureURL: pictureURL,
		Popularity: a.Popularity,
	}
}

type MiniArtistItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`

	// Handle *string `json:"handle"`

	Type string `json:"type"`

	Picture *string `json:"picture"`
}
