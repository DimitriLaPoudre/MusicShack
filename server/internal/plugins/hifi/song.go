package hifi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/mitchellh/mapstructure"
)

type SongData struct {
	Id           string `mapstructure:"id"`
	Title        string `mapstructure:"title"`
	Duration     uint   `mapstructure:"duration"`
	ReleaseDate  string `mapstructure:"streamStartDate"`
	TrackNumber  uint   `mapstructure:"trackNumber"`
	VolumeNumber uint   `mapstructure:"volumeNumber"`
	AudioQuality string `mapstructure:"audioQuality"`
	Artist       struct {
		Id   string `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artist"`
	Artists []struct {
		Id   string `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artists"`
	Album struct {
		Id    string `mapstructure:"id"`
		Title string `mapstructure:"title"`
	} `mapstructure:"album"`
	BitDepth    uint   `mapstructure:"bitDepth"`
	SampleRate  uint   `mapstructure:"sampleRate"`
	DownloadUrl string `mapstructure:"OriginalTrackUrl"`
}

func (p *Hifi) Song(id string) (models.SongData, error) {
	apiInstance, err := repository.GetApiInstanceByApi("hifi")
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
		log.Println("json decode")
		return models.SongData{}, err
	}

	data := make(map[string]any)
	for _, item := range items {
		for key, value := range item {
			data[key] = value
		}
	}

	var songData SongData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &songData,
		TagName: "mapstructure",
	})
	if err != nil {
		log.Println("decode into tmp songdata")
		return models.SongData{}, err
	}
	decoder.Decode(data)

	var normalizeSongData SongData
	decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeSongData,
		TagName:          "",
		WeaklyTypedInput: true,
	})
	if err != nil {
		log.Println("decode into songdata")
		return models.SongData{}, err
	}
	decoder.Decode(data)

	log.Println("finish")
	return models.SongData(normalizeSongData), nil
}
