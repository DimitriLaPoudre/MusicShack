package job

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type cleanExpiredSessionRepository interface {
	DeleteExpiredSessions(ctx context.Context) error
}

func cleanExpiredSession(ctx context.Context, repo cleanExpiredSessionRepository) error {
	if err := repo.DeleteExpiredSessions(ctx); err != nil {
		return err
	}
	return nil
}

func CleanExpiredSession(c *cron.Cron, ctx context.Context, l *zerolog.Logger, repo cleanExpiredSessionRepository) error {
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
