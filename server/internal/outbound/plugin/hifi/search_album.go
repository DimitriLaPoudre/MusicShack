package hifi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getSearchAlbum(ctx context.Context, urls []string, q string, limit int, offset int) (searchAlbumResponse, error) {
	searchAlbums, err := hifi_utils.MultiFetchTyped[searchAlbumResponse](ctx, urls, fmt.Sprintf("/search/?al=%s&limit=%d&offset=%d", url.QueryEscape(q), limit, offset), &p.limiters)
	if err != nil {
		return searchAlbumResponse{}, fmt.Errorf("fetch search album with url list: %w", err)
	}

	return searchAlbums, nil
}

func (p *Hifi) SearchAlbum(ctx context.Context, instances []model.Instance, q string, limit int, offset int) (model.PaginatedAlbums, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	searchAlbums, err := p.getSearchAlbum(ctx, urls, q, limit, offset)
	if err != nil {
		return model.PaginatedAlbums{}, err
	}

	albums := []model.AlbumInfo{}
	for _, album := range searchAlbums.Data.Albums.Albums {
		albums = append(albums, album.ToAlbumInfo(640))
	}

	return model.PaginatedAlbums{
		Pagination: model.Pagination{
			Limit:              limit,
			Offset:             offset,
			TotalNumberOfItems: searchAlbums.Data.Albums.TotalNumberOfItems},
		Albums: albums,
	}, nil
}
