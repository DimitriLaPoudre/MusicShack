package hifi

import (
	"context"
	"fmt"
	"log/slog"
	lib_url "net/url"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func getSearchISRC(ctx context.Context, limiters *sync.Map, urls []string, isrc string) (searchSongResponse, error) {
	searchSong, err := hifi_utils.FetchTypeSequential[searchSongResponse](ctx, urls, "/search/?i="+lib_url.QueryEscape(isrc), limiters)
	if err != nil {
		return searchSongResponse{}, fmt.Errorf("fetch search song with ISRC %s with url list: %w", isrc, err)
	}

	return searchSong, nil
}

func (p *Hifi) SongByISRC(ctx context.Context, instances []model.Instance, isrc string) (model.Song, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	songData, err := getSearchISRC(ctx, &p.limiters, urls, isrc)
	if err != nil {
		return model.Song{}, err
	}

	if len(songData.Data.Songs) == 0 {
		return model.Song{}, model.ErrPluginDataNotFound
	}

	songs := []model.Song{}
	for _, song := range songData.Data.Songs {
		releaseDate, err := time.Parse(StreamStartDateLayout, song.ReleaseDate)
		if err != nil {
			slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse song with ISRC %s releaseDate %s", isrc, song.ReleaseDate), slog.String("err", err.Error()))
		}

		audioQuality := LOW
		switch song.AudioQuality {
		case "LOW":
			audioQuality = LOW
		case "HIGH":
			audioQuality = HIGH
		case "LOSSLESS":
			audioQuality = LOSSLESS
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

		songs = append(songs,
			model.Song{
				Id:           strconv.FormatUint(uint64(song.Id), 10),
				Title:        song.Title,
				Duration:     song.Duration,
				AudioQuality: audioQuality,
				Popularity:   song.Popularity,
				Explicit:     song.Explicit,
				Isrc:         song.Isrc,
				Artists:      artists,
				Album: model.Album{
					Id:       strconv.FormatUint(uint64(song.Album.Id), 10),
					Title:    song.Album.Title,
					CoverUrl: hifi_utils.GetImageURL(song.Album.CoverUrl, 640),
				},
			})
		songs = append(songs,
			model.Song{
				Id:           strconv.FormatUint(uint64(song.Id), 10),
				Title:        song.Title,
				Duration:     song.Duration,
				ReleaseDate:  releaseDate,
				TrackNumber:  song.TrackNumber,
				VolumeNumber: song.VolumeNumber,
				AudioQuality: audioQuality,
				Explicit:     song.Explicit,
				Popularity:   song.Popularity,
				Isrc:         song.Isrc,
				Artists:      artists,
				Album: model.Album{
					Id:       strconv.FormatUint(uint64(song.Album.Id), 10),
					Title:    song.Album.Title,
					CoverUrl: hifi_utils.GetImageURL(song.Album.CoverUrl, 1280),
				},
			})
	}

	slices.SortFunc(songs, func(a, b model.Song) int {
		if a.ReleaseDate.After(b.ReleaseDate) {
			return -1
		}
		if a.ReleaseDate.Before(b.ReleaseDate) {
			return 1
		}
		if a.Id > b.Id {
			return -1
		}
		if a.Id < b.Id {
			return 1
		}
		return 0
	})

	return songs[0], nil
}
