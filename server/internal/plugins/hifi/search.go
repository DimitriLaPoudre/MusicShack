package hifi

import (
	"context"
	"encoding/json"
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
			Id         uint   `json:"id"`
			Name       string `json:"name"`
			PictureUrl string `json:"picture"`
		} `json:"items"`
	} `json:"artists"`
}

func (p *Hifi) getSearchSong(ctx context.Context, urlApi string, ch chan<- searchSongData, song string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?s="+url.QueryEscape(song), nil)
	if err != nil {
		println(err.Error())
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

func (p *Hifi) getSearchAlbum(ctx context.Context, urlApi string, ch chan<- searchAlbumData, album string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?al="+url.QueryEscape(album), nil)
	if err != nil {
		println(err.Error())
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

func (p *Hifi) getSearchArtist(ctx context.Context, urlApi string, ch chan<- searchArtistData, artist string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?a="+url.QueryEscape(artist), nil)
	if err != nil {
		println(err.Error())
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

func (p *Hifi) Search(ctx context.Context, song, album, artist string) (models.SearchData, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	var songData searchSongData
	var albumData searchAlbumData
	var artistData searchArtistData

	apiInstances, err := repository.ListApiInstancesByApi(p.Name())
	if err != nil {
		return models.SearchData{}, err
	}

	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchSongData)

		for _, instance := range apiInstances {
			go p.getSearchSong(routineCtx, instance.Url, ch, song)
		}
		select {
		case find := <-ch:
			cancel()
			songData = find
		case <-ctx.Done():
			cancel()
		}

	}()
	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchAlbumData)

		for _, instance := range apiInstances {
			go p.getSearchAlbum(routineCtx, instance.Url, ch, album)
		}
		select {
		case find := <-ch:
			cancel()
			albumData = find
		case <-ctx.Done():
			cancel()
		}
	}()
	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchArtistData)

		for _, instance := range apiInstances {
			go p.getSearchArtist(routineCtx, instance.Url, ch, artist)
		}
		select {
		case find := <-ch:
			cancel()
			artistData = find
		case <-ctx.Done():
			cancel()
		}
	}()

	wg.Wait()

	var result models.SearchData

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

	for index, value := range albumData.Section.Albums {
		if value.CoverUrl != "" {
			albumData.Section.Albums[index].CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.CoverUrl, "-", "/") + "/160x160.jpg"
		}
	}
	decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
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

	for index, value := range artistData.Section.Artists {
		if value.PictureUrl != "" {
			artistData.Section.Artists[index].PictureUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.PictureUrl, "-", "/") + "/160x160.jpg"
		}
	}
	decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
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

	return result, nil
}
