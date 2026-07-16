package service

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

type PluginStoreService struct {
	name     map[string]model.Plugin
	provider map[string][]model.Plugin
}

func NewPluginStoreService() PluginStoreService {
	return PluginStoreService{
		name:     map[string]model.Plugin{},
		provider: map[string][]model.Plugin{},
	}
}

func (store *PluginStoreService) Register(plugins ...model.Plugin) {
	for _, p := range plugins {
		store.name[p.Name()] = p
		store.provider[p.Provider()] = append(store.provider[p.Provider()], p)
		// slices.SortFunc(store.provider[p.Provider()], func(a, b models.Plugin) int {
		// 	return b.Priority() - a.Priority()
		// })
	}
}

func (store *PluginStoreService) GetPluginByName(name string) (model.Plugin, bool) {
	p, ok := store.name[name]
	return p, ok
}

func (store *PluginStoreService) ListPluginsByName() map[string]model.Plugin {
	return store.name
}

func (store *PluginStoreService) GetPluginByProvider(provider string) ([]model.Plugin, bool) {
	p, ok := store.provider[provider]
	return p, ok
}

func (store *PluginStoreService) ListPluginsByProvider() map[string][]model.Plugin {
	return store.provider
}
