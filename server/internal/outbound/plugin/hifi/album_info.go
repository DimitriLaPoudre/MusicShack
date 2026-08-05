package hifi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getAlbumInfo(ctx context.Context, urls []string, id string) (albumResponse, error) {
	album, err := hifi_utils.MultiFetchTyped[albumResponse](ctx, urls, "/album/?id="+url.QueryEscape(id), &p.limiters)
	if err != nil {
		return albumResponse{}, fmt.Errorf("fetch album info with url list: %w", err)
	}

	return album, nil
}

func (p *Hifi) AlbumInfo(ctx context.Context, instances []model.Instance, id string) (model.AlbumInfo, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	album, err := p.getAlbumInfo(ctx, urls, id)
	if err != nil {
		return model.AlbumInfo{}, err
	}

	return album.Data.ToAlbumInfo(640), nil
}
