package plugins

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

type plugStore struct {
	name     map[string]models.Plugin
	provider map[string][]models.Plugin
}

var store = plugStore{
	name:     make(map[string]models.Plugin),
	provider: make(map[string][]models.Plugin),
}

func Register(p models.Plugin) {
	store.name[p.Name()] = p
	store.provider[p.Provider()] = append(store.provider[p.Provider()], p)
	slices.SortFunc(store.provider[p.Provider()], func(a, b models.Plugin) int {
		return b.Priority() - a.Priority()
	})
}

func GetPluginByName(name string) (models.Plugin, bool) {
	p, ok := store.name[name]
	return p, ok
}

func GetAllPluginsByName() map[string]models.Plugin {
	return store.name
}

func GetPluginByProvider(provider string) ([]models.Plugin, bool) {
	p, ok := store.provider[provider]
	return p, ok
}

func GetAllPluginsByProvider() map[string][]models.Plugin {
	return store.provider
}

func GetSong(ctx context.Context, userId uint, provider string, id string) (models.SongData, error) {
	plugins, ok := GetPluginByProvider(provider)
	if !ok {
		return models.SongData{}, fmt.Errorf("services.GetSong: %w", errors.New("invalid provider name"))
	}

	var data models.SongData
	var err error
	for _, plugin := range plugins {
		data, err = plugin.Song(ctx, userId, id)
		if err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return models.SongData{}, err
	} else {
		return data, nil
	}
}

func GetPlaylist(ctx context.Context, userId uint, provider string, id string) (models.PlaylistData, error) {
	plugins, ok := GetPluginByProvider(provider)
	if !ok {
		return models.PlaylistData{}, fmt.Errorf("services.GetPlaylist: %w", errors.New("invalid provider name"))
	}

	var data models.PlaylistData
	var err error
	for _, plugin := range plugins {
		data, err = plugin.Playlist(ctx, userId, id)
		if err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return models.PlaylistData{}, err
	} else {
		return data, nil
	}
}

func GetAlbum(ctx context.Context, userId uint, provider string, id string) (models.AlbumData, error) {
	plugins, ok := GetPluginByProvider(provider)
	if !ok {
		return models.AlbumData{}, fmt.Errorf("services.GetAlbum: %w", errors.New("invalid provider name"))
	}

	var data models.AlbumData
	var err error
	for _, plugin := range plugins {
		data, err = plugin.Album(ctx, userId, id)
		if err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return models.AlbumData{}, err
	} else {
		return data, nil
	}
}

func GetArtist(ctx context.Context, userId uint, provider string, id string) (models.ArtistData, error) {
	plugins, ok := GetPluginByProvider(provider)
	if !ok {
		return models.ArtistData{}, fmt.Errorf("services.GetArtist: %w", errors.New("invalid provider name"))
	}

	var data models.ArtistData
	var err error
	for _, plugin := range plugins {
		data, err = plugin.Artist(ctx, userId, id)
		if err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return models.ArtistData{}, err
	} else {
		return data, nil
	}
}

func Download(ctx context.Context, userId uint, provider string, id string) (io.ReadCloser, string, error) {
	plugins, ok := GetPluginByProvider(provider)
	if !ok {
		return nil, "", fmt.Errorf("services.Download: %w", errors.New("invalid provider name"))
	}

	var reader io.ReadCloser
	var extension string
	var err error
	for _, plugin := range plugins {
		reader, extension, err = plugin.Download(ctx, userId, id)
		if err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return nil, "", err
	} else {
		return reader, extension, nil
	}
}
