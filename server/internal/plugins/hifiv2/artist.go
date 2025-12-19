package hifiv2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/mitchellh/mapstructure"
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

	var data artistAlbums
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) Artist(ctx context.Context, userId uint, id string) (models.ArtistData, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var data artistData
	var albums artistAlbums

	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.ArtistData{}, fmt.Errorf("HifiV2.Artist: %w", err)
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
				albums = find
			}
		case <-ctx.Done():
			cancel()
		}
	}()

	wg.Wait()

	select {
	case <-ctx.Done():
		return models.ArtistData{}, fmt.Errorf("HifiV2.Artist: %w", context.Canceled)
	default:
	}

	var normalizeArtistData models.ArtistData

	if data.Artist.Id == 0 {
		return models.ArtistData{}, fmt.Errorf("HifiV2.Artist: %w", errors.New("can't fetch"))
	}

	if data.Artist.PictureUrl == "" {
		data.Artist.PictureUrl = data.Artist.PictureUrlFallback
	}
	data.Artist.PictureUrl = utils.GetImageURL(data.Artist.PictureUrl, 750)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeArtistData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.ArtistData{}, fmt.Errorf("HifiV2.Artist: mapstructure.NewDecoder: %w", err)
	}
	if err := decoder.Decode(data.Artist); err != nil {
		return models.ArtistData{}, fmt.Errorf("HifiV2.Artist: mapstructure.Decode: %w", err)
	}

	if albums.Albums.Id == "" {
		return normalizeArtistData, nil
	}

	for i, item := range albums.Albums.Rows[0].Modules[0].PagedList.Items {
		albums.Albums.Rows[0].Modules[0].PagedList.Items[i].CoverUrl = utils.GetImageURL(item.CoverUrl, 640)
		albums.Albums.Rows[0].Modules[0].PagedList.Items[i].ReleaseDate = item.StreamStartDate[:10]
		albums.Albums.Rows[0].Modules[0].PagedList.Items[i].StreamStartDate = item.StreamStartDate[:10]
	}

	decoder, _ = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeArtistData.Albums,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	decoder.Decode(albums.Albums.Rows[0].Modules[0].PagedList.Items)

	return normalizeArtistData, nil
}
