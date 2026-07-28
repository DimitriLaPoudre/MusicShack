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

func (p *Hifi) getSearchSong(ctx context.Context, urls []string, q string, limit int, offset int) (searchSongResponse, error) {
	searchSong, err := hifi_utils.MultiFetchTyped[searchSongResponse](ctx, urls, fmt.Sprintf("/search/?s=%s&limit=%d&offset=%d", url.QueryEscape(q), limit, offset), &p.limiters)
	if err != nil {
		return searchSongResponse{}, fmt.Errorf("fetch search song with url list: %w", err)
	}

	return searchSong, nil
}

func (p *Hifi) getSearchAlbum(ctx context.Context, urls []string, q string, limit int, offset int) (searchAlbumResponse, error) {
	searchAlbum, err := hifi_utils.MultiFetchTyped[searchAlbumResponse](ctx, urls, fmt.Sprintf("/search/?al=%s&limit=%d&offset=%d", url.QueryEscape(q), limit, offset), &p.limiters)
	if err != nil {
		return searchAlbumResponse{}, fmt.Errorf("fetch search album with url list: %w", err)
	}

	return searchAlbum, nil
}

func (p *Hifi) getSearchArtist(ctx context.Context, urls []string, q string, limit int, offset int) (searchArtistResponse, error) {
	searchArtist, err := hifi_utils.MultiFetchTyped[searchArtistResponse](ctx, urls, fmt.Sprintf("/search/?a=%s&limit=%d&offset=%d", url.QueryEscape(q), limit, offset), &p.limiters)
	if err != nil {
		return searchArtistResponse{}, fmt.Errorf("fetch search artist with url list: %w", err)
	}

	return searchArtist, nil
}

func (p *Hifi) getSearchPlaylist(ctx context.Context, urls []string, q string, limit int, offset int) (searchPlaylistResponse, error) {
	searchPlaylist, err := hifi_utils.MultiFetchTyped[searchPlaylistResponse](ctx, urls, fmt.Sprintf("/search/?p=%s&limit=%d&offset=%d", url.QueryEscape(q), limit, offset), &p.limiters)
	if err != nil {
		return searchPlaylistResponse{}, fmt.Errorf("fetch search playlist with url list: %w", err)
	}

	return searchPlaylist, nil
}

func (p *Hifi) getSearch(ctx context.Context, urls []string, song, album, artist, playlist string, limit int, offset int) (searchSongResponse, searchAlbumResponse, searchArtistResponse, searchPlaylistResponse, error) {
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
		searchSong, searchSongErr = p.getSearchSong(ctx, urls, song, limit, offset)
	})
	wg.Go(func() {
		searchAlbum, searchAlbumErr = p.getSearchAlbum(ctx, urls, album, limit, offset)
	})
	wg.Go(func() {
		searchArtist, searchArtistErr = p.getSearchArtist(ctx, urls, artist, limit, offset)
	})
	wg.Go(func() {
		searchPlaylist, searchPlaylistErr = p.getSearchPlaylist(ctx, urls, playlist, limit, offset)
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

func (p *Hifi) Search(ctx context.Context, instances []model.Instance, song, album, artist, playlist string, limit int, offset int) (model.Search, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	songData, albumData, artistData, playlistData, err := p.getSearch(ctx, urls, song, album, artist, playlist, limit, offset)
	if err != nil {
		return model.Search{}, err
	}

	songs := []model.SongInfo{}
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

		artists := []model.ArtistInfo{}
		for _, artist := range song.Artists {
			artists = append(artists, model.ArtistInfo{
				ID:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		songs = append(songs,
			model.SongInfo{
				ID:           strconv.FormatUint(uint64(song.ID), 10),
				Title:        song.Title,
				Duration:     song.Duration,
				AudioQuality: audioQuality,
				Popularity:   song.Popularity,
				Explicit:     song.Explicit,
				Isrc:         song.ISRC,
				Artists:      artists,
				Album: model.AlbumInfo{
					ID:       strconv.FormatUint(uint64(song.Album.ID), 10),
					Title:    song.Album.Title,
					CoverURL: hifi_utils.GetImageURL(song.Album.Cover, 640),
				},
			})
	}

	albums := []model.AlbumInfo{}
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

		artists := []model.ArtistInfo{}
		for _, artist := range album.Artists {
			artists = append(artists, model.ArtistInfo{
				ID:   strconv.FormatUint(uint64(artist.ID), 10),
				Name: artist.Name,
			})
		}

		albums = append(albums,
			model.AlbumInfo{
				ID:           strconv.FormatUint(uint64(album.ID), 10),
				Title:        album.Title,
				Duration:     album.Duration,
				CoverURL:     hifi_utils.GetImageURL(album.Cover, 640),
				AudioQuality: audioQuality,
				Explicit:     album.Explicit,
				Popularity:   album.Popularity,
				Artists:      artists,
			})
	}

	artists := []model.ArtistInfo{}
	for _, artist := range artistData.Data.Artists.Artists {
		artists = append(artists,
			model.ArtistInfo{
				ID:         strconv.FormatUint(uint64(artist.ID), 10),
				Name:       artist.Name,
				PictureURL: hifi_utils.GetImageURL(artist.Picture, 750),
				Popularity: artist.Popularity,
			})
	}

	playlists := []model.PlaylistInfo{}
	for _, playlist := range playlistData.Data.Playlists.Playlists {
		playlists = append(playlists,
			model.PlaylistInfo{
				ID:       playlist.UUID,
				Title:    playlist.Title,
				Duration: playlist.Duration,
				CoverURL: hifi_utils.GetImageURL(playlist.SquareImage, 640),
			})
	}

	return model.Search{
		Songs:     model.PaginatedSongs{Pagination: model.Pagination{Limit: songData.Data.Limit, Offset: songData.Data.Offset, TotalNumberOfItems: songData.Data.TotalNumberOfItems}, Songs: songs},
		Albums:    model.PaginatedAlbums{Pagination: model.Pagination{Limit: albumData.Data.Albums.Limit, Offset: albumData.Data.Albums.Offset, TotalNumberOfItems: albumData.Data.Albums.TotalNumberOfItems}, Albums: albums},
		Artists:   model.PaginatedArtists{Pagination: model.Pagination{Limit: artistData.Data.Artists.Limit, Offset: artistData.Data.Artists.Offset, TotalNumberOfItems: artistData.Data.Artists.TotalNumberOfItems}, Artists: artists},
		Playlists: model.PaginatedPlaylists{Pagination: model.Pagination{Limit: playlistData.Data.Playlists.Limit, Offset: playlistData.Data.Playlists.Offset, TotalNumberOfItems: playlistData.Data.Playlists.TotalNumberOfItems}, Playlists: playlists},
	}, nil
}
