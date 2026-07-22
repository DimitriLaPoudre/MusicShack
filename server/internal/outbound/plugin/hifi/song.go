package hifi

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func getSongInfo(ctx context.Context, limiters *sync.Map, urls []string, id string) (songResponse, error) {
	songInfo, err := hifi_utils.FetchTypeSequential[songResponse](ctx, urls, "/info/?id="+url.QueryEscape(id), limiters)
	if err != nil {
		return songResponse{}, fmt.Errorf("fetch song info with url list: %w", err)
	}

	return songInfo, nil
}

func getSongDownloadInfo(ctx context.Context, limiters *sync.Map, urls []string, id string) (downloadResponse, error) {
	downloadInfo, err := hifi_utils.FetchTypeSequential[downloadResponse](ctx, urls, "/track/?id="+url.QueryEscape(id), limiters)
	if err != nil {
		return downloadResponse{}, fmt.Errorf("fetch song download info with url list: %w", err)
	}

	return downloadInfo, nil
}

func getSong(ctx context.Context, limiters *sync.Map, urls []string, id string) (songResponse, downloadResponse, error) {
	var songInfo songResponse
	var songInfoErr error
	var downloadInfo downloadResponse
	// var downloadInfoErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		songInfo, songInfoErr = getSongInfo(ctx, limiters, urls, id)
	})
	wg.Go(func() {
		downloadInfo, _ = getSongDownloadInfo(ctx, limiters, urls, id)
	})
	wg.Wait()

	if songInfoErr != nil {
		return songResponse{}, downloadResponse{}, songInfoErr
	}

	return songInfo, downloadInfo, nil
}

func (p *Hifi) Song(ctx context.Context, instances []model.Instance, id string) (model.Song, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	songInfo, downloadInfo, err := getSong(ctx, &p.limiters, urls, id)
	if err != nil {
		return model.Song{}, err
	}

	releaseDate, err := time.Parse(StreamStartDateLayout, songInfo.Data.ReleaseDate)
	if err != nil {
		slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse song releaseDate %s", songInfo.Data.ReleaseDate), slog.String("err", err.Error()))
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
		Peak:            songInfo.Data.Peak,
		AlbumReplayGain: songInfo.Data.ReplayGain,
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
