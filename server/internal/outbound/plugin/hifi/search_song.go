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

func (p *Hifi) getSearchSong(ctx context.Context, urls []string, q string, limit int, offset int) (searchSongResponse, error) {
	searchSongs, err := hifi_utils.MultiFetchTyped[searchSongResponse](ctx, urls, fmt.Sprintf("/search/?s=%s&limit=%d&offset=%d", url.QueryEscape(q), limit, offset), &p.limiters)
	if err != nil {
		return searchSongResponse{}, fmt.Errorf("fetch search song with url list: %w", err)
	}

	return searchSongs, nil
}

func (p *Hifi) SearchSong(ctx context.Context, instances []model.Instance, q string, limit int, offset int) (model.PaginatedSongs, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	searchSongs, err := p.getSearchSong(ctx, urls, q, limit, offset)
	if err != nil {
		return model.PaginatedSongs{}, err
	}

	songs := []model.SongInfo{}
	for _, song := range searchSongs.Data.Songs {
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
				CoverURL: hifi_utils.GetImageURL(song.Album.Cover, 1280),
			},
		})
	}

	return model.PaginatedSongs{
		Pagination: model.Pagination{
			Limit:              limit,
			Offset:             offset,
			TotalNumberOfItems: searchSongs.Data.TotalNumberOfItems,
		},
		Songs: songs,
	}, nil
}
