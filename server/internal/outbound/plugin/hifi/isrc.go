package hifi

import (
	"context"
	"fmt"
	"net/url"
	"slices"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getSearchISRC(ctx context.Context, urls []string, isrc string) (searchSongResponse, error) {
	searchSong, err := hifi_utils.MultiFetchTyped[searchSongResponse](ctx, urls, "/search/?i="+url.QueryEscape(isrc), &p.limiters)
	if err != nil {
		return searchSongResponse{}, fmt.Errorf("fetch search song with ISRC %s with url list: %w", isrc, err)
	}

	return searchSong, nil
}

func (p *Hifi) SongInfoByISRC(ctx context.Context, instances []model.Instance, isrc string) (model.SongInfo, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	songData, err := p.getSearchISRC(ctx, urls, isrc)
	if err != nil {
		return model.SongInfo{}, err
	}

	if len(songData.Data.Songs) == 0 {
		return model.SongInfo{}, model.ErrPluginDataNotFound
	}

	songs := []model.SongInfo{}
	for _, song := range songData.Data.Songs {
		songs = append(songs, song.ToSongInfo(1280))
	}

	slices.SortFunc(songs, func(a, b model.SongInfo) int {
		if a.ReleaseDate.After(b.ReleaseDate) {
			return -1
		}
		if a.ReleaseDate.Before(b.ReleaseDate) {
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

	return songs[0], nil
}
