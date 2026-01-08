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

func fetchAlbum(ctx context.Context, url string, id string) (albumData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/album/?id="+id, nil)
	if err != nil {
		return albumData{}, fmt.Errorf("fetchAlbum: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return albumData{}, fmt.Errorf("fetchAlbum: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return albumData{}, fmt.Errorf("fetchSong: %w", errors.New("http error "+strconv.FormatInt(400, 10)))
	}

	var data albumData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return albumData{}, fmt.Errorf("fetchAlbum: json.Decode: %w", err)
	}

	return data, nil
}

func getAlbum(ctx context.Context, instances []models.Instance, id string) (albumData, error) {
	type res struct {
		data albumData
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchAlbum(ctx, url, id)
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
			return albumData{}, ctx.Err()
		}
	}
	return albumData{}, fmt.Errorf("getAlbum: %w", lastErr)
}

func (p *Hifi) Album(ctx context.Context, userId uint, id string) (models.AlbumData, error) {
	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.AlbumData{}, fmt.Errorf("Hifi.Album: %w", err)
	}

	data, err := getAlbum(ctx, instances, id)
	if err != nil {
		return models.AlbumData{}, fmt.Errorf("Hifi.Album: %w", err)
	}

	normalizeAlbumData := models.AlbumData{
		Provider:      p.Provider(),
		Api:           p.Name(),
		Id:            strconv.FormatUint(uint64(data.Data.Id), 10),
		Title:         data.Data.Title,
		Duration:      data.Data.Duration,
		ReleaseDate:   data.Data.ReleaseDate,
		NumberTracks:  data.Data.NumberOfTracks,
		NumberVolumes: data.Data.NumberOfVolumes,
		AudioQuality:  models.QualityHigh,
		CoverUrl:      utils.GetImageURL(data.Data.CoverUrl, 1280),
		Explicit:      data.Data.Explicit,
		Songs:         make([]models.AlbumDataSong, 0),
		Artists:       make([]models.AlbumDataArtist, 0),
	}

	for _, quality := range data.Data.MediaMetadata.Tags {
		switch quality {
		case "HIRES_LOSSLESS":
			normalizeAlbumData.AudioQuality = max(normalizeAlbumData.AudioQuality, models.QualityHiresLossless)
		case "LOSSLESS", "DOLBY_ATMOS":
			normalizeAlbumData.AudioQuality = max(normalizeAlbumData.AudioQuality, models.QualityLossless)
		}
	}

	for _, artist := range data.Data.Artists {
		normalizeAlbumData.Artists = append(normalizeAlbumData.Artists,
			models.AlbumDataArtist{
				Id:   strconv.FormatUint(uint64(artist.Id), 10),
				Name: artist.Name,
			})
	}

	for _, item := range data.Data.Items {
		rawSong := item.Item

		song := models.AlbumDataSong{
			Id:           strconv.FormatUint(uint64(rawSong.Id), 10),
			Title:        rawSong.Title,
			Duration:     rawSong.Duration,
			TrackNumber:  rawSong.TrackNumber,
			VolumeNumber: rawSong.VolumeNumber,
			AudioQuality: models.QualityHigh,
			Explicit:     rawSong.Explicit,
			Artists:      make([]models.SongDataArtist, 0),
		}

		for _, quality := range rawSong.MediaMetadata.Tags {
			switch quality {
			case "HIRES_LOSSLESS":
				song.AudioQuality = max(song.AudioQuality, models.QualityHiresLossless)
			case "LOSSLESS", "DOLBY_ATMOS":
				song.AudioQuality = max(song.AudioQuality, models.QualityLossless)
			}
		}

		for _, artist := range rawSong.Artists {
			song.Artists = append(song.Artists, models.SongDataArtist{
				Id:   strconv.FormatUint(uint64(artist.Id), 10),
				Name: artist.Name,
			})
		}

		normalizeAlbumData.Songs = append(normalizeAlbumData.Songs, song)
	}

	return normalizeAlbumData, nil
}
