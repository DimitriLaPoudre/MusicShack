package hifi

import (
	"context"
	"fmt"
	lib_url "net/url"
	"strconv"
	"sync"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	hifi_utils "github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/plugin/hifi/utils"
	"golang.org/x/time/rate"
)

func getSongInfo(ctx context.Context, limiters map[string]*rate.Limiter, urls []string, id string) (songResponse, error) {
	songInfo, err := hifi_utils.FetchTypeSequential[songResponse](ctx, urls, "/info/?id="+lib_url.QueryEscape(id), limiters)
	if err != nil {
		return songResponse{}, fmt.Errorf("getSongInfo: %w", err)
	}

	return songInfo, nil
}

func getSong(ctx context.Context, limiters map[string]*rate.Limiter, urls []string, id string) (songResponse, downloadResponse, error) {
	var songInfo songResponse
	var songInfoErr error
	var downloadInfo downloadResponse
	var downloadInfoErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		songInfo, songInfoErr = getSongInfo(ctx, limiters, urls, id)
	})
	wg.Go(func() {
		// downloadInfo, downloadInfoErr := getDownloadInfo(ctx, limiters, urls, id, "")
	})
	wg.Wait()

	if songInfoErr != nil {
		return songResponse{}, downloadResponse{}, fmt.Errorf("getSong: %w", songInfoErr)
	}
	if downloadInfoErr != nil {
		return songResponse{}, downloadResponse{}, fmt.Errorf("getSong: %w", downloadInfoErr)
	}

	return songInfo, downloadInfo, nil
}

func (p *Hifi) Song(ctx context.Context, instances []model.Instance, id string) (model.Song, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	songInfo, downloadInfo, err := getSong(ctx, p.limiters, urls, id)
	if err != nil {
		return model.Song{}, fmt.Errorf("Hifi.Song: %w", err)
	}

	releaseDate, err := time.Parse(StreamStartDateLayout, songInfo.Data.ReleaseDate)
	if err != nil {
		p.l.Warn().Msg(fmt.Sprintf("Hifi.Song: time.Parse(%s): %s", songInfo.Data.ReleaseDate, err.Error()))
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
