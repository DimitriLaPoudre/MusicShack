package services

import (
	"context"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
)

func TestApi(ctx context.Context, url string) models.Plugin {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ch := make(chan models.Plugin, 1)

	for _, plugin := range plugins.GetAll() {
		go func(p models.Plugin) {
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
	case api := <-ch:
		return api
	case <-ctx.Done():
		return nil
	}
}
