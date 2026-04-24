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
	userId   uint
	provider string
	albumId  string
}

func getNewReleasesOfArtist(ctx context.Context, userId uint, provider string, id string, lastFetchDate string) ([]release, error) {
	var newReleases []release
	plugins, ok := plugins.GetPluginByProvider(provider)
	if !ok {
		return newReleases, fmt.Errorf("invalid provider name")
	}

	var artist models.ArtistData
	var err error
	for _, plugin := range plugins {
		artist, err = plugin.Artist(ctx, userId, id)
		if err == nil {
			break
		}
	}
	if err != nil {
		return newReleases, err
	}

	for _, album := range artist.Albums {
		if album.ReleaseDate > lastFetchDate {
			newReleases = append(newReleases, release{
				userId:   userId,
				provider: artist.Provider,
				albumId:  album.Id,
			})
		}
	}

	for _, ep := range artist.Ep {
		if ep.ReleaseDate > lastFetchDate {
			newReleases = append(newReleases, release{
				userId:   userId,
				provider: artist.Provider,
				albumId:  ep.Id,
			})
		}
	}

	for _, single := range artist.Singles {
		if single.ReleaseDate > lastFetchDate {
			newReleases = append(newReleases, release{
				userId:   userId,
				provider: artist.Provider,
				albumId:  single.Id,
			})
		}
	}

	return newReleases, nil
}

func getNewReleases(ctx context.Context, follows []models.Follow, lastFetchDate string) ([]release, error) {
	var newReleases []release
	for _, follow := range follows {
		tmpReleases, _ := getNewReleasesOfArtist(ctx, follow.UserId, follow.Provider, follow.ArtistId, lastFetchDate)
		newReleases = append(newReleases, tmpReleases...)
	}
	return newReleases, nil
}

func fetch(ctx context.Context, lastFetchDate string) error {
	follows, err := repository.ListFollows()
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	releases, err := getNewReleases(ctx, follows, lastFetchDate)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	for _, release := range releases {
		services.DownloadManager.AddAlbum(release.userId, release.provider, release.albumId)
	}
	return nil
}

func AutoFetch(ctx context.Context) *cron.Cron {
	c := cron.New()
	if _, err := c.AddFunc("0 1 * * *", func() {
		lastFetchDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		for try := range 3 {
			if err := fetch(ctx, lastFetchDate); err != nil {
				log.Println("AutoFetch: try ", try, ": ", err)
			} else {
				log.Println("AutoFetch: success")
				break
			}
		}
	}); err != nil {
		log.Fatalln("AutoFetch: AddFunc: ", err)
	}

	c.Start()
	return c
}
