package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	lib_url "net/url"
	"strconv"
	"sync"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	hifi_utils "github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/plugin/hifi/utils"
	"golang.org/x/time/rate"
)

func fetchSong(ctx context.Context, limiter *rate.Limiter, url string, id string) (songData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := hifi_utils.Fetch(ctx, url+"/info/?id="+lib_url.QueryEscape(id), limiter)
	if err != nil {
		return songData{}, fmt.Errorf("fetchSong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return songData{}, fmt.Errorf("fetchSong: http: %w", errors.New(resp.Status))
	}

	var data songData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return songData{}, fmt.Errorf("fetchSong: json.Decode: %w", err)
	}

	return data, nil
}

func getSong(ctx context.Context, limiter *rate.Limiter, url string, id string) (songData, downloadData, error) {
	var songInfo songData
	var songInfoErr error
	var downloadInfo downloadData
	var downloadInfoErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		songInfo, songInfoErr = fetchSong(ctx, limiter, url, id)
	}()
	go func() {
		defer wg.Done()
		// downloadInfo, downloadInfoErr := getDownloadInfo(ctx, url, id, "")
	}()
	wg.Wait()

	if songInfoErr != nil {
		return songData{}, downloadData{}, fmt.Errorf("getSongData: %w", songInfoErr)
	}
	if downloadInfoErr != nil {
		return songData{}, downloadData{}, fmt.Errorf("getSongData: %w", downloadInfoErr)
	}

	return songInfo, downloadInfo, nil
}

func (p *Hifi) Song(ctx context.Context, url string, id string) (model.Song, error) {
	songInfo, downloadInfo, err := getSong(ctx, p.limiter, url, id)
	if err != nil {
		return model.Song{}, fmt.Errorf("Hifi.Song: %w", err)
	}

	releaseDate, err := time.Parse(StreamStartDateLayout, songInfo.Data.ReleaseDate)
	if err != nil {
		return model.Song{}, fmt.Errorf("Hifi.Song: time.Parse: %w", err)
	}

	audioQuality := LOW
	switch songInfo.Data.AudioQuality {
	case "LOW":
		audioQuality = LOW
	case "HIGH":
		audioQuality = HIGH
	case "LOSSLESS":
		audioQuality = LOSSLESS
	}
	for _, quality := range songInfo.Data.MediaMetadata.Tags {
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
	for _, artist := range songInfo.Data.Artists {
		artists = append(artists, model.Artist{
			Id:   strconv.FormatUint(uint64(artist.Id), 10),
			Name: artist.Name,
		})
	}

	return model.Song{
		Id:              strconv.FormatUint(uint64(songInfo.Data.Id), 10),
		Title:           songInfo.Data.Title,
		Duration:        songInfo.Data.Duration,
		ReplayGain:      downloadInfo.Data.TrackReplayGain,
		Peak:            downloadInfo.Data.TrackPeakAmplitude,
		AlbumReplayGain: downloadInfo.Data.AlbumReplayGain,
		AlbumPeak:       downloadInfo.Data.AlbumPeakAmplitude,
		ReleaseDate:     releaseDate,
		TrackNumber:     songInfo.Data.TrackNumber,
		VolumeNumber:    songInfo.Data.VolumeNumber,
		AudioQuality:    audioQuality,
		Explicit:        songInfo.Data.Explicit,
		Popularity:      songInfo.Data.Popularity,
		Isrc:            songInfo.Data.Isrc,
		Artists:         artists,
		Album: model.Album{
			Id:       strconv.FormatUint(uint64(songInfo.Data.Album.Id), 10),
			Title:    songInfo.Data.Album.Title,
			CoverUrl: hifi_utils.GetImageURL(songInfo.Data.Album.CoverUrl, 1280),
		},
	}, nil
}
