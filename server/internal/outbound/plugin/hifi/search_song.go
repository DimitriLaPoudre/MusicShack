package hifi

import (
	"context"
	"fmt"
	"net/url"

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
		songs = append(songs, song.ToSongInfo(1280))
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
