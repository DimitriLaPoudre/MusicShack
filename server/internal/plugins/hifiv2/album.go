package hifiv2

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

	var data albumData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) Album(ctx context.Context, userId uint, id string) (models.AlbumData, error) {
	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", err)
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
			return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", errors.New("can't fetch"))
		} else {
			data = find
		}
	case <-ctx.Done():
		cancel()
		return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", context.Canceled)
	}

	var normalizeAlbumData models.AlbumData
	{
		if len(data.Data.Items) == 0 {
			return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", errors.New("album songs not found"))
		}

		firstSong := data.Data.Items[0].Item

		normalizeAlbumData.Id = strconv.FormatUint(uint64(firstSong.Album.Id), 10)
		normalizeAlbumData.Title = firstSong.Album.Title
		normalizeAlbumData.CoverUrl = utils.GetImageURL(firstSong.Album.CoverUrl, 640)
		normalizeAlbumData.Artists = append(normalizeAlbumData.Artists,
			models.AlbumDataArtist{
				Id:   strconv.FormatUint(uint64(firstSong.Artist.Id), 10),
				Name: firstSong.Artist.Name,
			})

		songs := make([]models.AlbumDataSong, 0)
		for _, item := range data.Data.Items {
			song := item.Item

			normalizeAlbumData.Duration += song.Duration
			if normalizeAlbumData.ReleaseDate < song.ReleaseDate {
				normalizeAlbumData.ReleaseDate = song.ReleaseDate
			}
			normalizeAlbumData.NumberTracks++
			if normalizeAlbumData.NumberVolumes < song.VolumeNumber {
				normalizeAlbumData.NumberVolumes = song.VolumeNumber
			}
			if normalizeAlbumData.MaximalAudioQuality > song.AudioQuality {
				normalizeAlbumData.MaximalAudioQuality = song.AudioQuality
			}

			artists := make([]models.SongDataArtist, 0)
			for _, artist := range song.Artists {
				artists = append(artists, models.SongDataArtist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
			}

			songs = append(songs, models.AlbumDataSong{
				Id:                  strconv.FormatUint(uint64(song.Id), 10),
				Title:               song.Title,
				Duration:            song.Duration,
				TrackNumber:         song.TrackNumber,
				VolumeNumber:        song.VolumeNumber,
				MaximalAudioQuality: song.AudioQuality,
				Artists:             artists,
			})
		}
		normalizeAlbumData.Songs = songs
	}

	return normalizeAlbumData, nil
}
