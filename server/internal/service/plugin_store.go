package service

import (
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/rs/zerolog"
)

type PluginStoreService struct {
	l    *zerolog.Logger
	name map[string]model.Plugin
}

func NewPluginStoreService(l *zerolog.Logger) PluginStoreService {
	return PluginStoreService{
		l:    l,
		name: map[string]model.Plugin{},
	}
}

func (store *PluginStoreService) Register(p model.Plugin) {
	store.name[p.Name()] = p
}

func (store *PluginStoreService) GetPluginByName(name string) (model.Plugin, bool) {
	p, ok := store.name[name]
	return p, ok
}

func (store *PluginStoreService) GetAllPluginsByName() map[string]model.Plugin {
	return store.name
}
