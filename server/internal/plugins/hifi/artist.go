package hifi

import (
	"context"
	"encoding/json"
	"errors"
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

func (p *Hifi) getArtistData(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- artistData, id string) {
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

	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return
	}
	var item map[string]any
	if err := json.Unmarshal(raw[0], &item); err != nil {
		return
	}

	var data artistData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &data,
		TagName: "mapstructure",
	})
	if err != nil {
		return
	}
	if err := decoder.Decode(item); err != nil {
		return
	}

	ch <- data
}

func (p *Hifi) getArtistAlbums(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- artistAlbums, id string) {
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

	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return
	}
	var item map[string]any
	if err := json.Unmarshal(raw[0], &item); err != nil {
		return
	}

	var data artistAlbums
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &data,
		TagName: "mapstructure",
	})
	if err != nil {
		return
	}
	if err := decoder.Decode(item); err != nil {
		return
	}

	ch <- data
}

func (p *Hifi) Artist(ctx context.Context, id string) (models.ArtistData, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var data artistData
	var albums artistAlbums

	apiInstances, err := repository.ListApiInstancesByApi(p.Name())
	if err != nil {
		return models.ArtistData{}, err
	}
	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan artistData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(apiInstances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range apiInstances {
			go p.getArtistData(routineCtx, &wgRoutine, instance.Url, ch, id)
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

		wgRoutine.Add(len(apiInstances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range apiInstances {
			go p.getArtistAlbums(routineCtx, &wgRoutine, instance.Url, ch, id)
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
		return models.ArtistData{}, errors.New("connection closed")
	default:
	}

	var normalizeArtistData models.ArtistData

	if data.Id == 0 {
		return models.ArtistData{}, errors.New("Artist can't be fetch")
	}

	if data.PictureUrl != "" {
		data.PictureUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(data.PictureUrl, "-", "/") + "/160x160.jpg"
	} else if data.PictureUrlFallback != "" {
		data.PictureUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(data.PictureUrlFallback, "-", "/") + "/160x160.jpg"
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeArtistData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.ArtistData{}, err
	}
	if err := decoder.Decode(data); err != nil {
		return models.ArtistData{}, err
	}

	if albums.Id != "" {
		for i := range albums.Rows[0].Modules[0].PagedList.Items {
			if albums.Rows[0].Modules[0].PagedList.Items[i].CoverUrl != "" {
				albums.Rows[0].Modules[0].PagedList.Items[i].CoverUrl =
					"https://resources.tidal.com/images/" + strings.ReplaceAll(albums.Rows[0].Modules[0].PagedList.Items[i].CoverUrl, "-", "/") + "/640x640.jpg"
			}
		}

		decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           &normalizeArtistData.Albums,
			TagName:          "useless",
			WeaklyTypedInput: true,
		})
		decoder.Decode(albums.Rows[0].Modules[0].PagedList.Items)
	}

	return normalizeArtistData, nil
}
