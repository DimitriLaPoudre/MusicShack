package hifiv2_2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func getArtistInfo(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- artistData, id string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/artist/?id="+id, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return
	}

	var data artistData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func getArtistAlbums(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- artistAlbums, id string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/artist/?f="+id, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return
	}

	var data artistAlbums
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *Hifi) Artist(ctx context.Context, userId uint, id string) (models.ArtistData, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var data artistData
	var albumsData artistAlbums

	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.ArtistData{}, fmt.Errorf("Hifi.Artist: %w", err)
	}

	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan artistData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(instances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range instances {
			go getArtistInfo(routineCtx, &wgRoutine, instance.Url, ch, id)
		}
		select {
		case find, ok := <-ch:
			cancel()
			if ok {
				data = find
			}
		case <-ctx.Done():
			cancel()
		}
	}()
	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan artistAlbums)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(instances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range instances {
			go getArtistAlbums(routineCtx, &wgRoutine, instance.Url, ch, id)
		}
		select {
		case find, ok := <-ch:
			cancel()
			if ok {
				albumsData = find
			}
		case <-ctx.Done():
			cancel()
		}
	}()

	wg.Wait()

	select {
	case <-ctx.Done():
		return models.ArtistData{}, fmt.Errorf("Hifi.Artist: %w", context.Canceled)
	default:
	}

	if data.Artist.Id == 0 {
		return models.ArtistData{}, fmt.Errorf("Hifi.Artist: %w", errors.New("can't fetch"))
	}

	normalizeArtistData := models.ArtistData{
		Api:     p.Name(),
		Id:      strconv.FormatUint(uint64(data.Artist.Id), 10),
		Name:    data.Artist.Name,
		Albums:  make([]models.ArtistDataAlbum, 0),
		Ep:      make([]models.ArtistDataAlbum, 0),
		Singles: make([]models.ArtistDataAlbum, 0),
	}
	{
		if data.Artist.PictureUrl == "" {
			data.Artist.PictureUrl = data.Artist.PictureUrlFallback
		}
		normalizeArtistData.PictureUrl = utils.GetImageURL(data.Artist.PictureUrl, 750)

		best := make(map[albumItemComparaison]*albumItem)
		for _, album := range albumsData.Albums.Items {
			if bestVersion, ok := best[albumItemComparaison{
				Title:       album.Title,
				ReleaseDate: album.ReleaseDate,
			}]; !ok || (!bestVersion.Explicit && album.Explicit) || (bestVersion.Explicit == album.Explicit && len(bestVersion.MediaMetadata.Tags) < len(album.MediaMetadata.Tags)) {
				best[albumItemComparaison{
					Title:       album.Title,
					ReleaseDate: album.ReleaseDate,
				}] = &album
			}
		}

		list := []*albumItem{}
		for _, album := range best {
			list = append(list, album)
		}
		slices.SortFunc(list, func(a, b *albumItem) int {
			if a.ReleaseDate > b.ReleaseDate {
				return -1
			}
			if a.ReleaseDate < b.ReleaseDate {
				return 1
			}
			return 0
		})

		for _, album := range list {
			audioQuality := models.QualityHigh
			for _, quality := range album.MediaMetadata.Tags {
				switch quality {
				case "HIRES_LOSSLESS":
					audioQuality = max(audioQuality, models.QualityHiresLossless)
				case "LOSSLESS", "DOLBY_ATMOS":
					audioQuality = max(audioQuality, models.QualityLossless)
				}
			}

			artists := make([]models.AlbumDataArtist, 0)
			for _, artist := range album.Artists {
				artists = append(artists, models.AlbumDataArtist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
			}

			switch album.Type {
			case "ALBUM":
				normalizeArtistData.Albums = append(normalizeArtistData.Albums, models.ArtistDataAlbum{
					Id:           strconv.FormatUint(uint64(album.Id), 10),
					Title:        album.Title,
					Duration:     album.Duration,
					ReleaseDate:  album.ReleaseDate,
					CoverUrl:     utils.GetImageURL(album.CoverUrl, 640),
					AudioQuality: audioQuality,
					Artists:      artists,
				})
			case "EP":
				normalizeArtistData.Ep = append(normalizeArtistData.Ep, models.ArtistDataAlbum{
					Id:           strconv.FormatUint(uint64(album.Id), 10),
					Title:        album.Title,
					Duration:     album.Duration,
					ReleaseDate:  album.ReleaseDate,
					CoverUrl:     utils.GetImageURL(album.CoverUrl, 640),
					AudioQuality: audioQuality,
					Artists:      artists,
				})
			case "SINGLE":
				normalizeArtistData.Singles = append(normalizeArtistData.Singles, models.ArtistDataAlbum{
					Id:           strconv.FormatUint(uint64(album.Id), 10),
					Title:        album.Title,
					Duration:     album.Duration,
					ReleaseDate:  album.ReleaseDate,
					CoverUrl:     utils.GetImageURL(album.CoverUrl, 640),
					AudioQuality: audioQuality,
					Artists:      artists,
				})
			}
		}
	}

	return normalizeArtistData, nil
}
