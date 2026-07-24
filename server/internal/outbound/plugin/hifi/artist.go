package hifi

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getArtistInfo(ctx context.Context, urls []string, id string) (artistResponse, error) {
	info, err := hifi_utils.CachedMultiFetchTyped[artistResponse](ctx, urls, "/artist/?id="+url.QueryEscape(id), &p.limiters, &p.cache)
	if err != nil {
		return artistResponse{}, fmt.Errorf("fetch artist info with url list: %w", err)
	}

	return info, nil
}

func (p *Hifi) getArtistAlbums(ctx context.Context, urls []string, id string) (artistAlbumsResponse, error) {
	albums, err := hifi_utils.CachedMultiFetchTyped[artistAlbumsResponse](ctx, urls, "/artist/?f="+url.QueryEscape(id)+"&skip_tracks=1", &p.limiters, &p.cache)
	if err != nil {
		return artistAlbumsResponse{}, fmt.Errorf("fetch artist albums with url list: %w", err)
	}

	return albums, nil
}

func (p *Hifi) getArtist(ctx context.Context, urls []string, id string) (artistResponse, artistAlbumsResponse, error) {
	var artist artistResponse
	var artistErr error
	var albums artistAlbumsResponse
	var albumsErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		artist, artistErr = p.getArtistInfo(ctx, urls, id)
	})
	wg.Go(func() {
		albums, albumsErr = p.getArtistAlbums(ctx, urls, id)
	})
	wg.Wait()

	if artistErr != nil {
		return artistResponse{}, artistAlbumsResponse{}, artistErr
	}
	if albumsErr != nil {
		return artistResponse{}, artistAlbumsResponse{}, albumsErr
	}

	return artist, albums, nil
}

func (p *Hifi) Artist(ctx context.Context, instances []model.Instance, id string) (model.Artist, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	artistInfo, artistAlbums, err := p.getArtist(ctx, urls, id)
	if err != nil {
		return model.Artist{}, err
	}

	pictureURL := artistInfo.Artist.PictureUrl
	if pictureURL == "" {
		pictureURL = artistInfo.Artist.PictureUrlFallback
	}
	pictureURL = hifi_utils.GetImageURL(pictureURL, 750)

	// type albumItemComparaison struct {
	// 	Title       string
	// 	ReleaseDate string
	// 	TrackNumber uint
	// }

	// best := make(map[albumItemComparaison]albumItem)
	// for _, album := range artistAlbums.Albums.Items {
	// 	if bestVersion, ok := best[albumItemComparaison{
	// 		Title:       strings.ToLower(album.Title),
	// 		ReleaseDate: album.ReleaseDate,
	// 		TrackNumber: album.NumberOfTracks,
	// 	}]; !ok || (!bestVersion.Explicit && album.Explicit) || (bestVersion.Explicit == album.Explicit && len(bestVersion.MediaMetadata.Tags) < len(album.MediaMetadata.Tags)) {
	// 		best[albumItemComparaison{
	// 			Title:       strings.ToLower(album.Title),
	// 			ReleaseDate: album.ReleaseDate,
	// 			TrackNumber: album.NumberOfTracks,
	// 		}] = album
	// 	}
	// }

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

	// list := []albumItem{}
	// for _, album := range best {
	// 	list = append(list, album)
	// }

	list := artistAlbums.Albums.Items

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
		releaseDate, err := time.Parse(ReleaseDateLayout, album.ReleaseDate)
		if err != nil {
			slog.Warn(fmt.Sprintf("plugin [hifi]: failed to parse artist's album %s releaseDate %s", album.Title, album.ReleaseDate), slog.String("err", err.Error()))
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
		Id:         strconv.FormatUint(uint64(artistInfo.Artist.Id), 10),
		Name:       artistInfo.Artist.Name,
		PictureUrl: pictureURL,
		Popularity: artistInfo.Artist.Popularity,
		Albums:     albums,
		Ep:         eps,
		Singles:    singles,
	}, nil
}
