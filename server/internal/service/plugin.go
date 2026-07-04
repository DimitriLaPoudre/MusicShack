package service

import (
	"context"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/rs/zerolog"
)

type PluginService struct {
	l     *zerolog.Logger
	store *PluginStoreService
}

func NewPluginService(l *zerolog.Logger, store *PluginStoreService) PluginService {
	return PluginService{
		l:     l,
		store: store,
	}
}

func (s *PluginService) GetOriginalPlugin(ctx context.Context, url string) (model.Plugin, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ch := make(chan model.Plugin, 1)

	for _, plugin := range s.store.ListPluginsByName() {
		go func(p model.Plugin) {
			if err := p.Status(ctx, url); err == nil {
				select {
				case ch <- p:
					cancel()
				default:
				}
			}
		}(plugin)
	}

	select {
	case plugin := <-ch:
		return plugin, nil
	case <-ctx.Done():
		return nil, model.ErrInvalidUrlPlugin
	}
}
