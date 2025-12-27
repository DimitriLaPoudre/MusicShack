package hifiv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func getSearchSong(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchSongData, song string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?s="+url.QueryEscape(song), nil)
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

	var data searchSongData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func getSearchAlbum(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchAlbumData, album string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?al="+url.QueryEscape(album), nil)
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

	var data searchAlbumData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func getSearchArtist(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- searchArtistData, artist string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?a="+url.QueryEscape(artist), nil)
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

	var data searchArtistData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) Search(ctx context.Context, userId uint, song, album, artist string) (models.SearchData, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	var songData searchSongData
	var albumData searchAlbumData
	var artistData searchArtistData

	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.SearchData{}, fmt.Errorf("HifiV2.Search: %w", err)
	}

	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchSongData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(instances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range instances {
			go getSearchSong(routineCtx, &wgRoutine, instance.Url, ch, song)
		}
		select {
		case find, ok := <-ch:
			cancel()
			if ok {
				songData = find
			}
		case <-ctx.Done():
			cancel()
		}

	}()
	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchAlbumData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(instances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range instances {
			go getSearchAlbum(routineCtx, &wgRoutine, instance.Url, ch, album)
		}
		select {
		case find, ok := <-ch:
			cancel()
			if ok {
				albumData = find
			}
		case <-ctx.Done():
			cancel()
		}
	}()
	go func() {
		defer wg.Done()

		routineCtx, cancel := context.WithCancel(context.Background())
		ch := make(chan searchArtistData)
		var wgRoutine sync.WaitGroup

		wgRoutine.Add(len(instances))
		go func() {
			wgRoutine.Wait()
			close(ch)
		}()
		for _, instance := range instances {
			go getSearchArtist(routineCtx, &wgRoutine, instance.Url, ch, artist)
		}
		select {
		case find, ok := <-ch:
			cancel()
			if ok {
				artistData = find
			}
		case <-ctx.Done():
			cancel()
		}
	}()

	wg.Wait()

	select {
	case <-ctx.Done():
		return models.SearchData{}, fmt.Errorf("HifiV2.Search: %w", context.Canceled)
	default:
	}

	var result models.SearchData

	if len(songData.Data.Songs) != 0 {
		songs := make([]models.SearchDataSong, 0)
		for _, song := range songData.Data.Songs {
			var audioQuality models.Quality
			switch song.AudioQuality {
			case "HI_RES_LOSSLESS":
				audioQuality = 4
			case "LOSSLESS":
				audioQuality = 3
			case "HIGH":
				audioQuality = 2
			case "LOW":
				audioQuality = 1
			}

			artists := make([]models.SongDataArtist, 0)
			for _, artist := range song.Artists {
				artists = append(artists, models.SongDataArtist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
			}

			songs = append(songs, models.SearchDataSong{
				Id:           strconv.FormatUint(uint64(song.Id), 10),
				Title:        song.Title,
				Duration:     song.Duration,
				AudioQuality: audioQuality,
				Popularity:   song.Popularity,
				Artists:      artists,
				Album: models.SongDataAlbum{
					Id:       strconv.FormatUint(uint64(song.Album.Id), 10),
					Title:    song.Album.Title,
					CoverUrl: utils.GetImageURL(song.Album.CoverUrl, 640),
				},
			})
		}
		result.Songs = songs
	}

	if len(albumData.Data.Albums.Albums) != 0 {
		albums := make([]models.SearchDataAlbum, 0)
		for _, album := range albumData.Data.Albums.Albums {
			var audioQuality models.Quality
			switch album.AudioQuality {
			case "HI_RES_LOSSLESS":
				audioQuality = 4
			case "LOSSLESS":
				audioQuality = 3
			case "HIGH":
				audioQuality = 2
			case "LOW":
				audioQuality = 1
			}

			artists := make([]models.AlbumDataArtist, 0)
			for _, artist := range album.Artists {
				artists = append(artists, models.AlbumDataArtist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
			}

			albums = append(albums, models.SearchDataAlbum{
				Id:           strconv.FormatUint(uint64(album.Id), 10),
				Title:        album.Title,
				Duration:     album.Duration,
				CoverUrl:     utils.GetImageURL(album.CoverUrl, 640),
				AudioQuality: audioQuality,
				Popularity:   album.Popularity,
				Artists:      artists,
			})
		}
		result.Albums = albums
	}

	if len(artistData.Data.Artists.Artists) != 0 {
		artists := make([]models.SearchDataArtist, 0)
		for _, artist := range artistData.Data.Artists.Artists {
			artists = append(artists, models.SearchDataArtist{
				Id:         strconv.FormatUint(uint64(artist.Id), 10),
				Name:       artist.Name,
				PictureUrl: utils.GetImageURL(artist.PictureUrl, 750),
				Popularity: artist.Popularity,
			})
		}
		result.Artists = artists
	}

	return result, nil
}
