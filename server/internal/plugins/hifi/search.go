package hifi

import (
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

func (p *Hifi) getSearchSong(wg *sync.WaitGroup, data *searchSongData, song string) {
	defer wg.Done()

	apiInstance, err := repository.GetApiInstanceByApi(p.Name())
	if err != nil {
		return
	}

	resp, err := http.Get(apiInstance.Url + "/search/?s=" + url.QueryEscape(song))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return
	}
}

func (p *Hifi) getSearchAlbum(wg *sync.WaitGroup, data *searchAlbumData, album string) {
	defer wg.Done()

	apiInstance, err := repository.GetApiInstanceByApi(p.Name())
	if err != nil {
		return
	}

	resp, err := http.Get(apiInstance.Url + "/search/?al=" + url.QueryEscape(album))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return
	}
}

func (p *Hifi) getSearchArtist(wg *sync.WaitGroup, data *searchArtistData, artist string) {
	defer wg.Done()

	apiInstance, err := repository.GetApiInstanceByApi(p.Name())
	if err != nil {
		return
	}

	resp, err := http.Get(apiInstance.Url + "/search/?a=" + url.QueryEscape(artist))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var tmp []searchArtistData
	if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
		return
	}

	*data = tmp[0]
}

func (p *Hifi) Search(song, album, artist string) (models.SearchData, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	var searchSongData searchSongData
	var searchAlbumData searchAlbumData
	var searchArtistData searchArtistData

	go p.getSearchSong(&wg, &searchSongData, song)
	go p.getSearchAlbum(&wg, &searchAlbumData, album)
	go p.getSearchArtist(&wg, &searchArtistData, artist)

	wg.Wait()

	var result models.SearchData

	for index, value := range searchSongData.Songs {
		if value.Album.CoverUrl != "" {
			searchSongData.Songs[index].CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.Album.CoverUrl, "-", "/") + "/160x160.jpg"
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
	if err := decoder.Decode(searchSongData.Songs); err != nil {
		return models.SearchData{}, err
	}

	for index, value := range searchAlbumData.Section.Albums {
		if value.CoverUrl != "" {
			searchAlbumData.Section.Albums[index].CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.CoverUrl, "-", "/") + "/160x160.jpg"
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
	if err := decoder.Decode(searchAlbumData.Section.Albums); err != nil {
		return models.SearchData{}, err
	}

	for index, value := range searchArtistData.Section.Artists {
		if value.PictureUrl != "" {
			searchArtistData.Section.Artists[index].PictureUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(value.PictureUrl, "-", "/") + "/160x160.jpg"
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
	if err := decoder.Decode(searchArtistData.Section.Artists); err != nil {
		return models.SearchData{}, err
	}

	return result, nil
}
