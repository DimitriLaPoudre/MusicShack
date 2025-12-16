package autofetch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/robfig/cron/v3"
)

func getNewReleasesOfArtist(ctx context.Context, api string, id string, lastFetchDate string) ([]models.AlbumData, error) {
	p, ok := plugins.Get(api)
	if !ok {
		return []models.AlbumData{}, fmt.Errorf("api name invalid")
	}

	artist, err := p.Artist(ctx, id)
	if err != nil {
		return []models.AlbumData{}, err
	}

	var newReleases []models.AlbumData
	for _, album := range artist.Albums {
		if album.ReleaseDate > lastFetchDate {
			newReleases = append(newReleases, album)
		}
	}

	return newReleases, nil
}

func getNewReleases(ctx context.Context, follows []models.Follow, lastFetchDate string) ([]models.AlbumData, error) {
	var newReleases []models.AlbumData
	for _, follow := range follows {
		tmpReleases, _ := getNewReleasesOfArtist(ctx, follow.Api, follow.ArtistId, lastFetchDate)
		newReleases = append(newReleases, tmpReleases...)
	}
	return newReleases, nil
}

func fetch(ctx context.Context) error {
	lastFetchDate := time.Now().Format("2006-01-02")
	follows, err := repository.ListFollows()

	releases, err := getNewReleases(ctx, follows, lastFetchDate)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	releases = releases
	return nil
}

func AutoFetch(ctx context.Context) *cron.Cron {
	c := cron.New()

	c.AddFunc("* * * * *", func() {
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
