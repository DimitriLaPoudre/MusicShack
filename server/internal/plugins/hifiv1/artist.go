package hifiv1

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/mitchellh/mapstructure"
)

type artistData struct {
	Id                 uint   `mapstructure:"id"`
	Name               string `mapstructure:"name"`
	PictureUrl         string `mapstructure:"picture"`
	PictureUrlFallback string `mapstructure:"selectedAlbumCoverFallback"`
}

type artistAlbums struct {
	Albums []struct {
		Id       uint
		Title    string
		CoverUrl string
	}
}

func (p *HifiV1) getArtistData(wg *sync.WaitGroup, artistData *artistData, errPtr *error, id string) {
	defer wg.Done()

	apiInstance, err := repository.GetApiInstanceByApi(p.Name())
	if err != nil {
		*errPtr = err
		return
	}

	resp, err := http.Get(apiInstance.Url + "/artist/?id=" + id)
	if err != nil {
		*errPtr = err
		return
	}
	defer resp.Body.Close()

	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		*errPtr = err
		return
	}

	var item map[string]any
	if err := json.Unmarshal(raw[0], &item); err != nil {
		*errPtr = err
		return
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  artistData,
		TagName: "mapstructure",
	})
	if err != nil {
		*errPtr = err
		return
	}
	if err := decoder.Decode(item); err != nil {
		*errPtr = err
		return
	}
}

func (p *HifiV1) getArtistAlbums(wg *sync.WaitGroup, artistAlbums *artistAlbums, id string) {
	defer wg.Done()

	apiInstance, err := repository.GetApiInstanceByApi(p.Name())
	if err != nil {
		return
	}

	resp, err := http.Get(apiInstance.Url + "/artist/?f=" + id)
	if err != nil {
		return
	}
	defer resp.Body.Close()

}

func (p *HifiV1) Artist(id string) (models.ArtistData, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var artistData artistData
	var err error = nil
	var artistAlbums artistAlbums

	p.getArtistData(&wg, &artistData, &err, id)
	p.getArtistAlbums(&wg, &artistAlbums, id)

	wg.Wait()
	if err != nil {
		return models.ArtistData{}, err
	}

	if artistData.PictureUrl != "" {
		artistData.PictureUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(artistData.PictureUrl, "-", "/") + "/640x640.jpg"
	} else {
		artistData.PictureUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(artistData.PictureUrlFallback, "-", "/") + "/640x640.jpg"
	}

	var normalizeArtistData models.ArtistData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeArtistData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.ArtistData{}, err
	}
	if err := decoder.Decode(artistData); err != nil {
		return models.ArtistData{}, err
	}

	return normalizeArtistData, nil
}
