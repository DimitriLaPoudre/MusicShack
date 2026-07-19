package usecase

import (
	"context"
	"fmt"

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
		return model.EnrichedSong{}, fmt.Errorf("list instances of user %s for provider %s: %w", user.ID.String(), provider, err)
	}

	pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

	song, err := u.plugin.GetSong(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedSong{}, fmt.Errorf("get song %s from each plugin: %w", id, err)
	}

	return song, nil
}

func (u *PluginUseCase) GetAlbum(ctx context.Context, user model.User, provider string, id string) (model.EnrichedAlbum, error) {
	instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.EnrichedAlbum{}, fmt.Errorf("list instances of user %s for provider %s: %w", user.ID.String(), provider, err)
	}

	pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

	album, err := u.plugin.GetAlbum(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedAlbum{}, fmt.Errorf("get album %s from each plugin: %w", id, err)
	}

	return album, nil
}

func (u *PluginUseCase) GetArtist(ctx context.Context, user model.User, provider string, id string) (model.EnrichedArtist, error) {
	instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.EnrichedArtist{}, fmt.Errorf("list instances of user %s for provider %s: %w", user.ID.String(), provider, err)
	}

	pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

	artist, err := u.plugin.GetArtist(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedArtist{}, fmt.Errorf("get artist %s from each plugin: %w", id, err)
	}

	return artist, nil
}

func (u *PluginUseCase) GetPlaylist(ctx context.Context, user model.User, provider string, id string) (model.EnrichedPlaylist, error) {
	instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.EnrichedPlaylist{}, fmt.Errorf("list instances of user %s for provider %s: %w", user.ID.String(), provider, err)
	}

	pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

	playlist, err := u.plugin.GetPlaylist(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedPlaylist{}, fmt.Errorf("get playlist %s from each: %w", id, err)
	}

	return playlist, nil
}

func (u *PluginUseCase) GetSearch(ctx context.Context, user model.User, q string) (model.SearchResult, error) {
	if q == "" {
		return model.SearchResult{}, model.ErrPluginSearchEmptyQuery
	}

	providerResult := map[string]model.EnrichedSearch{}

	for provider := range u.store.ListPluginsByProvider() {
		instances, err := u.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
		if err != nil {
			continue
		}

		pluginInstances := u.plugin.InstancesToMapPluginInstances(instances)

		if item, err := u.plugin.Url(ctx, pluginInstances, q); err == nil {
			return model.SearchResult{ItemFound: item}, nil
		}

		result, err := u.plugin.Search(ctx, pluginInstances, q)
		if err != nil {
			continue
		}

		providerResult[provider] = result
	}

	return model.SearchResult{ProviderResult: providerResult}, nil
}
