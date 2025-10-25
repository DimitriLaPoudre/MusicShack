package hifiv1

import (
	"encoding/json"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/mitchellh/mapstructure"
)

type songData struct {
	Id           uint   `mapstructure:"id"`
	Title        string `mapstructure:"title"`
	Duration     uint   `mapstructure:"duration"`
	ReleaseDate  string `mapstructure:"streamStartDate"`
	TrackNumber  uint   `mapstructure:"trackNumber"`
	VolumeNumber uint   `mapstructure:"volumeNumber"`
	AudioQuality string `mapstructure:"audioQuality"`
	Artist       struct {
		Id   uint   `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artist"`
	Artists []struct {
		Id   uint   `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artists"`
	Album struct {
		Id    uint   `mapstructure:"id"`
		Title string `mapstructure:"title"`
	} `mapstructure:"album"`
	BitDepth    uint   `mapstructure:"bitDepth"`
	SampleRate  uint   `mapstructure:"sampleRate"`
	DownloadUrl string `mapstructure:"OriginalTrackUrl"`
}

func (p *HifiV1) Song(id string) (models.SongData, error) {
	apiInstance, err := repository.GetApiInstanceByApi(p.Name())
	if err != nil {
		return models.SongData{}, err
	}
	resp, err := http.Get(apiInstance.Url + "/track/?id=" + id)
	if err != nil {
		return models.SongData{}, err
	}
	defer resp.Body.Close()

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return models.SongData{}, err
	}

	data := make(map[string]any)
	for _, item := range items {
		for key, value := range item {
			data[key] = value
		}
	}

	var songData songData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &songData,
		TagName: "mapstructure",
	})
	if err != nil {
		return models.SongData{}, err
	}
	if err := decoder.Decode(data); err != nil {
		return models.SongData{}, err
	}

	var normalizeSongData models.SongData
	decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeSongData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.SongData{}, err
	}
	if err := decoder.Decode(songData); err != nil {
		return models.SongData{}, err
	}

	return normalizeSongData, nil
}
