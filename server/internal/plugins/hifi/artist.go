package hifi

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
	Id    string
	Title string
	Rows  []struct {
		Modules []struct {
			PagedList struct {
				Limit               uint
				Offset              uint
				TotalNumbersOfItems uint
				Items               []struct {
					Id       uint
					Title    string
					CoverUrl string `mapstructure:"cover"`
					Artists  []struct {
						Id   uint
						Name string
					}
				}
			}
		}
	}
}

func (p *Hifi) getArtistData(wg *sync.WaitGroup, artistData *artistData, errPtr *error, id string) {
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

func (p *Hifi) getArtistAlbums(wg *sync.WaitGroup, artistAlbums *artistAlbums, id string) {
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

	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return
	}
	var item map[string]any
	if err := json.Unmarshal(raw[0], &item); err != nil {
		return
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  artistAlbums,
		TagName: "mapstructure",
	})
	if err != nil {
		return
	}
	if err := decoder.Decode(item); err != nil {
		return
	}
}

func (p *Hifi) Artist(id string) (models.ArtistData, error) {
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

	if artistAlbums.Id != "" {
		items := artistAlbums.Rows[0].Modules[0].PagedList.Items
		for i := range items {
			items[i].CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(items[i].CoverUrl, "-", "/") + "/640x640.jpg"
		}
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

	decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeArtistData.Albums,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.ArtistData{}, err
	}
	if err := decoder.Decode(artistAlbums.Rows[0].Modules[0].PagedList.Items); err != nil {
		return models.ArtistData{}, err
	}

	return normalizeArtistData, nil
}
