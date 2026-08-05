package hifi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getSearchPlaylist(ctx context.Context, urls []string, q string, limit int, offset int) (searchPlaylistResponse, error) {
	searchPlaylists, err := hifi_utils.MultiFetchTyped[searchPlaylistResponse](ctx, urls, fmt.Sprintf("/search/?p=%s&limit=%d&offset=%d", url.QueryEscape(q), limit, offset), &p.limiters)
	if err != nil {
		return searchPlaylistResponse{}, fmt.Errorf("fetch search playlist with url list: %w", err)
	}

	return searchPlaylists, nil
}

func (p *Hifi) SearchPlaylist(ctx context.Context, instances []model.Instance, q string, limit int, offset int) (model.PaginatedPlaylists, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	searchPlaylists, err := p.getSearchPlaylist(ctx, urls, q, limit, offset)
	if err != nil {
		return model.PaginatedPlaylists{}, err
	}

	playlists := []model.PlaylistInfo{}
	for _, playlist := range searchPlaylists.Data.Playlists.Playlists {
		playlists = append(playlists, playlist.ToPlaylistInfo(640))
	}

	return model.PaginatedPlaylists{
		Pagination: model.Pagination{
			Limit:              limit,
			Offset:             offset,
			TotalNumberOfItems: searchPlaylists.Data.Playlists.TotalNumberOfItems},
		Playlists: playlists,
	}, nil
}
