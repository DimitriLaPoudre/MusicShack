package job

import (
	"context"
	"log/slog"

	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/robfig/cron/v3"
)

func CleanExpiredSession(c *cron.Cron, ctx context.Context, auth *service.AuthService) error {
	if _, err := c.AddFunc("0 0 * * *", func() {
		if err := auth.CleanExpiredSession(ctx); err != nil {
			slog.Error("failed to clean expired session", slog.String("err", err.Error()))
		}
	}); err != nil {
		return err
	}

	return nil
}
