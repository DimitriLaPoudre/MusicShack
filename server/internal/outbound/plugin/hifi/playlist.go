package hifi

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getPlaylist(ctx context.Context, urls []string, id string) (playlistResponse, error) {
	playlist, err := hifi_utils.CachedMultiFetchTyped[playlistResponse](ctx, urls, "/playlist/?id="+url.QueryEscape(id), &p.limiters, &p.cache)
	if err != nil {
		return playlistResponse{}, fmt.Errorf("fetch playlist info with url list: %w", err)
	}

	return playlist, nil
}

func (p *Hifi) Playlist(ctx context.Context, instances []model.Instance, id string) (model.Playlist, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	data, err := p.getPlaylist(ctx, urls, id)
	if err != nil {
		return model.Playlist{}, err
	}

	lastUpdated, err := time.Parse(StreamStartDateLayout, data.Playlist.LastUpdated)
	if err != nil {
		slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse playlist last update date: %s", data.Playlist.LastUpdated), slog.String("err", err.Error()))
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
				ID:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		songs = append(songs, model.Song{
			ID:       strconv.FormatUint(uint64(song.ID), 10),
			Title:    song.Title,
			Duration: song.Duration,
			Explicit: song.Explicit,
			Isrc:     song.ISRC,
			Artists:  artists,
		})
	}

	return model.Playlist{
		ID:             data.Playlist.UUID,
		Title:          data.Playlist.Title,
		Description:    data.Playlist.Description,
		Duration:       data.Playlist.Duration,
		NumberOfTracks: data.Playlist.NumberOfTracks,
		CoverURL:       hifi_utils.GetImageURL(data.Playlist.SquareImage, 640),
		LastUpdated:    lastUpdated,
		Songs:          songs,
	}, nil
}
