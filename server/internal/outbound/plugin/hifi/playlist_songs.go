package hifi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getPlaylistSongs(ctx context.Context, urls []string, id string, limit int, offset int) (playlistResponse, error) {
	playlistSongs, err := hifi_utils.MultiFetchTyped[playlistResponse](ctx, urls, fmt.Sprintf("/playlist/?id=%s&limit=%d&offset=%d", url.QueryEscape(id), limit, offset), &p.limiters)
	if err != nil {
		return playlistResponse{}, fmt.Errorf("fetch playlist songs with url list: %w", err)
	}

	return playlistSongs, nil
}

func (p *Hifi) PlaylistSongs(ctx context.Context, instances []model.Instance, id string, limit int, offset int) (model.PaginatedSongs, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	playlistSongs, err := p.getPlaylistSongs(ctx, urls, id, limit, offset)
	if err != nil {
		return model.PaginatedSongs{}, err
	}

	songs := []model.SongInfo{}
	for _, item := range playlistSongs.Items {
		if item.Type != "track" {
			continue
		}
		songs = append(songs, item.Item.ToSongInfo(1280))
	}

	return model.PaginatedSongs{
		Pagination: model.Pagination{
			Limit:              limit,
			Offset:             offset,
			TotalNumberOfItems: playlistSongs.Playlist.NumberOfTracks,
		},
		Songs: songs,
	}, nil
}
