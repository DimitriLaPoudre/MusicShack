package hifiv2

import (
	"context"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"strings"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/mitchellh/mapstructure"
)

func getAlbum(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- albumData, id string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/album/?id="+id, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return
	}

	tmp := make(map[string]any)
	for _, item := range items {
		maps.Copy(tmp, item)
	}

	var data albumData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &data,
		TagName: "mapstructure",
	})
	if err != nil {
		return
	}
	if err := decoder.Decode(tmp); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) Album(ctx context.Context, id string) (models.AlbumData, error) {
	apiInstances, err := repository.ListApiInstancesByApi(p.Name())
	if err != nil {
		return models.AlbumData{}, err
	}

	var data albumData
	routineCtx, cancel := context.WithCancel(context.Background())
	ch := make(chan albumData)
	var wg sync.WaitGroup

	wg.Add(len(apiInstances))
	go func() {
		wg.Wait()
		close(ch)
	}()
	for _, instance := range apiInstances {
		go getAlbum(routineCtx, &wg, instance.Url, ch, id)
	}
	select {
	case find, ok := <-ch:
		cancel()
		if !ok {
			return models.AlbumData{}, errors.New("Album not found")
		} else {
			data = find
		}
	case <-ctx.Done():
		cancel()
	}

	for _, item := range data.DirtySongs {
		data.Songs = append(data.Songs, struct {
			Id           uint
			Title        string
			Duration     uint
			TrackNumber  uint
			VolumeNumber uint
			Artists      []struct {
				Id   uint
				Name string
			}
		}(
			item.SongData,
		))
	}
	if data.CoverUrl != "" {
		data.CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(data.CoverUrl, "-", "/") + "/640x640.jpg"
	}

	var normalizeAlbumData models.AlbumData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeAlbumData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.AlbumData{}, err
	}
	if err := decoder.Decode(data); err != nil {
		return models.AlbumData{}, err
	}

	return normalizeAlbumData, nil
}
