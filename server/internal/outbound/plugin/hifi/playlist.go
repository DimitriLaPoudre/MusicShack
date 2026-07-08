package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	lib_url "net/url"
	"strconv"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	hifi_utils "github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/plugin/hifi/utils"
	"golang.org/x/time/rate"
)

func getPlaylist(ctx context.Context, limiter *rate.Limiter, url string, id string) (playlistData, error) {
	playlist, err := hifi_utils.FetchType[playlistData](ctx, url+"/playlist/?id="+lib_url.QueryEscape(id), limiter)
	if err != nil {
		return playlistData{}, fmt.Errorf("getPlaylist: %w", err)
	}

	return playlist, nil
}

func (p *Hifi) Playlist(ctx context.Context, url string, id string) (model.Playlist, error) {
	data, err := getPlaylist(ctx, p.limiter, url, id)
	if err != nil {
		return model.Playlist{}, fmt.Errorf("Hifi.Playlist: %w", err)
	}

	lastUpdated, err := time.Parse(StreamStartDateLayout, data.Playlist.LastUpdated)
	if err != nil {
		p.l.Warn().Msg(fmt.Sprintf("Hifi.Playlist: time.Parse(%s): %s", data.Playlist.LastUpdated, err.Error()))
	}

	songs := []model.Song{}
	for _, item := range data.Items {
		if item.Type != "track" {
			continue
		}
		song := item.Item

		audioQuality := LOW
		switch song.AudioQuality {
		case "LOW":
			audioQuality = LOW
		case "HIGH":
			audioQuality = HIGH
		case "LOSSLESS":
			audioQuality = LOSSLESS
		default:
			audioQuality = LOW
		}
		for _, quality := range song.MediaMetadata.Tags {
			switch quality {
			case "HIRES_LOSSLESS":
				audioQuality = HIRES
			case "LOSSLESS", "DOLBY_ATMOS":
				if audioQuality != HIRES {
					audioQuality = LOSSLESS
				}
			}
		}

		artists := []model.Artist{}
		for _, artist := range song.Artists {
			artists = append(artists, model.Artist{
				Id:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		songs = append(songs, model.Song{
			Id:       strconv.FormatUint(uint64(song.ID), 10),
			Title:    song.Title,
			Duration: song.Duration,
			Explicit: song.Explicit,
			Isrc:     song.ISRC,
			Artists:  artists,
		})
	}

	return model.Playlist{
		Id:             data.Playlist.UUID,
		Title:          data.Playlist.Title,
		Description:    data.Playlist.Description,
		Duration:       data.Playlist.Duration,
		NumberOfTracks: data.Playlist.NumberOfTracks,
		CoverURL:       hifi_utils.GetImageURL(data.Playlist.SquareImage, 640),
		LastUpdated:    lastUpdated,
		Songs:          songs,
	}, nil
}
