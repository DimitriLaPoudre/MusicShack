package hifiv2_2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func fetchSong(ctx context.Context, url string, id string) (songData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/info/?id="+id, nil)
	if err != nil {
		return songData{}, fmt.Errorf("fetchSong: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return songData{}, fmt.Errorf("fetchSong: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return songData{}, fmt.Errorf("fetchSong: %w", errors.New("http error "+strconv.FormatInt(400, 10)))
	}

	var data songData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return songData{}, fmt.Errorf("fetchSong: json.Decode: %w", err)
	}

	return data, nil
}

func getSong(ctx context.Context, instances []models.ApiInstance, id string) (songData, error) {
	type res struct {
		data songData
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchSong(ctx, url, id)
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
			return songData{}, ctx.Err()
		}
	}
	return songData{}, fmt.Errorf("getSong: %w", lastErr)
}

func (p *Hifi) Song(ctx context.Context, userId uint, id string) (models.SongData, error) {
	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.SongData{}, fmt.Errorf("Hifi.Song: %w", err)
	}

	data, err := getSong(ctx, instances, id)
	if err != nil {
		return models.SongData{}, fmt.Errorf("Hifi.Song: %w", err)
	}

	normalizeSongData := models.SongData{
		Api:          p.Name(),
		Id:           strconv.FormatUint(uint64(data.Data.Id), 10),
		Title:        data.Data.Title,
		Duration:     data.Data.Duration,
		ReplayGain:   data.Data.ReplayGain,
		Peak:         data.Data.Peak,
		ReleaseDate:  data.Data.ReleaseDate[:10],
		TrackNumber:  data.Data.TrackNumber,
		VolumeNumber: data.Data.VolumeNumber,
		AudioQuality: models.QualityHigh,
		Explicit:     data.Data.Explicit,
		Popularity:   data.Data.Popularity,
		Isrc:         data.Data.Isrc,
		Artists:      make([]models.SongDataArtist, 0),
		Album: models.SongDataAlbum{
			Id:       strconv.FormatUint(uint64(data.Data.Album.Id), 10),
			Title:    data.Data.Album.Title,
			CoverUrl: utils.GetImageURL(data.Data.Album.CoverUrl, 1280),
		},
	}

	for _, quality := range data.Data.MediaMetadata.Tags {
		switch quality {
		case "HIRES_LOSSLESS":
			normalizeSongData.AudioQuality = max(normalizeSongData.AudioQuality, models.QualityHiresLossless)
		case "LOSSLESS", "DOLBY_ATMOS":
			normalizeSongData.AudioQuality = max(normalizeSongData.AudioQuality, models.QualityLossless)
		}
	}

	for _, artist := range data.Data.Artists {
		normalizeSongData.Artists = append(normalizeSongData.Artists, models.SongDataArtist{
			Id:   strconv.FormatUint(uint64(artist.Id), 10),
			Name: artist.Name,
		})
	}

	return normalizeSongData, nil
}
