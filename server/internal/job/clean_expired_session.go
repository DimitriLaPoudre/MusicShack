package job

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

func cleanExpiredSession(ctx context.Context, repo model.SessionRepository) error {
	return repo.DeleteSessionExpired(ctx)
}

func CleanExpiredSession(c *cron.Cron, ctx context.Context, l *zerolog.Logger, repo model.SessionRepository) error {
	if _, err := c.AddFunc("0 0 * * *", func() {
		if err := cleanExpiredSession(ctx, repo); err != nil {
			l.Err(err).Msg("")
		}
	}); err != nil {
		return err
	}

	c.Start()
	return nil
}
