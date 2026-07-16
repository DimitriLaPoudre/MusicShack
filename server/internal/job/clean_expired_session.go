package job

import (
	"context"
	"log/slog"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/robfig/cron/v3"
)

func cleanExpiredSession(ctx context.Context, repo model.SessionRepository) error {
	return repo.DeleteSessionExpired(ctx)
}

func CleanExpiredSession(c *cron.Cron, ctx context.Context, repo model.SessionRepository) error {
	if _, err := c.AddFunc("0 0 * * *", func() {
		if err := cleanExpiredSession(ctx, repo); err != nil {
			slog.Error("failed to clean expired session", slog.String("err", err.Error()))
		}
	}); err != nil {
		return err
	}

	c.Start()
	return nil
}
