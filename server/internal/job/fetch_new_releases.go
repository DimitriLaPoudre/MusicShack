package job

import (
	"context"
	"log/slog"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/robfig/cron/v3"
)

func DownloadNewReleases(c *cron.Cron, ctx context.Context, s *service.FetchNewReleasesService) error {
	if _, err := c.AddFunc("0 1 * * *", func() {
		lastFetchDate := time.Now().Add(-1 * time.Hour * 24)
		if errList := s.FetchNewReleases(ctx, lastFetchDate); len(errList) > 0 {
			slog.Error("download new releases", slog.Any("err", errList))
		}
	}); err != nil {
		return err
	}

	return nil
}
