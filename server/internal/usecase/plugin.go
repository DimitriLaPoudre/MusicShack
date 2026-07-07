package usecase

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/service"
	"github.com/rs/zerolog"
)

type PluginUseCase struct {
	l        *zerolog.Logger
	plugin   *service.PluginService
	store    *service.PluginStoreService
	instance model.InstanceRepository
}

func NewPluginUseCase(l *zerolog.Logger, plugin *service.PluginService, store *service.PluginStoreService, instance model.InstanceRepository) PluginUseCase {
	return PluginUseCase{
		l:        l,
		plugin:   plugin,
		store:    store,
		instance: instance,
	}
}

func (u *PluginUseCase) GetSong(ctx context.Context, user model.User, provider string, id string) (model.EnrichedSong, error) {
	instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.EnrichedSong{}, err
	}

	song, err := u.plugin.GetSong(ctx, instances, id)
	if err != nil {
		return model.EnrichedSong{}, err
	}

	return song, nil
}

func (u *PluginUseCase) GetAlbum(ctx context.Context, user model.User, provider string, id string) (model.EnrichedAlbum, error) {
	instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.EnrichedAlbum{}, err
	}

	album, err := u.plugin.GetAlbum(ctx, instances, id)
	if err != nil {
		return model.EnrichedAlbum{}, err
	}

	return album, nil
}

func (u *PluginUseCase) GetArtist(ctx context.Context, user model.User, provider string, id string) (model.EnrichedArtist, error) {
	instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.EnrichedArtist{}, err
	}

	artist, err := u.plugin.GetArtist(ctx, instances, id)
	if err != nil {
		return model.EnrichedArtist{}, err
	}

	return artist, nil
}

func (u *PluginUseCase) GetPlaylist(ctx context.Context, user model.User, provider string, id string) (model.EnrichedPlaylist, error) {
	instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.EnrichedPlaylist{}, err
	}

	playlist, err := u.plugin.GetPlaylist(ctx, instances, id)
	if err != nil {
		return model.EnrichedPlaylist{}, err
	}

	return playlist, nil
}

func (u *PluginUseCase) GetSearch(ctx context.Context, user model.User, q string) (map[string]model.EnrichedSearch, error) {
	providerResult := make(map[string]model.EnrichedSearch)

	for provider := range u.store.ListPluginsByProvider() {
		instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
		if err != nil {
			continue
		}

		result, err := u.plugin.Search(ctx, instances, q)
		if err != nil {
			continue
		}

		providerResult[provider] = result
	}

	return providerResult, nil
}
