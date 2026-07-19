package hifi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

func checkSong(arr []string) (model.UrlItem, error) {
	if len(arr) != 2 {
		return model.UrlItem{}, fmt.Errorf("/track without id path")
	}

	return model.UrlItem{
		Type: model.TypeSong,
		Id:   arr[1],
	}, nil
}

func checkAlbum(arr []string) (model.UrlItem, error) {
	if len(arr) == 2 {
		return model.UrlItem{
			Type: model.TypeAlbum,
			Id:   arr[1],
		}, nil
	} else if len(arr) == 4 {
		if item, err := checkSong(arr[2:]); err != nil {
			return model.UrlItem{}, fmt.Errorf("/album: %w", err)
		} else {
			return item, nil
		}
	} else {
		return model.UrlItem{}, fmt.Errorf("/album without id path or id/track/id path")
	}
}

func checkArtist(arr []string) (model.UrlItem, error) {
	if len(arr) != 2 {
		return model.UrlItem{}, fmt.Errorf("/artist without id path")
	}

	return model.UrlItem{
		Type: model.TypeArtist,
		Id:   arr[1],
	}, nil
}

func checkPlaylist(arr []string) (model.UrlItem, error) {
	if len(arr) != 1 {
		return model.UrlItem{}, fmt.Errorf("/playlist without id path")
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
	if len(arr) == 1 {
		return model.UrlItem{}, fmt.Errorf("url contain 1 sub path: %v", arr)
	}

	var item model.UrlItem
	var err error
	switch arr[0] {
	case "track":
		item, err = checkSong(arr)
	case "album":
		item, err = checkAlbum(arr)
	case "artist":
		item, err = checkArtist(arr)
	case "playlist":
		item, err = checkPlaylist(arr[1:])
	default:
		item, err = model.UrlItem{}, fmt.Errorf("url contain unknown sub path: %v", arr)
	}

	if err != nil {
		return model.UrlItem{}, fmt.Errorf("url %s: %w", url, err)
	}

	return item, nil
}
