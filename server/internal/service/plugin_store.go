package service

import (
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/rs/zerolog"
)

type PluginStoreService struct {
	l        *zerolog.Logger
	name     map[string]model.Plugin
	provider map[string][]model.Plugin
}

func NewPluginStoreService(l *zerolog.Logger) PluginStoreService {
	return PluginStoreService{
		l:        l,
		name:     map[string]model.Plugin{},
		provider: map[string][]model.Plugin{},
	}
}

func (store *PluginStoreService) Register(p model.Plugin) {
	store.name[p.Name()] = p
	store.provider[p.Provider()] = append(store.provider[p.Provider()], p)
	// slices.SortFunc(store.provider[p.Provider()], func(a, b models.Plugin) int {
	// 	return b.Priority() - a.Priority()
	// })
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
