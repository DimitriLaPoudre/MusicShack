package hifi

import (
	"context"
	"fmt"
	"net/url"

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

func (p *Hifi) AlbumSongs(ctx context.Context, instances []model.Instance, id string, limit int, offset int) (model.PaginatedSongs, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	album, err := p.getAlbumSongs(ctx, urls, id, limit, offset)
	if err != nil {
		return model.PaginatedSongs{}, err
	}

	songs := []model.SongInfo{}
	for _, item := range album.Data.Items {
		songs = append(songs, item.Item.ToSongInfo(1280))
	}

	return model.PaginatedSongs{
		Pagination: model.Pagination{
			Limit:              limit,
			Offset:             offset,
			TotalNumberOfItems: album.Data.NumberOfTracks,
		},
		Songs: songs,
	}, nil
}
