package hifiv2_2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
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

	if resp.StatusCode >= 400 {
		return
	}

	var data albumData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *Hifi) Album(ctx context.Context, userId uint, id string) (models.AlbumData, error) {
	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.AlbumData{}, fmt.Errorf("Hifi.Album: %w", err)
	}

	routineCtx, cancel := context.WithCancel(context.Background())
	ch := make(chan albumData)
	var wg sync.WaitGroup

	wg.Add(len(instances))
	go func() {
		wg.Wait()
		close(ch)
	}()
	for _, instance := range instances {
		go getAlbum(routineCtx, &wg, instance.Url, ch, id)
	}

	var data albumData
	select {
	case find, ok := <-ch:
		cancel()
		if !ok {
			return models.AlbumData{}, fmt.Errorf("Hifi.Album: %w", errors.New("can't fetch"))
		} else {
			data = find
		}
	case <-ctx.Done():
		cancel()
		return models.AlbumData{}, fmt.Errorf("Hifi.Album: %w", context.Canceled)
	}

	normalizeAlbumData := models.AlbumData{
		Api:           p.Name(),
		Id:            strconv.FormatUint(uint64(data.Data.Id), 10),
		Title:         data.Data.Title,
		Duration:      data.Data.Duration,
		ReleaseDate:   data.Data.ReleaseDate,
		NumberTracks:  data.Data.NumberOfTracks,
		NumberVolumes: data.Data.NumberOfVolumes,
		CoverUrl:      utils.GetImageURL(data.Data.CoverUrl, 640),
		Songs:         make([]models.AlbumDataSong, 0),
	}
	{
		for _, quality := range data.Data.MediaMetadata.Tags {
			switch quality {
			case "HIRES_LOSSLESS":
				normalizeAlbumData.AudioQuality = max(normalizeAlbumData.AudioQuality, models.QualityHiresLossless)
			case "LOSSLESS", "DOLBY_ATMOS":
				normalizeAlbumData.AudioQuality = max(normalizeAlbumData.AudioQuality, models.QualityLossless)
			}
		}

		for _, artist := range data.Data.Artists {
			normalizeAlbumData.Artists = append(normalizeAlbumData.Artists,
				models.AlbumDataArtist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
		}

		if len(data.Data.Items) == 0 {
			return models.AlbumData{}, fmt.Errorf("Hifi.Album: %w", errors.New("album songs not found"))
		}

		for _, item := range data.Data.Items {
			song := item.Item

			audioQuality := models.QualityHigh
			for _, quality := range song.MediaMetadata.Tags {
				switch quality {
				case "HIRES_LOSSLESS":
					audioQuality = max(audioQuality, models.QualityHiresLossless)
				case "LOSSLESS", "DOLBY_ATMOS":
					audioQuality = max(audioQuality, models.QualityLossless)
				}
			}

			artists := make([]models.SongDataArtist, 0)
			for _, artist := range song.Artists {
				artists = append(artists, models.SongDataArtist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
			}

			normalizeAlbumData.Songs = append(normalizeAlbumData.Songs, models.AlbumDataSong{
				Id:           strconv.FormatUint(uint64(song.Id), 10),
				Title:        song.Title,
				Duration:     song.Duration,
				TrackNumber:  song.TrackNumber,
				VolumeNumber: song.VolumeNumber,
				AudioQuality: audioQuality,
				Artists:      artists,
			})
		}
	}

	return normalizeAlbumData, nil
}
