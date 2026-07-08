package hifi

import (
	"context"
	"fmt"
	lib_url "net/url"
	"strconv"
	"sync"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	hifi_utils "github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/plugin/hifi/utils"
	"golang.org/x/time/rate"
)

func getSearchSong(ctx context.Context, limiter *rate.Limiter, url string, q string) (searchSongData, error) {
	searchSong, err := hifi_utils.FetchType[searchSongData](ctx, url+"/search/?s="+lib_url.QueryEscape(q), limiter)
	if err != nil {
		return searchSongData{}, fmt.Errorf("getSearchSong: %w", err)
	}

	return searchSong, nil
}

func getSearchAlbum(ctx context.Context, limiter *rate.Limiter, url string, q string) (searchAlbumData, error) {
	searchAlbum, err := hifi_utils.FetchType[searchAlbumData](ctx, url+"/search/?al="+lib_url.QueryEscape(q), limiter)
	if err != nil {
		return searchAlbumData{}, fmt.Errorf("getSearchAlbum: %w", err)
	}

	return searchAlbum, nil
}

func getSearchArtist(ctx context.Context, limiter *rate.Limiter, url string, q string) (searchArtistData, error) {
	searchArtist, err := hifi_utils.FetchType[searchArtistData](ctx, url+"/search/?a="+lib_url.QueryEscape(q), limiter)
	if err != nil {
		return searchArtistData{}, fmt.Errorf("getSearchArtist: %w", err)
	}

	return searchArtist, nil
}

func getSearchPlaylist(ctx context.Context, limiter *rate.Limiter, url string, q string) (searchPlaylistData, error) {
	searchPlaylist, err := hifi_utils.FetchType[searchPlaylistData](ctx, url+"/search/?a="+lib_url.QueryEscape(q), limiter)
	if err != nil {
		return searchPlaylistData{}, fmt.Errorf("getSearchPlaylist: %w", err)
	}

	return searchPlaylist, nil
}

func getSearchData(ctx context.Context, limiter *rate.Limiter, url, song, album, artist, playlist string) (searchSongData, searchAlbumData, searchArtistData, searchPlaylistData, error) {
	var searchSong searchSongData
	var searchSongErr error
	var searchAlbum searchAlbumData
	var searchAlbumErr error
	var searchArtist searchArtistData
	var searchArtistErr error
	var searchPlaylist searchPlaylistData
	var searchPlaylistErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		searchSong, searchSongErr = getSearchSong(ctx, limiter, url, song)
	})
	wg.Go(func() {
		searchAlbum, searchAlbumErr = getSearchAlbum(ctx, limiter, url, album)
	})
	wg.Go(func() {
		searchArtist, searchArtistErr = getSearchArtist(ctx, limiter, url, artist)
	})
	wg.Go(func() {
		searchPlaylist, searchPlaylistErr = getSearchPlaylist(ctx, limiter, url, playlist)
	})

	if searchSongErr != nil {
		return searchSongData{}, searchAlbumData{}, searchArtistData{}, searchPlaylistData{}, fmt.Errorf("getSearchData: %w", searchSongErr)
	}
	if searchAlbumErr != nil {
		return searchSongData{}, searchAlbumData{}, searchArtistData{}, searchPlaylistData{}, fmt.Errorf("getSearchData: %w", searchAlbumErr)
	}
	if searchArtistErr != nil {
		return searchSongData{}, searchAlbumData{}, searchArtistData{}, searchPlaylistData{}, fmt.Errorf("getSearchData: %w", searchArtistErr)
	}
	if searchPlaylistErr != nil {
		return searchSongData{}, searchAlbumData{}, searchArtistData{}, searchPlaylistData{}, fmt.Errorf("getSearchData: %w", searchPlaylistErr)
	}

	return searchSong, searchAlbum, searchArtist, searchPlaylist, nil
}

func (p *Hifi) Search(ctx context.Context, url, song, album, artist, playlist string) (model.Search, error) {
	songData, albumData, artistData, playlistData, err := getSearchData(ctx, p.limiter, url, song, album, artist, playlist)
	if err != nil {
		return model.Search{}, fmt.Errorf("Hifi.Search: %w", err)
	}

	songs := []model.Song{}
	if len(songData.Data.Songs) != 0 {
		for _, song := range songData.Data.Songs {
			audioQuality := LOW
			switch song.AudioQuality {
			case "LOW":
				audioQuality = LOW
			case "HIGH":
				audioQuality = HIGH
			case "LOSSLESS":
				audioQuality = LOSSLESS
			}
			for _, quality := range song.MediaMetadata.Tags {
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
			for _, artist := range song.Artists {
				artists = append(artists, model.Artist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
			}

			songs = append(songs,
				model.Song{
					Id:           strconv.FormatUint(uint64(song.Id), 10),
					Title:        song.Title,
					Duration:     song.Duration,
					AudioQuality: audioQuality,
					Popularity:   song.Popularity,
					Explicit:     song.Explicit,
					Isrc:         song.Isrc,
					Artists:      artists,
					Album: model.Album{
						Id:       strconv.FormatUint(uint64(song.Album.Id), 10),
						Title:    song.Album.Title,
						CoverUrl: hifi_utils.GetImageURL(song.Album.CoverUrl, 640),
					},
				})
		}
	}

	albums := []model.Album{}
	if len(albumData.Data.Albums.Albums) != 0 {
		for _, album := range albumData.Data.Albums.Albums {
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

			albums = append(albums,
				model.Album{
					Id:           strconv.FormatUint(uint64(album.Id), 10),
					Title:        album.Title,
					Duration:     album.Duration,
					CoverUrl:     hifi_utils.GetImageURL(album.CoverUrl, 640),
					AudioQuality: audioQuality,
					Explicit:     album.Explicit,
					Popularity:   album.Popularity,
					Artists:      artists,
				})
		}
	}

	artists := []model.Artist{}
	if len(artistData.Data.Artists.Artists) != 0 {
		for _, artist := range artistData.Data.Artists.Artists {
			artists = append(artists,
				model.Artist{
					Id:         strconv.FormatUint(uint64(artist.Id), 10),
					Name:       artist.Name,
					PictureUrl: hifi_utils.GetImageURL(artist.PictureUrl, 750),
					Popularity: artist.Popularity,
				})
		}
	}

	playlists := []model.Playlist{}
	if len(playlistData.Data.Playlists.Playlists) != 0 {
		for _, playlist := range playlistData.Data.Playlists.Playlists {
			playlists = append(playlists,
				model.Playlist{
					Id:       playlist.UUID,
					Title:    playlist.Title,
					Duration: playlist.Duration,
					CoverURL: hifi_utils.GetImageURL(playlist.SquareImage, 640),
				})
		}
	}

	return model.Search{
		Songs:     songs,
		Albums:    albums,
		Artists:   artists,
		Playlists: playlists,
	}, nil
}
