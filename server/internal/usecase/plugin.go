package usecase

import (
	"context"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
)

type PluginUseCase struct {
	plugin *service.PluginService
}

func NewPluginUseCase(plugin *service.PluginService) PluginUseCase {
	return PluginUseCase{
		plugin: plugin,
	}
}

func (u *PluginUseCase) GetSong(ctx context.Context, user model.User, provider string, id string) (model.EnrichedSong, error) {
	song, err := u.plugin.GetSong(ctx, user, provider, id)
	if err != nil {
		return model.EnrichedSong{}, err
	}

	return song, nil
}

func (u *PluginUseCase) GetAlbum(ctx context.Context, user model.User, provider string, id string) (model.EnrichedAlbum, error) {
	album, err := u.plugin.GetAlbum(ctx, user, provider, id)
	if err != nil {
		return model.EnrichedAlbum{}, err
	}

	return album, nil
}

func (u *PluginUseCase) GetArtist(ctx context.Context, user model.User, provider string, id string) (model.EnrichedArtist, error) {
	artist, err := u.plugin.GetArtist(ctx, user, provider, id)
	if err != nil {
		return model.EnrichedArtist{}, err
	}

	return artist, nil
}

func (u *PluginUseCase) GetPlaylist(ctx context.Context, user model.User, provider string, id string) (model.EnrichedPlaylist, error) {
	playlist, err := u.plugin.GetPlaylist(ctx, user, provider, id)
	if err != nil {
		return model.EnrichedPlaylist{}, err
	}

	return playlist, nil
}

func (u *PluginUseCase) GetSearch(ctx context.Context, user model.User, q string) (model.SearchResult, error) {
	searchResult, err := u.plugin.Search(ctx, user, q)
	if err != nil {
		return model.SearchResult{}, err
	}

	return searchResult, nil
}
