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

func (p *Hifi) getAlbumInfo(ctx context.Context, urls []string, id string) (albumResponse, error) {
	album, err := hifi_utils.MultiFetchTyped[albumResponse](ctx, urls, "/album/?id="+url.QueryEscape(id), &p.limiters)
	if err != nil {
		return albumResponse{}, fmt.Errorf("fetch album info with url list: %w", err)
	}

	return album, nil
}

func (p *Hifi) AlbumInfo(ctx context.Context, instances []model.Instance, id string) (model.AlbumInfo, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	album, err := p.getAlbumInfo(ctx, urls, id)
	if err != nil {
		return model.AlbumInfo{}, err
	}
	releaseDate, err := time.Parse(ReleaseDateLayout, album.Data.ReleaseDate)
	if err != nil {
		slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse album releaseDate %s", album.Data.ReleaseDate), slog.String("err", err.Error()))
	}

	audioQuality := LOW
	switch album.Data.AudioQuality {
	case "LOW":
		audioQuality = LOW
	case "HIGH":
		audioQuality = HIGH
	case "LOSSLESS":
		audioQuality = LOSSLESS
	default:
		audioQuality = LOW
	}
	for _, quality := range album.Data.MediaMetadata.Tags {
		switch quality {
		case "HIRES_LOSSLESS":
			audioQuality = HIRES
		case "LOSSLESS", "DOLBY_ATMOS":
			if audioQuality != HIRES {
				audioQuality = LOSSLESS
			}
		}
	}

	artists := []model.ArtistInfo{}
	for _, artist := range album.Data.Artists {
		artists = append(artists, model.ArtistInfo{
			ID:   strconv.FormatUint(uint64(artist.ID), 10),
			Name: artist.Name,
		})
	}

	return model.AlbumInfo{
		ID:            strconv.FormatUint(uint64(album.Data.ID), 10),
		Title:         album.Data.Title,
		Duration:      album.Data.Duration,
		ReleaseDate:   releaseDate,
		NumberTracks:  album.Data.NumberOfTracks,
		NumberVolumes: album.Data.NumberOfVolumes,
		CoverUrl:      hifi_utils.GetImageURL(album.Data.CoverUrl, 640),
		AudioQuality:  audioQuality,
		Explicit:      album.Data.Explicit,
		Artists:       artists,
	}, nil
}
