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

func (p *Hifi) getAlbumSongs(ctx context.Context, urls []string, id string, limit int, offset int) (albumResponse, error) {
	album, err := hifi_utils.MultiFetchTyped[albumResponse](ctx, urls, fmt.Sprintf("/album/?id=%s&limit=%d&offset=%d", url.QueryEscape(id), limit, offset), &p.limiters)
	if err != nil {
		return albumResponse{}, fmt.Errorf("fetch album songs with url list: %w", err)
	}

	return album, nil
}

func (p *Hifi) AlbumSongs(ctx context.Context, instances []model.Instance, id string, limit int, offset int) (model.AlbumSongs, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	album, err := p.getAlbumSongs(ctx, urls, id, limit, offset)
	if err != nil {
		return model.AlbumSongs{}, err
	}

	songs := []model.SongInfo{}
	for _, item := range album.Data.Items {
		song := item.Item

		releaseDate, err := time.Parse(StreamStartDateLayout, song.StreamStartDate)
		if err != nil {
			slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse album's song %s releaseDate %s", song.Title, song.StreamStartDate), slog.String("err", err.Error()))
		}

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

		artists := []model.ArtistInfo{}
		for _, artist := range song.Artists {
			artists = append(artists, model.ArtistInfo{
				ID:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		songs = append(songs, model.SongInfo{
			ID:              strconv.FormatUint(uint64(song.ID), 10),
			Title:           song.Title,
			Duration:        song.Duration,
			Peak:            song.Peak,
			AlbumReplayGain: song.ReplayGain,
			ReleaseDate:     releaseDate,
			TrackNumber:     song.TrackNumber,
			VolumeNumber:    song.VolumeNumber,
			AudioQuality:    audioQuality,
			Explicit:        song.Explicit,
			Popularity:      song.Popularity,
			Isrc:            song.ISRC,
			Artists:         artists,
			Album: model.AlbumInfo{
				ID:       strconv.FormatUint(uint64(song.Album.ID), 10),
				Title:    song.Album.Title,
				CoverUrl: hifi_utils.GetImageURL(song.Album.Cover, 1280),
			},
		})
	}

	return model.AlbumSongs{
		Songs: songs,
	}, nil
}
