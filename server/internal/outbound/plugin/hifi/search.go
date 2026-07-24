package hifi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getSearchSong(ctx context.Context, urls []string, q string) (searchSongResponse, error) {
	searchSong, err := hifi_utils.CachedMultiFetchTyped[searchSongResponse](ctx, urls, "/search/?s="+url.QueryEscape(q), &p.limiters, &p.cache)
	if err != nil {
		return searchSongResponse{}, fmt.Errorf("fetch search song with url list: %w", err)
	}

	return searchSong, nil
}

func (p *Hifi) getSearchAlbum(ctx context.Context, urls []string, q string) (searchAlbumResponse, error) {
	searchAlbum, err := hifi_utils.CachedMultiFetchTyped[searchAlbumResponse](ctx, urls, "/search/?al="+url.QueryEscape(q), &p.limiters, &p.cache)
	if err != nil {
		return searchAlbumResponse{}, fmt.Errorf("fetch search album with url list: %w", err)
	}

	return searchAlbum, nil
}

func (p *Hifi) getSearchArtist(ctx context.Context, urls []string, q string) (searchArtistResponse, error) {
	searchArtist, err := hifi_utils.CachedMultiFetchTyped[searchArtistResponse](ctx, urls, "/search/?a="+url.QueryEscape(q), &p.limiters, &p.cache)
	if err != nil {
		return searchArtistResponse{}, fmt.Errorf("fetch search artist with url list: %w", err)
	}

	return searchArtist, nil
}

func (p *Hifi) getSearchPlaylist(ctx context.Context, urls []string, q string) (searchPlaylistResponse, error) {
	searchPlaylist, err := hifi_utils.CachedMultiFetchTyped[searchPlaylistResponse](ctx, urls, "/search/?p="+url.QueryEscape(q), &p.limiters, &p.cache)
	if err != nil {
		return searchPlaylistResponse{}, fmt.Errorf("fetch search playlist with url list: %w", err)
	}

	return searchPlaylist, nil
}

func (p *Hifi) getSearch(ctx context.Context, urls []string, song, album, artist, playlist string) (searchSongResponse, searchAlbumResponse, searchArtistResponse, searchPlaylistResponse, error) {
	var searchSong searchSongResponse
	var searchSongErr error
	var searchAlbum searchAlbumResponse
	var searchAlbumErr error
	var searchArtist searchArtistResponse
	var searchArtistErr error
	var searchPlaylist searchPlaylistResponse
	var searchPlaylistErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		searchSong, searchSongErr = p.getSearchSong(ctx, urls, song)
	})
	wg.Go(func() {
		searchAlbum, searchAlbumErr = p.getSearchAlbum(ctx, urls, album)
	})
	wg.Go(func() {
		searchArtist, searchArtistErr = p.getSearchArtist(ctx, urls, artist)
	})
	wg.Go(func() {
		searchPlaylist, searchPlaylistErr = p.getSearchPlaylist(ctx, urls, playlist)
	})
	wg.Wait()

	if searchSongErr != nil {
		return searchSongResponse{}, searchAlbumResponse{}, searchArtistResponse{}, searchPlaylistResponse{}, searchSongErr
	}
	if searchAlbumErr != nil {
		return searchSongResponse{}, searchAlbumResponse{}, searchArtistResponse{}, searchPlaylistResponse{}, searchAlbumErr
	}
	if searchArtistErr != nil {
		return searchSongResponse{}, searchAlbumResponse{}, searchArtistResponse{}, searchPlaylistResponse{}, searchArtistErr
	}
	if searchPlaylistErr != nil {
		return searchSongResponse{}, searchAlbumResponse{}, searchArtistResponse{}, searchPlaylistResponse{}, searchPlaylistErr
	}

	return searchSong, searchAlbum, searchArtist, searchPlaylist, nil
}

func (p *Hifi) Search(ctx context.Context, instances []model.Instance, song, album, artist, playlist string) (model.Search, error) {
	urls := hifi_utils.InstancesToUrls(instances)

	songData, albumData, artistData, playlistData, err := p.getSearch(ctx, urls, song, album, artist, playlist)
	if err != nil {
		return model.Search{}, err
	}

	songs := []model.Song{}
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
				ID:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		songs = append(songs,
			model.Song{
				ID:           strconv.FormatUint(uint64(song.ID), 10),
				Title:        song.Title,
				Duration:     song.Duration,
				AudioQuality: audioQuality,
				Popularity:   song.Popularity,
				Explicit:     song.Explicit,
				Isrc:         song.Isrc,
				Artists:      artists,
				Album: model.Album{
					ID:       strconv.FormatUint(uint64(song.Album.ID), 10),
					Title:    song.Album.Title,
					CoverUrl: hifi_utils.GetImageURL(song.Album.CoverUrl, 640),
				},
			})
	}

	albums := []model.Album{}
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
				ID:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		albums = append(albums,
			model.Album{
				ID:           strconv.FormatUint(uint64(album.ID), 10),
				Title:        album.Title,
				Duration:     album.Duration,
				CoverUrl:     hifi_utils.GetImageURL(album.CoverUrl, 640),
				AudioQuality: audioQuality,
				Explicit:     album.Explicit,
				Popularity:   album.Popularity,
				Artists:      artists,
			})
	}

	artists := []model.Artist{}
	for _, artist := range artistData.Data.Artists.Artists {
		artists = append(artists,
			model.Artist{
				ID:         strconv.FormatUint(uint64(artist.ID), 10),
				Name:       artist.Name,
				PictureUrl: hifi_utils.GetImageURL(artist.PictureUrl, 750),
				Popularity: artist.Popularity,
			})
	}

	playlists := []model.Playlist{}
	for _, playlist := range playlistData.Data.Playlists.Playlists {
		playlists = append(playlists,
			model.Playlist{
				ID:       playlist.UUID,
				Title:    playlist.Title,
				Duration: playlist.Duration,
				CoverURL: hifi_utils.GetImageURL(playlist.SquareImage, 640),
			})
	}

	return model.Search{
		Songs:     songs,
		Albums:    albums,
		Artists:   artists,
		Playlists: playlists,
	}, nil
}
