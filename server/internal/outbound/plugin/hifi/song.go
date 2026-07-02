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
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/network"
)

func fetchSong(ctx context.Context, url string, id string) (songData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := network.Fetch(ctx, url+"/info/?id="+lib_url.QueryEscape(id), nil)
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

func getSongData(ctx context.Context, url string, id string) (songData, downloadData, error) {
	var songInfo songData
	var songInfoErr error
	var downloadInfo downloadData
	var downloadInfoErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		songInfo, songInfoErr = fetchSong(ctx, url, id)
	}()
	go func() {
		defer wg.Done()
		// downloadInfo, downloadInfoErr := getDownloadInfo(ctx, url, id, "")
	}()
	wg.Wait()
	err := songInfoErr
	if err == nil {
		err = downloadInfoErr
	}

	if err != nil {
		return songData{}, downloadData{}, fmt.Errorf("getSongData: %w", err)
	}

	return songInfo, downloadInfo, nil
}

func (p *Hifi) Song(ctx context.Context, url string, id string) (model.Song, error) {
	data, downloadInfo, err := getSongData(ctx, url, id)
	if err != nil {
		return model.Song{}, fmt.Errorf("Hifi.Song: %w", err)
	}

	normalizeSongData := model.Song{
		Provider:        p.Provider(),
		Api:             p.Name(),
		Id:              strconv.FormatUint(uint64(data.Data.Id), 10),
		Title:           data.Data.Title,
		Duration:        data.Data.Duration,
		ReplayGain:      downloadInfo.Data.TrackReplayGain,
		Peak:            downloadInfo.Data.TrackPeakAmplitude,
		AlbumReplayGain: downloadInfo.Data.AlbumReplayGain,
		AlbumPeak:       downloadInfo.Data.AlbumPeakAmplitude,
		ReleaseDate:     data.Data.ReleaseDate[:10],
		TrackNumber:     data.Data.TrackNumber,
		VolumeNumber:    data.Data.VolumeNumber,
		Explicit:        data.Data.Explicit,
		Popularity:      data.Data.Popularity,
		Isrc:            data.Data.Isrc,
		Artists:         make([]model.SongArtist, 0),
		Album: model.SongAlbum{
			Id:       strconv.FormatUint(uint64(data.Data.Album.Id), 10),
			Title:    data.Data.Album.Title,
			CoverUrl: hifi_utils.GetImageURL(data.Data.Album.CoverUrl, 1280),
		},
	}

	switch data.Data.AudioQuality {
	case "LOW":
		normalizeSongData.AudioQuality = LOW
	case "HIGH":
		normalizeSongData.AudioQuality = HIGH
	case "LOSSLESS":
		normalizeSongData.AudioQuality = LOSSLESS
	}
	for _, quality := range data.Data.MediaMetadata.Tags {
		switch quality {
		case "HIRES_LOSSLESS":
			normalizeSongData.AudioQuality = HIRES
		case "LOSSLESS", "DOLBY_ATMOS":
			if normalizeSongData.AudioQuality != HIRES {
				normalizeSongData.AudioQuality = LOSSLESS
			}
		}
	}

	for _, artist := range data.Data.Artists {
		normalizeSongData.Artists = append(normalizeSongData.Artists, model.SongArtist{
			Id:   strconv.FormatUint(uint64(artist.Id), 10),
			Name: artist.Name,
		})
	}

	return normalizeSongData, nil
}
