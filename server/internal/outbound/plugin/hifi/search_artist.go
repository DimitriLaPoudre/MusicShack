package hifi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getSearchArtist(ctx context.Context, urls []string, q string, limit int, offset int) (searchArtistResponse, error) {
	searchArtists, err := hifi_utils.MultiFetchTyped[searchArtistResponse](ctx, urls, fmt.Sprintf("/search/?a=%s&limit=%d&offset=%d", url.QueryEscape(q), limit, offset), &p.limiters)
	if err != nil {
		return searchArtistResponse{}, fmt.Errorf("fetch search artist with url list: %w", err)
	}

	return searchArtists, nil
}

func (p *Hifi) SearchArtist(ctx context.Context, instances []model.Instance, q string, limit int, offset int) (model.PaginatedArtists, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	searchArtists, err := p.getSearchArtist(ctx, urls, q, limit, offset)
	if err != nil {
		return model.PaginatedArtists{}, err
	}

	artists := []model.ArtistInfo{}
	for _, artist := range searchArtists.Data.Artists.Artists {
		artists = append(artists, artist.ToArtistInfo(750))
	}

	return model.PaginatedArtists{
		Pagination: model.Pagination{
			Limit:              limit,
			Offset:             offset,
			TotalNumberOfItems: searchArtists.Data.Artists.TotalNumberOfItems},
		Artists: artists,
	}, nil
}
