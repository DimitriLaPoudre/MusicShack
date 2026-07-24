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

func (p *Hifi) getSongInfo(ctx context.Context, urls []string, id string) (songResponse, error) {
	songInfo, err := hifi_utils.CachedMultiFetchTyped[songResponse](ctx, urls, "/info/?id="+url.QueryEscape(id), &p.limiters, &p.cache)
	if err != nil {
		return songResponse{}, fmt.Errorf("fetch song info with url list: %w", err)
	}

	return songInfo, nil
}

func (p *Hifi) getSongDownloadInfo(ctx context.Context, urls []string, id string) (downloadResponse, error) {
	downloadInfo, err := hifi_utils.CachedMultiFetchTyped[downloadResponse](ctx, urls, "/track/?id="+url.QueryEscape(id), &p.limiters, &p.cache)
	if err != nil {
		return downloadResponse{}, fmt.Errorf("fetch song download info with url list: %w", err)
	}

	return downloadInfo, nil
}

func (p *Hifi) getSong(ctx context.Context, urls []string, id string) (songResponse, downloadResponse, error) {
	var songInfo songResponse
	var songInfoErr error
	var downloadInfo downloadResponse
	// var downloadInfoErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		songInfo, songInfoErr = p.getSongInfo(ctx, urls, id)
	})
	wg.Go(func() {
		downloadInfo, _ = p.getSongDownloadInfo(ctx, urls, id)
	})
	wg.Wait()

	if songInfoErr != nil {
		return songResponse{}, downloadResponse{}, songInfoErr
	}

	return songInfo, downloadInfo, nil
}

func (p *Hifi) Song(ctx context.Context, instances []model.Instance, id string) (model.Song, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	songInfo, downloadInfo, err := p.getSong(ctx, urls, id)
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
			ID:   strconv.FormatUint(uint64(artist.ID), 10),
			Name: artist.Name,
		})
	}

	return model.Song{
		ID:              strconv.FormatUint(uint64(songInfo.Data.ID), 10),
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
			ID:       strconv.FormatUint(uint64(songInfo.Data.Album.ID), 10),
			Title:    songInfo.Data.Album.Title,
			CoverUrl: hifi_utils.GetImageURL(songInfo.Data.Album.CoverUrl, 1280),
		},
	}, nil
}
