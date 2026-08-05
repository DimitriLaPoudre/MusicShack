package hifi

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getArtistAlbums(ctx context.Context, urls []string, id string, limit int, offset int) (artistAlbumsResponse, error) {
	_, _ = limit, offset
	albums, err := hifi_utils.MultiFetchTyped[artistAlbumsResponse](ctx, urls, "/artist/?f="+url.QueryEscape(id)+"&skip_tracks=1", &p.limiters)
	if err != nil {
		return artistAlbumsResponse{}, fmt.Errorf("fetch artist albums with url list: %w", err)
	}

	return albums, nil
}

func (p *Hifi) ArtistAlbums(ctx context.Context, instances []model.Instance, id string, limit int, offset int) (model.ArtistPaginatedAlbums, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	artistAlbums, err := p.getArtistAlbums(ctx, urls, id, limit, offset)
	if err != nil {
		return model.ArtistPaginatedAlbums{}, err
	}

	type albumItemComparaison struct {
		Title       string
		ReleaseDate string
		TrackNumber int
	}

	best := make(map[albumItemComparaison]albumItem)
	for _, album := range artistAlbums.Albums.Items {
		if bestVersion, ok := best[albumItemComparaison{
			Title:       strings.ToLower(album.Title),
			ReleaseDate: album.ReleaseDate,
			TrackNumber: album.NumberOfTracks,
		}]; !ok || (!bestVersion.Explicit && album.Explicit) || (bestVersion.Explicit == album.Explicit && len(bestVersion.MediaMetadata.Tags) < len(album.MediaMetadata.Tags)) {
			best[albumItemComparaison{
				Title:       strings.ToLower(album.Title),
				ReleaseDate: album.ReleaseDate,
				TrackNumber: album.NumberOfTracks,
			}] = album
		}
	}

	// type albumItemComparaisonExtension struct {
	// 	Title       string
	// 	ReleaseDate string
	// }

	// extension := make(map[albumItemComparaisonExtension][]*albumItem)
	// for _, album := range best {
	// 	newExtension := []*albumItem{}
	// 	bestExtension, ok := extension[albumItemComparaisonExtension{
	// 		Title:       strings.ToLower(album.Title),
	// 		ReleaseDate: album.ReleaseDate,
	// 	}]
	// 	if !ok {
	// 		extension[albumItemComparaisonExtension{
	// 			Title:       strings.ToLower(album.Title),
	// 			ReleaseDate: album.ReleaseDate,
	// 		}] = append(newExtension, album)
	// 		continue
	// 	}
	//
	// 	for i, tmp := range bestExtension {
	//
	// 	}
	// }

	list := []albumItem{}
	for _, album := range best {
		list = append(list, album)
	}

	slices.SortFunc(list, func(a, b albumItem) int {
		if a.ReleaseDate > b.ReleaseDate {
			return -1
		}
		if a.ReleaseDate < b.ReleaseDate {
			return 1
		}
		if a.ID > b.ID {
			return -1
		}
		if a.ID < b.ID {
			return 1
		}
		return 0
	})

	albums := []model.AlbumInfo{}
	eps := []model.AlbumInfo{}
	singles := []model.AlbumInfo{}
	for _, album := range list {
		releaseDate, err := time.Parse(ReleaseDateLayout, album.ReleaseDate)
		if err != nil {
			slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse artist's album %s releaseDate %s", album.Title, album.ReleaseDate), slog.String("err", err.Error()))
		}

		audioQuality := LOW
		switch album.AudioQuality {
		case "LOW":
			audioQuality = LOW
		case "HIGH":
			audioQuality = HIGH
		case "LOSSLESS":
			audioQuality = LOSSLESS
		}
		for _, quality := range album.MediaMetadata.Tags {
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
		for _, artist := range album.Artists {
			artists = append(artists, model.ArtistInfo{
				ID:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		newAlbum := model.AlbumInfo{
			ID:           strconv.FormatUint(uint64(album.ID), 10),
			Title:        album.Title,
			Duration:     album.Duration,
			ReleaseDate:  releaseDate,
			CoverURL:     hifi_utils.GetImageURL(album.Cover, 1280),
			AudioQuality: audioQuality,
			Explicit:     album.Explicit,
			Artists:      artists,
		}

		switch album.Type {
		case "ALBUM":
			albums = append(albums, newAlbum)
		case "EP":
			eps = append(eps, newAlbum)
		case "SINGLE":
			singles = append(singles, newAlbum)
		}
	}

	return model.ArtistPaginatedAlbums{
		Albums:  model.PaginatedAlbums{Pagination: model.Pagination{Limit: limit, Offset: offset, TotalNumberOfItems: len(albums)}, Albums: paginate(albums, limit, offset)},
		EPs:     model.PaginatedAlbums{Pagination: model.Pagination{Limit: limit, Offset: offset, TotalNumberOfItems: len(eps)}, Albums: paginate(eps, limit, offset)},
		Singles: model.PaginatedAlbums{Pagination: model.Pagination{Limit: limit, Offset: offset, TotalNumberOfItems: len(singles)}, Albums: paginate(singles, limit, offset)},
	}, nil
}

func paginate[T any](items []T, limit, offset int) []T {
	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		return []T{}
	}

	if offset >= len(items) {
		return []T{}
	}

	end := offset + limit
	if end > len(items) {
		end = len(items)
	}

	return items[offset:end]
}
