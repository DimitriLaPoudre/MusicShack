package hifiv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/mitchellh/mapstructure"
)

func getSearchSong(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchSongData, song string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?s="+url.QueryEscape(song), nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data searchSongData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func getSearchAlbum(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchAlbumData, album string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?al="+url.QueryEscape(album), nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data searchAlbumData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func getSearchArtist(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchArtistData, artist string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?a="+url.QueryEscape(artist), nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data searchArtistData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) Search(ctx context.Context, userId uint, song, album, artist string) (models.SearchData, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	var songData searchSongData
	var albumData searchAlbumData
	var artistData searchArtistData

	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.SearchData{}, fmt.Errorf("HifiV2.Search: %w", err)
	}

	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchSongData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(instances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range instances {
			go getSearchSong(routineCtx, &wgRoutine, instance.Url, ch, song)
		}
		select {
		case find, ok := <-ch:
			cancel()
			if ok {
				songData = find
			}
		case <-ctx.Done():
			cancel()
		}

	}()
	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchAlbumData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(instances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range instances {
			go getSearchAlbum(routineCtx, &wgRoutine, instance.Url, ch, album)
		}
		select {
		case find, ok := <-ch:
			cancel()
			if ok {
				albumData = find
			}
		case <-ctx.Done():
			cancel()
		}
	}()
	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchArtistData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(instances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range instances {
			go getSearchArtist(routineCtx, &wgRoutine, instance.Url, ch, artist)
		}
		select {
		case find, ok := <-ch:
			cancel()
			if ok {
				artistData = find
			}
		case <-ctx.Done():
			cancel()
		}
	}()

	wg.Wait()

	select {
	case <-ctx.Done():
		return models.SearchData{}, fmt.Errorf("HifiV2.Search: %w", context.Canceled)
	default:
	}

	var result models.SearchData

	if len(songData.Data.Songs) != 0 {
		for index, value := range songData.Data.Songs {
			songData.Data.Songs[index].CoverUrl = utils.GetImageURL(value.Album.CoverUrl, 640)
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           &result.Songs,
			TagName:          "useless",
			WeaklyTypedInput: true,
		})
		if err != nil {
			return models.SearchData{}, fmt.Errorf("HifiV2.Search: Song: mapstructure.NewDecoder: %w", err)
		}
		if err := decoder.Decode(songData.Data.Songs); err != nil {
			return models.SearchData{}, fmt.Errorf("HifiV2.Search: Song: mapstructure.Decode: %w", err)
		}
	}

	if len(albumData.Data.Albums.Albums) != 0 {
		for index, value := range albumData.Data.Albums.Albums {
			albumData.Data.Albums.Albums[index].CoverUrl = utils.GetImageURL(value.CoverUrl, 640)
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           &result.Albums,
			TagName:          "useless",
			WeaklyTypedInput: true,
		})
		if err != nil {
			return models.SearchData{}, fmt.Errorf("HifiV2.Search: Album: mapstructure.NewDecoder: %w", err)
		}
		if err := decoder.Decode(albumData.Data.Albums.Albums); err != nil {
			return models.SearchData{}, fmt.Errorf("HifiV2.Search: Album: mapstructure.Decode: %w", err)
		}
	}

	if len(artistData.Data.Artists.Artists) != 0 {
		for index, value := range artistData.Data.Artists.Artists {
			if value.PictureUrl == "" {
				value.PictureUrl = value.PictureUrlFallback
			}
			artistData.Data.Artists.Artists[index].PictureUrl = utils.GetImageURL(value.PictureUrl, 750)
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           &result.Artists,
			TagName:          "useless",
			WeaklyTypedInput: true,
		})
		if err != nil {
			return models.SearchData{}, fmt.Errorf("HifiV2.Search: Artist: mapstructure.NewDecoder: %w", err)
		}
		if err := decoder.Decode(artistData.Data.Artists.Artists); err != nil {
			return models.SearchData{}, fmt.Errorf("HifiV2.Search: Artist: mapstructure.Decode: %w", err)
		}
	}

	return result, nil
}
