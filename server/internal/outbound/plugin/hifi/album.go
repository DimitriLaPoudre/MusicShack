package hifi

import (
	"context"
	"fmt"
	lib_url "net/url"
	"strconv"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	hifi_utils "github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/plugin/hifi/utils"
	"golang.org/x/time/rate"
)

func getAlbum(ctx context.Context, limiters map[string]*rate.Limiter, urls []string, id string) (albumResponse, error) {
	album, err := hifi_utils.FetchTypeSequential[albumResponse](ctx, urls, "/album/?id="+lib_url.QueryEscape(id), limiters)
	if err != nil {
		return albumResponse{}, fmt.Errorf("getAlbum: %w", err)
	}

	return album, nil
}

func (p *Hifi) Album(ctx context.Context, instances []model.Instance, id string) (model.Album, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	album, err := getAlbum(ctx, p.limiters, urls, id)
	if err != nil {
		return model.Album{}, fmt.Errorf("Hifi.Album: %w", err)
	}
	releaseDate, err := time.Parse(StreamStartDateLayout, album.Data.ReleaseDate)
	if err != nil {
		p.l.Warn().Msg(fmt.Sprintf("Hifi.Album: time.Parse(%s): %s", album.Data.ReleaseDate, err.Error()))
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

	songs := []model.Song{}
	for _, item := range album.Data.Items {
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
				Id:   strconv.FormatUint(uint64(artist.Id), 10),
				Name: artist.Name,
			})
		}

		songs = append(songs, model.Song{
			Id:           strconv.FormatUint(uint64(song.Id), 10),
			Title:        song.Title,
			Duration:     song.Duration,
			TrackNumber:  song.TrackNumber,
			VolumeNumber: song.VolumeNumber,
			AudioQuality: audioQuality,
			Explicit:     song.Explicit,
			Isrc:         song.Isrc,
			Artists:      artists,
		})
	}

	artists := []model.Artist{}
	for _, artist := range album.Data.Artists {
		artists = append(artists, model.Artist{
			Id:   strconv.FormatUint(uint64(artist.Id), 10),
			Name: artist.Name,
		})
	}

	return model.Album{
		Id:            strconv.FormatUint(uint64(album.Data.Id), 10),
		Title:         album.Data.Title,
		Duration:      album.Data.Duration,
		ReleaseDate:   releaseDate,
		NumberTracks:  album.Data.NumberOfTracks,
		NumberVolumes: album.Data.NumberOfVolumes,
		CoverUrl:      hifi_utils.GetImageURL(album.Data.CoverUrl, 640),
		AudioQuality:  audioQuality,
		Explicit:      album.Data.Explicit,
		Songs:         songs,
		Artists:       artists,
	}, nil
}
