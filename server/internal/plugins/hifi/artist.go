package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func fetchArtistInfo(ctx context.Context, url string, id string) (artistInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/artist/?id="+id, nil)
	if err != nil {
		return artistInfo{}, fmt.Errorf("fetchArtistInfo: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return artistInfo{}, fmt.Errorf("fetchArtistInfo: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return artistInfo{}, fmt.Errorf("fetchArtistInfo: %w", errors.New("http error "+strconv.FormatInt(400, 10)))
	}

	var data artistInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return artistInfo{}, fmt.Errorf("fetchArtistInfo: json.Decode: %w", err)
	}

	return data, nil
}

func getArtistInfo(ctx context.Context, instances []models.Instance, id string) (artistInfo, error) {
	type res struct {
		data artistInfo
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchArtistInfo(ctx, url, id)
			ch <- res{data: data, err: err}
		}(instance.Url)
	}

	var lastErr error
	for range instances {
		select {
		case res := <-ch:
			if res.err == nil {
				return res.data, nil
			}
			lastErr = res.err
		case <-ctx.Done():
			return artistInfo{}, ctx.Err()
		}
	}
	return artistInfo{}, fmt.Errorf("getArtistInfo: %w", lastErr)
}

func fetchArtistAlbums(ctx context.Context, url string, id string) (artistAlbums, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/artist/?f="+id+"&skip_tracks=1", nil)
	if err != nil {
		return artistAlbums{}, fmt.Errorf("fetchArtistAlbums: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return artistAlbums{}, fmt.Errorf("fetchArtistAlbums: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return artistAlbums{}, fmt.Errorf("fetchArtistAlbums: %w", errors.New("http error "+strconv.FormatInt(400, 10)))
	}

	var data artistAlbums
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return artistAlbums{}, fmt.Errorf("fetchArtistAlbums: json.Decode: %w", err)
	}

	return data, nil
}

func getArtistAlbums(ctx context.Context, instances []models.Instance, id string) (artistAlbums, error) {
	type res struct {
		data artistAlbums
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchArtistAlbums(ctx, url, id)
			ch <- res{data: data, err: err}
		}(instance.Url)
	}

	var lastErr error
	for range instances {
		select {
		case res := <-ch:
			if res.err == nil {
				return res.data, nil
			}
			lastErr = res.err
		case <-ctx.Done():
			return artistAlbums{}, ctx.Err()
		}
	}
	return artistAlbums{}, fmt.Errorf("getArtistInfo: %w", lastErr)
}

func getArtistData(ctx context.Context, instances []models.Instance, id string) (artistInfo, artistAlbums, error) {
	type res struct {
		data any
		err  error
	}

	ch := make(chan res, 2)
	go func() {
		info, err := getArtistInfo(ctx, instances, id)
		ch <- res{data: info, err: err}
	}()
	go func() {
		albums, err := getArtistAlbums(ctx, instances, id)
		ch <- res{data: albums, err: err}
	}()

	var info artistInfo
	var albums artistAlbums
	var infoErr, albumsErr error
	for range 2 {
		res := <-ch
		switch v := res.data.(type) {
		case artistInfo:
			info = v
			infoErr = res.err
		case artistAlbums:
			albums = v
			albumsErr = res.err
		}
	}

	if infoErr != nil {
		return artistInfo{}, artistAlbums{}, fmt.Errorf("getArtistData: %w", infoErr)
	}
	if albumsErr != nil {
		return artistInfo{}, artistAlbums{}, fmt.Errorf("getArtistData: %w", albumsErr)
	}

	return info, albums, nil
}

func (p *Hifi) Artist(ctx context.Context, userId uint, id string) (models.ArtistData, error) {
	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.ArtistData{}, fmt.Errorf("Hifi.Artist: %w", err)
	}

	artistInfo, artistAlbums, err := getArtistData(ctx, instances, id)
	if err != nil {
		return models.ArtistData{}, fmt.Errorf("Hifi.Artist: %w", err)
	}

	normalizeArtistData := models.ArtistData{
		Provider: p.Provider(),
		Api:      p.Name(),
		Id:       strconv.FormatUint(uint64(artistInfo.Artist.Id), 10),
		Name:     artistInfo.Artist.Name,
		Albums:   make([]models.ArtistDataAlbum, 0),
		Ep:       make([]models.ArtistDataAlbum, 0),
		Singles:  make([]models.ArtistDataAlbum, 0),
	}

	if artistInfo.Artist.PictureUrl == "" {
		artistInfo.Artist.PictureUrl = artistInfo.Artist.PictureUrlFallback
	}
	normalizeArtistData.PictureUrl = utils.GetImageURL(artistInfo.Artist.PictureUrl, 750)

	best := make(map[albumItemComparaison]*albumItem)
	for _, album := range artistAlbums.Albums.Items {
		if bestVersion, ok := best[albumItemComparaison{
			Title:       strings.ToLower(album.Title),
			ReleaseDate: album.ReleaseDate,
			TrackNumber: album.NumberOfTracks,
		}]; !ok || (!bestVersion.Explicit && album.Explicit) || (bestVersion.Explicit == album.Explicit && len(bestVersion.MediaMetadata.Tags) < len(album.MediaMetadata.Tags)) {
			best[albumItemComparaison{
				Title:       strings.ToLower(album.Title),
				ReleaseDate: album.ReleaseDate,
				TrackNumber: album.NumberOfTracks,
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

	for _, rawAlbum := range list {
		album := models.ArtistDataAlbum{
			Id:           strconv.FormatUint(uint64(rawAlbum.Id), 10),
			Title:        rawAlbum.Title,
			Duration:     rawAlbum.Duration,
			ReleaseDate:  rawAlbum.ReleaseDate,
			CoverUrl:     utils.GetImageURL(rawAlbum.CoverUrl, 1280),
			AudioQuality: models.QualityHigh,
			Explicit:     rawAlbum.Explicit,
			Artists:      make([]models.AlbumDataArtist, 0),
		}

		for _, quality := range rawAlbum.MediaMetadata.Tags {
			switch quality {
			case "HIRES_LOSSLESS":
				album.AudioQuality = max(album.AudioQuality, models.QualityHiresLossless)
			case "LOSSLESS", "DOLBY_ATMOS":
				album.AudioQuality = max(album.AudioQuality, models.QualityLossless)
			}
		}

		for _, artist := range rawAlbum.Artists {
			album.Artists = append(album.Artists, models.AlbumDataArtist{
				Id:   strconv.FormatUint(uint64(artist.Id), 10),
				Name: artist.Name,
			})
		}

		switch rawAlbum.Type {
		case "ALBUM":
			normalizeArtistData.Albums = append(normalizeArtistData.Albums, album)
		case "EP":
			normalizeArtistData.Ep = append(normalizeArtistData.Ep, album)
		case "SINGLE":
			normalizeArtistData.Singles = append(normalizeArtistData.Singles, album)
		}
	}

	return normalizeArtistData, nil
}
