package hifi

import (
	"context"
	"fmt"
	lib_url "net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	hifi_utils "github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/plugin/hifi/utils"
	"golang.org/x/time/rate"
)

func getArtistInfo(ctx context.Context, limiter *rate.Limiter, url string, id string) (artistInfo, error) {
	info, err := hifi_utils.FetchType[artistInfo](ctx, url+"/artist/?id="+lib_url.QueryEscape(id), limiter)
	if err != nil {
		return artistInfo{}, fmt.Errorf("getArtistInfo: %w", err)
	}

	return info, nil
}

func getArtistAlbums(ctx context.Context, limiter *rate.Limiter, url string, id string) (artistAlbums, error) {
	albums, err := hifi_utils.FetchType[artistAlbums](ctx, url+"/artist/?f="+lib_url.QueryEscape(id)+"&skip_tracks=1", limiter)
	if err != nil {
		return artistAlbums{}, fmt.Errorf("getArtistAlbums: %w", err)
	}

	return albums, nil
}

func getArtist(ctx context.Context, limiter *rate.Limiter, url string, id string) (artistInfo, artistAlbums, error) {
	var artist artistInfo
	var artistErr error
	var albums artistAlbums
	var albumsErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		defer wg.Done()
		artist, artistErr = getArtistInfo(ctx, limiter, url, id)
	})
	wg.Go(func() {
		defer wg.Done()
		albums, albumsErr = getArtistAlbums(ctx, limiter, url, id)
	})
	wg.Wait()

	if artistErr != nil {
		return artistInfo{}, artistAlbums{}, fmt.Errorf("getArtistData: %w", artistErr)
	}
	if albumsErr != nil {
		return artistInfo{}, artistAlbums{}, fmt.Errorf("getArtistData: %w", albumsErr)
	}

	return artist, albums, nil
}

func (p *Hifi) Artist(ctx context.Context, url string, id string) (model.Artist, error) {
	artistInfo, artistAlbums, err := getArtist(ctx, p.limiter, url, id)
	if err != nil {
		return model.Artist{}, fmt.Errorf("Hifi.Artist: %w", err)
	}

	pictureURL := artistInfo.Artist.PictureUrl
	if pictureURL == "" {
		pictureURL = artistInfo.Artist.PictureUrlFallback
	}
	pictureURL = hifi_utils.GetImageURL(pictureURL, 750)

	type albumItemComparaison struct {
		Title       string
		ReleaseDate string
		TrackNumber uint
	}

	best := make(map[albumItemComparaison]albumItem)
	for _, album := range artistAlbums.Albums.Items {
		if bestVersion, ok := best[albumItemComparaison{
			Title:       strings.ToLower(album.Title),
			ReleaseDate: album.ReleaseDate,
			TrackNumber: album.NumberOfTracks,
		}]; !ok || (!bestVersion.Explicit && album.Explicit) || (bestVersion.Explicit == album.Explicit && len(bestVersion.MediaMetadata.Tags) < len(album.MediaMetadata.Tags)) {
			best[albumItemComparaison{
				Title:       strings.ToLower(album.Title),
				ReleaseDate: album.ReleaseDate,
				TrackNumber: album.NumberOfTracks,
			}] = album
		}
	}

	// type albumItemComparaisonExtension struct {
	// 	Title       string
	// 	ReleaseDate string
	// }

	// extension := make(map[albumItemComparaisonExtension][]*albumItem)
	// for _, album := range best {
	// 	newExtension := []*albumItem{}
	// 	bestExtension, ok := extension[albumItemComparaisonExtension{
	// 		Title:       strings.ToLower(album.Title),
	// 		ReleaseDate: album.ReleaseDate,
	// 	}]
	// 	if !ok {
	// 		extension[albumItemComparaisonExtension{
	// 			Title:       strings.ToLower(album.Title),
	// 			ReleaseDate: album.ReleaseDate,
	// 		}] = append(newExtension, album)
	// 		continue
	// 	}
	//
	// 	for i, tmp := range bestExtension {
	//
	// 	}
	// }

	list := []albumItem{}
	for _, album := range best {
		list = append(list, album)
	}

	slices.SortFunc(list, func(a, b albumItem) int {
		if a.ReleaseDate > b.ReleaseDate {
			return -1
		}
		if a.ReleaseDate < b.ReleaseDate {
			return 1
		}
		if a.Id > b.Id {
			return -1
		}
		if a.Id < b.Id {
			return 1
		}
		return 0
	})

	albums := []model.Album{}
	eps := []model.Album{}
	singles := []model.Album{}
	for _, album := range list {
		releaseDate, err := time.Parse(StreamStartDateLayout, album.ReleaseDate)
		if err != nil {
			p.l.Warn().Msg(fmt.Sprintf("Hifi.Artist: for album: time.Parse(%s): %s", album.ReleaseDate, err.Error()))
		}

		audioQuality := LOW
		switch album.AudioQuality {
		case "LOW":
			audioQuality = LOW
		case "HIGH":
			audioQuality = HIGH
		case "LOSSLESS":
			audioQuality = LOSSLESS
		}
		for _, quality := range album.MediaMetadata.Tags {
			switch quality {
			case "HIRES_LOSSLESS":
				audioQuality = HIRES
			case "LOSSLESS", "DOLBY_ATMOS":
				if audioQuality != HIRES {
					audioQuality = LOSSLESS
				}
			}
		}

		artists := []model.Artist{}
		for _, artist := range album.Artists {
			artists = append(artists, model.Artist{
				Id:   strconv.FormatUint(uint64(artist.Id), 10),
				Name: artist.Name,
			})
		}

		newAlbum := model.Album{
			Id:           strconv.FormatUint(uint64(album.Id), 10),
			Title:        album.Title,
			Duration:     album.Duration,
			ReleaseDate:  releaseDate,
			CoverUrl:     hifi_utils.GetImageURL(album.CoverUrl, 1280),
			AudioQuality: audioQuality,
			Explicit:     album.Explicit,
			Artists:      artists,
		}

		switch album.Type {
		case "ALBUM":
			albums = append(albums, newAlbum)
		case "EP":
			eps = append(eps, newAlbum)
		case "SINGLE":
			singles = append(singles, newAlbum)
		}
	}

	return model.Artist{
		Id:      strconv.FormatUint(uint64(artistInfo.Artist.Id), 10),
		Name:    artistInfo.Artist.Name,
		Albums:  albums,
		Ep:      eps,
		Singles: singles,
	}, nil
}
