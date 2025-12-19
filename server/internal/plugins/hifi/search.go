package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/mitchellh/mapstructure"
)

type searchSongData struct {
	Songs []struct {
		Id       uint   `json:"id"`
		Title    string `json:"title"`
		CoverUrl string
		Artists  []struct {
			Id   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"artists"`
		Album struct {
			CoverUrl string `json:"cover"`
		} `json:"album"`
	} `json:"items"`
}

type searchAlbumData struct {
	Section struct {
		Albums []struct {
			Id       uint   `json:"id"`
			Title    string `json:"title"`
			CoverUrl string `json:"cover"`
			Artists  []struct {
				Id   uint   `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"albums"`
}

type searchArtistData struct {
	Section struct {
		Artists []struct {
			Id                 uint   `json:"id"`
			Name               string `json:"name"`
			PictureUrl         string `json:"picture"`
			PictureUrlFallback string `json:"selectedAlbumCoverFallback"`
		} `json:"items"`
	} `json:"artists"`
}

func (p *Hifi) getSearchSong(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchSongData, song string) {
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

func (p *Hifi) getSearchAlbum(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchAlbumData, album string) {
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

func (p *Hifi) getSearchArtist(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchArtistData, artist string) {
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

	var tmp []searchArtistData
	if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
		return
	}
	ch <- tmp[0]
}

func (p *Hifi) Search(ctx context.Context, userId uint, song, album, artist string) (models.SearchData, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	var songData searchSongData
	var albumData searchAlbumData
	var artistData searchArtistData

	apiInstances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.SearchData{}, err
	}

	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchSongData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(apiInstances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range apiInstances {
			go p.getSearchSong(routineCtx, &wgRoutine, instance.Url, ch, song)
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

		wgRoutine.Add(len(apiInstances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range apiInstances {
			go p.getSearchAlbum(routineCtx, &wgRoutine, instance.Url, ch, album)
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

		wgRoutine.Add(len(apiInstances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range apiInstances {
			go p.getSearchArtist(routineCtx, &wgRoutine, instance.Url, ch, artist)
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
		return models.SearchData{}, errors.New("connection closed")
	default:
	}

	var result models.SearchData

	if len(songData.Songs) != 0 {
		for index, value := range songData.Songs {
			if value.Album.CoverUrl != "" {
				songData.Songs[index].CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.Album.CoverUrl, "-", "/") + "/160x160.jpg"
			}
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           &result.Songs,
			TagName:          "useless",
			WeaklyTypedInput: true,
		})
		if err != nil {
			return models.SearchData{}, err
		}
		if err := decoder.Decode(songData.Songs); err != nil {
			return models.SearchData{}, err
		}
	}

	if len(albumData.Section.Albums) != 0 {
		for index, value := range albumData.Section.Albums {
			if value.CoverUrl != "" {
				albumData.Section.Albums[index].CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.CoverUrl, "-", "/") + "/160x160.jpg"
			}
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           &result.Albums,
			TagName:          "useless",
			WeaklyTypedInput: true,
		})
		if err != nil {
			return models.SearchData{}, err
		}
		if err := decoder.Decode(albumData.Section.Albums); err != nil {
			return models.SearchData{}, err
		}
	}

	if len(artistData.Section.Artists) != 0 {
		for index, value := range artistData.Section.Artists {
			if value.PictureUrl != "" {
				artistData.Section.Artists[index].PictureUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.PictureUrl, "-", "/") + "/160x160.jpg"
			} else if value.PictureUrlFallback != "" {
				artistData.Section.Artists[index].PictureUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.PictureUrlFallback, "-", "/") + "/160x160.jpg"
			}
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           &result.Artists,
			TagName:          "useless",
			WeaklyTypedInput: true,
		})
		if err != nil {
			return models.SearchData{}, err
		}
		if err := decoder.Decode(artistData.Section.Artists); err != nil {
			return models.SearchData{}, err
		}
	}

	return result, nil
}
