package hifi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

func checkSong(arr []string) (model.URLItem, error) {
	if len(arr) != 1 {
		return model.URLItem{}, fmt.Errorf("/track/id path but extra path found")
	}

	return model.URLItem{
		Type: model.TypeSong,
		ID:   arr[1],
	}, nil
}

func checkAlbum(arr []string) (model.URLItem, error) {
	if len(arr) != 1 {
		return model.URLItem{}, fmt.Errorf("/album/id path but extra path found")
	}
	return model.URLItem{
		Type: model.TypeAlbum,
		ID:   arr[1],
	}, nil
}

func checkArtist(arr []string) (model.URLItem, error) {
	if len(arr) != 1 {
		return model.URLItem{}, fmt.Errorf("/artist/id path but extra path found")
	}

	return model.URLItem{
		Type: model.TypeArtist,
		ID:   arr[1],
	}, nil
}

func checkPlaylist(arr []string) (model.URLItem, error) {
	if len(arr) != 1 {
		return model.URLItem{}, fmt.Errorf("/playlist/id path but extra path found")
	}

	return model.URLItem{
		Type: model.TypePlaylist,
		ID:   arr[0],
	}, nil
}

func (p *Hifi) URL(ctx context.Context, url string) (model.URLItem, error) {
	var clean_url string
	if url, ok := strings.CutPrefix(url, "https://tidal.com/"); !ok {
		return model.URLItem{}, errors.New("url not contain 'https://tidal.com/'")
	} else {
		clean_url = url
	}

	arr := strings.Split(clean_url, "/")
	if len(arr) == 0 {
		return model.URLItem{}, fmt.Errorf("url %s: path missing", url)
	}
	if arr[len(arr)-1] == "u" {
		arr = arr[:len(arr)-1]
	}
	if len(arr) <= 1 {
		return model.URLItem{}, fmt.Errorf("url %s: path missing /type/id pattern", url)
	}

	var item model.URLItem
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
		item, err = model.URLItem{}, fmt.Errorf("/type/id bad type: %v", arr)
	}

	if err != nil {
		return model.URLItem{}, fmt.Errorf("url %s: %w", url, err)
	}

	return item, nil
}
