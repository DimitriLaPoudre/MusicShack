package autofetch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/robfig/cron/v3"
)

type release struct {
	userId  uint
	api     models.Plugin
	albumId string
}

func getNewReleasesOfArtist(ctx context.Context, userId uint, api string, id string, lastFetchDate string) ([]release, error) {
	var newReleases []release
	p, ok := plugins.Get(api)
	if !ok {
		return newReleases, fmt.Errorf("api name invalid")
	}

	artist, err := p.Artist(ctx, userId, id)
	if err != nil {
		return newReleases, err
	}

	for _, album := range artist.Albums {
		if album.ReleaseDate > lastFetchDate {
			newReleases = append(newReleases, release{
				userId:  userId,
				api:     p,
				albumId: album.Id,
			})
		}
	}

	for _, ep := range artist.Ep {
		if ep.ReleaseDate > lastFetchDate {
			newReleases = append(newReleases, release{
				userId:  userId,
				api:     p,
				albumId: ep.Id,
			})
		}
	}

	for _, single := range artist.Singles {
		if single.ReleaseDate > lastFetchDate {
			newReleases = append(newReleases, release{
				userId:  userId,
				api:     p,
				albumId: single.Id,
			})
		}
	}

	return newReleases, nil
}

func getNewReleases(ctx context.Context, follows []models.Follow, lastFetchDate string) ([]release, error) {
	var newReleases []release
	for _, follow := range follows {
		tmpReleases, _ := getNewReleasesOfArtist(ctx, follow.UserId, follow.Api, follow.ArtistId, lastFetchDate)
		newReleases = append(newReleases, tmpReleases...)
	}
	return newReleases, nil
}

func fetch(ctx context.Context) error {
	lastFetchDate := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02")
	follows, err := repository.ListFollows()

	releases, err := getNewReleases(ctx, follows, lastFetchDate)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	for _, release := range releases {
		services.DownloadManager.AddAlbum(release.userId, release.api, release.albumId, "")
	}
	return nil
}

func AutoFetch(ctx context.Context) *cron.Cron {
	c := cron.New()

	c.AddFunc("0 1 * * 5", func() {
		for try := range 3 {
			if err := fetch(ctx); err != nil {
				log.Println("AutoFetch: try ", try, ": ", err)
			} else {
				log.Println("AutoFetch: success")
				break
			}
		}
	})

	c.Start()
	return c
}
