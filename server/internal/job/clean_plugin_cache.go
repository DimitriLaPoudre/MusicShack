package job

import (
	"context"

	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/robfig/cron/v3"
)

func CleanExpiredPluginCache(c *cron.Cron, ctx context.Context, cache *service.PluginCacheService) error {
	if _, err := c.AddFunc("0 0 * * *", func() {
		cache.CleanExpired()
	}); err != nil {
		return err
	}

	return nil
}
