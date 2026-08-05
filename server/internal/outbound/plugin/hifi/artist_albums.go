package hifi

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/dto"
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

	best := make(map[albumItemComparaison]dto.AlbumItem)
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

	// extension := make(map[albumItemComparaisonExtension][]*dto.AlbumItem)
	// for _, album := range best {
	// 	newExtension := []*dto.AlbumItem{}
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

	list := []dto.AlbumItem{}
	for _, album := range best {
		list = append(list, album)
	}

	slices.SortFunc(list, func(a, b dto.AlbumItem) int {
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
		newAlbum := album.ToAlbumInfo(1280)

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
