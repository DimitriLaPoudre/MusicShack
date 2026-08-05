package hifi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getArtistInfo(ctx context.Context, urls []string, id string) (artistResponse, error) {
	info, err := hifi_utils.MultiFetchTyped[artistResponse](ctx, urls, "/artist/?id="+url.QueryEscape(id), &p.limiters)
	if err != nil {
		return artistResponse{}, fmt.Errorf("fetch artist info with url list: %w", err)
	}

	return info, nil
}

func (p *Hifi) ArtistInfo(ctx context.Context, instances []model.Instance, id string) (model.ArtistInfo, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	artistInfo, err := p.getArtistInfo(ctx, urls, id)
	if err != nil {
		return model.ArtistInfo{}, err
	}

	return artistInfo.Artist.ToArtistInfo(750), nil
}
