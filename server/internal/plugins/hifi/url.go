package hifi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func (p Hifi) checkArtist(arr []string) (models.UrlItem, error) {
	if len(arr) != 2 {
		return models.UrlItem{}, errors.New(fmt.Sprint("Hifi.checkArtist: url contain artist sub path but without valid pattern id:", arr))
	}
	return models.UrlItem{
		Provider: p.Provider(),
		Type:     models.TypeArtist,
		Id:       arr[1],
	}, nil
}

func (p Hifi) checkAlbum(arr []string) (models.UrlItem, error) {
	if len(arr) == 2 {
		return models.UrlItem{
			Provider: p.Provider(),
			Type:     models.TypeAlbum,
			Id:       arr[1],
		}, nil
	} else if len(arr) == 4 {
		return p.checkSong(arr[2:])
	} else {
		return models.UrlItem{}, errors.New(fmt.Sprint("Hifi.checkAlbum: url contain album sub path but without valid pattern id or id/track/id:", arr))
	}
}

func (p Hifi) checkSong(arr []string) (models.UrlItem, error) {
	if len(arr) != 2 {
		return models.UrlItem{}, errors.New(fmt.Sprint("Hifi.checkSong: url contain track sub path but without valid pattern id:", arr))
	}
	return models.UrlItem{
		Provider: p.Provider(),
		Type:     models.TypeSong,
		Id:       arr[1],
	}, nil
}

func (p *Hifi) Url(ctx context.Context, userId uint, url string) (models.UrlItem, error) {
	var clean_url string
	if url, ok := strings.CutPrefix(url, "https://tidal.com/"); !ok {
		return models.UrlItem{}, errors.New("Hifi.Url: url not contain \"https://tidal.com/\"")
	} else {
		clean_url = url
	}

	arr := strings.Split(clean_url, "/")
	if len(arr) == 1 {
		return models.UrlItem{}, errors.New(fmt.Sprint("Hifi.Url: url contain 1 sub path:", arr))
	}
	switch arr[0] {
	case "artist":
		return p.checkArtist(arr)
	case "album":
		return p.checkAlbum(arr)
	case "track":
		return p.checkSong(arr)
	default:
		return models.UrlItem{}, errors.New(fmt.Sprint("Hifi.Url: url contain unknown sub path:", arr))
	}
}
