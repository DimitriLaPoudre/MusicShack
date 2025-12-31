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

	apiList := plugins.GetRegistry()

	for _, api := range apiList {
		go func(api models.Plugin) {
			if err := api.Status(ctx, url); err == nil {
				select {
				case ch <- api:
					cancel()
				default:
				}
			}
		}(api)
	}

	select {
	case api := <-ch:
		return api
	case <-ctx.Done():
		return nil
	}
}
