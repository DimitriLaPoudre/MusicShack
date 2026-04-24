package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/server/internal/plugins/hifi/utils"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func fetchPlaylist(ctx context.Context, apiURL string, id string) (playlistData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := utils.Fetch(ctx, apiURL+"/playlist/?id="+url.QueryEscape(id))
	if err != nil {
		return playlistData{}, fmt.Errorf("fetchPlaylist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return playlistData{}, fmt.Errorf("fetchPlaylist: http: %w", errors.New(resp.Status))
	}

	var data playlistData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return playlistData{}, fmt.Errorf("fetchPlaylist: json.Decode: %w", err)
	}

	return data, nil
}

func getPlaylist(ctx context.Context, instances []models.Instance, id string) (playlistData, error) {
	type res struct {
		data playlistData
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchPlaylist(ctx, url, id)
			ch <- res{data: data, err: err}
		}(instance.Url)
	}

	var lastErr error
	for range instances {
		select {
		case res := <-ch:
			if res.err == nil {
				return res.data, nil
			}
			lastErr = res.err
		case <-ctx.Done():
			return playlistData{}, ctx.Err()
		}
	}
	return playlistData{}, fmt.Errorf("getPlaylist: %w", lastErr)
}

func (p *Hifi) Playlist(ctx context.Context, userId uint, id string) (models.PlaylistData, error) {
	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.PlaylistData{}, fmt.Errorf("Hifi.Playlist: %w", err)
	}

	data, err := getPlaylist(ctx, instances, id)
	if err != nil {
		return models.PlaylistData{}, fmt.Errorf("Hifi.Playlist: %w", err)
	}

	normalizePlaylistData := models.PlaylistData{
		Provider:       p.Provider(),
		Api:            p.Name(),
		Id:             data.Playlist.UUID,
		Title:          data.Playlist.Title,
		Description:    data.Playlist.Description,
		Duration:       data.Playlist.Duration,
		NumberOfTracks: data.Playlist.NumberOfTracks,
		CoverURL:       hifi_utils.GetImageURL(data.Playlist.SquareImage, 640),
		LastUpdated:    data.Playlist.LastUpdated,
		Songs:          make([]models.PlaylistDataSong, 0),
	}

	for _, item := range data.Items {
		if item.Type != "track" {
			continue
		}
		song := item.Item

		newSong := models.PlaylistDataSong{
			Id:       strconv.FormatUint(uint64(song.ID), 10),
			Title:    song.Title,
			Duration: song.Duration,
			Explicit: song.Explicit,
			Isrc:     song.ISRC,
			Artists:  make([]models.SongDataArtist, 0),
		}

		switch song.AudioQuality {
		case "LOW":
			newSong.AudioQuality = LOW
		case "HIGH":
			newSong.AudioQuality = HIGH
		case "LOSSLESS":
			newSong.AudioQuality = LOSSLESS
		default:
			newSong.AudioQuality = LOW
		}
		for _, quality := range song.MediaMetadata.Tags {
			switch quality {
			case "HIRES_LOSSLESS":
				newSong.AudioQuality = HIRES
			case "LOSSLESS", "DOLBY_ATMOS":
				if newSong.AudioQuality != HIRES {
					newSong.AudioQuality = LOSSLESS
				}
			}
		}

		for _, artist := range song.Artists {
			newSong.Artists = append(newSong.Artists, models.SongDataArtist{
				Id:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		normalizePlaylistData.Songs = append(normalizePlaylistData.Songs, newSong)
	}

	return normalizePlaylistData, nil
}
