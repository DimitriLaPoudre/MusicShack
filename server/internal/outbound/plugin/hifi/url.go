package hifi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

func checkSong(arr []string) (model.UrlItem, error) {
	if len(arr) != 1 {
		return model.UrlItem{}, fmt.Errorf("/track/id path but extra path found")
	}

	return model.UrlItem{
		Type: model.TypeSong,
		Id:   arr[1],
	}, nil
}

func checkAlbum(arr []string) (model.UrlItem, error) {
	if len(arr) != 1 {
		return model.UrlItem{}, fmt.Errorf("/album/id path but extra path found")
	}
	return model.UrlItem{
		Type: model.TypeAlbum,
		Id:   arr[1],
	}, nil
}

func checkArtist(arr []string) (model.UrlItem, error) {
	if len(arr) != 1 {
		return model.UrlItem{}, fmt.Errorf("/artist/id path but extra path found")
	}

	return model.UrlItem{
		Type: model.TypeArtist,
		Id:   arr[1],
	}, nil
}

func checkPlaylist(arr []string) (model.UrlItem, error) {
	if len(arr) != 1 {
		return model.UrlItem{}, fmt.Errorf("/playlist/id path but extra path found")
	}

	return model.UrlItem{
		Type: model.TypePlaylist,
		Id:   arr[0],
	}, nil
}

func (p *Hifi) Url(ctx context.Context, url string) (model.UrlItem, error) {
	var clean_url string
	if url, ok := strings.CutPrefix(url, "https://tidal.com/"); !ok {
		return model.UrlItem{}, errors.New("url not contain 'https://tidal.com/'")
	} else {
		clean_url = url
	}

	arr := strings.Split(clean_url, "/")
	if len(arr) == 0 {
		return model.UrlItem{}, fmt.Errorf("url %s: path missing", url)
	}
	if arr[len(arr)-1] == "u" {
		arr = arr[:len(arr)-1]
	}
	if len(arr) <= 1 {
		return model.UrlItem{}, fmt.Errorf("url %s: path missing /type/id pattern", url)
	}

	var item model.UrlItem
	var err error
	switch arr[0] {
	case "track":
		item, err = checkSong(arr[1:])
	case "album":
		item, err = checkAlbum(arr[1:])
	case "artist":
		item, err = checkArtist(arr[1:])
	case "playlist":
		item, err = checkPlaylist(arr[1:])
	default:
		item, err = model.UrlItem{}, fmt.Errorf("/type/id bad type: %v", arr)
	}

	if err != nil {
		return model.UrlItem{}, fmt.Errorf("url %s: %w", url, err)
	}

	return item, nil
}
