package hifi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getPlaylistInfo(ctx context.Context, urls []string, id string) (playlistResponse, error) {
	playlistInfo, err := hifi_utils.MultiFetchTyped[playlistResponse](ctx, urls, "/playlist/?id="+url.QueryEscape(id), &p.limiters)
	if err != nil {
		return playlistResponse{}, fmt.Errorf("fetch playlist info with url list: %w", err)
	}

	return playlistInfo, nil
}

func (p *Hifi) PlaylistInfo(ctx context.Context, instances []model.Instance, id string) (model.PlaylistInfo, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	playlistInfo, err := p.getPlaylistInfo(ctx, urls, id)
	if err != nil {
		return model.PlaylistInfo{}, err
	}

	return playlistInfo.Playlist.ToPlaylistInfo(640), nil
}
