package usecase

import (
	"context"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
)

type PluginUseCase struct {
	plugin   *service.PluginService
	store    *service.PluginStoreService
	instance model.InstanceRepository
}

func NewPluginUseCase(plugin *service.PluginService, store *service.PluginStoreService, instance model.InstanceRepository) PluginUseCase {
	return PluginUseCase{
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

	pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

	song, err := u.plugin.GetSong(ctx, pluginInstances, id)
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

	pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

	album, err := u.plugin.GetAlbum(ctx, pluginInstances, id)
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

	pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

	artist, err := u.plugin.GetArtist(ctx, pluginInstances, id)
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

	pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

	playlist, err := u.plugin.GetPlaylist(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedPlaylist{}, err
	}

	return playlist, nil
}

func (u *PluginUseCase) GetSearch(ctx context.Context, user model.User, q string) (map[string]model.EnrichedSearch, error) {
	if q == "" {
		return nil, model.ErrPluginSearchEmptyQuery
	}

	providerResult := map[string]model.EnrichedSearch{}

	for provider := range u.store.ListPluginsByProvider() {
		instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
		if err != nil {
			continue
		}

		pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

		result, err := u.plugin.Search(ctx, pluginInstances, q)
		if err != nil {
			continue
		}

		providerResult[provider] = result
	}

	return providerResult, nil
}
