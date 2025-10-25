package hifiv1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/mitchellh/mapstructure"
)

type albumData struct {
	Id            uint   `mapstructure:"id"`
	Title         string `mapstructure:"title"`
	Duration      uint   `mapstructure:"duration"`
	ReleaseDate   string `mapstructure:"releaseDate"`
	NumberTracks  uint   `mapstructure:"numberOfTracks"`
	NumberVolumes uint   `mapstructure:"numberOfVolumes"`
	Type          string `mapstructure:"type"`
	CoverUrl      string `mapstructure:"cover"`
	AudioQuality  string `mapstructure:"audioQuality"`
	Artist        struct {
		Id   uint   `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artist"`
	Artists []struct {
		Id   uint   `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artists"`
	Limit       uint `mapstructure:"limit"`
	Offset      uint `mapstructure:"offset"`
	NumberSongs uint `mapstructure:"totalNumberOfItems"`
	DirtySongs  []struct {
		SongData struct {
			Id           uint   `mapstructure:"id"`
			Title        string `mapstructure:"title"`
			Duration     uint   `mapstructure:"duration"`
			TrackNumber  uint   `mapstructure:"trackNumber"`
			VolumeNumber uint   `mapstructure:"volumeNumber"`
			Artists      []struct {
				Id   uint   `mapstructure:"id"`
				Name string `mapstructure:"name"`
			} `mapstructure:"artists"`
		} `mapstructure:"item"`
		Type string `mapstructure:"type"`
	} `mapstructure:"items"`
	Songs []struct {
		Id           uint
		Title        string
		Duration     uint
		TrackNumber  uint
		VolumeNumber uint
		Artists      []struct {
			Id   uint
			Name string
		}
	}
}

func (p *HifiV1) Album(id string) (models.AlbumData, error) {
	apiInstance, err := repository.GetApiInstanceByApi(p.Name())
	if err != nil {
		return models.AlbumData{}, err
	}
	resp, err := http.Get(apiInstance.Url + "/album/?id=" + id)
	if err != nil {
		return models.AlbumData{}, err
	}
	defer resp.Body.Close()

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return models.AlbumData{}, err
	}

	data := make(map[string]any)
	for _, item := range items {
		for key, value := range item {
			data[key] = value
		}
	}

	var albumData albumData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &albumData,
		TagName: "mapstructure",
	})
	if err != nil {
		return models.AlbumData{}, err
	}
	decoder.Decode(data)

	for _, item := range albumData.DirtySongs {
		albumData.Songs = append(albumData.Songs, struct {
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
	albumData.CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(albumData.CoverUrl, "-", "/") + "/640x640.jpg"

	var normalizeAlbumData models.AlbumData
	decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeAlbumData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.AlbumData{}, err
	}
	decoder.Decode(albumData)

	return normalizeAlbumData, nil
}
