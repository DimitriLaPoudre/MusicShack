package service

import (
	"context"
	"sync"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PluginService struct {
	l     *zerolog.Logger
	store *PluginStoreService
}

func NewPluginService(l *zerolog.Logger, store *PluginStoreService) PluginService {
	return PluginService{
		l:     l,
		store: store,
	}
}

func (s *PluginService) GetOriginalPlugin(ctx context.Context, url string) (model.Plugin, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ch := make(chan model.Plugin, 1)

	for _, plugin := range s.store.ListPluginsByName() {
		go func(p model.Plugin) {
			if err := p.Status(ctx, url); err == nil {
				select {
				case ch <- p:
					cancel()
				default:
				}
			}
		}(plugin)
	}

	select {
	case plugin := <-ch:
		return plugin, nil
	case <-ctx.Done():
		return nil, model.ErrInvalidUrlPlugin
	}
}

func (s *PluginService) GetSong(ctx context.Context, instances []model.Instance, id string) (model.EnrichedSong, error) {
	var song model.Song
	var err error
	var provider string
	for _, instance := range instances {
		plugin, ok := s.store.GetPluginByName(instance.Plugin)
		if !ok {
			err = model.ErrPluginNotFound
			continue
		}
		song, err = plugin.Song(ctx, instance.Url, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
	}
	if err != nil {
		return model.EnrichedSong{}, err
	}

	enrichedSong := s.enrichSong(ctx, provider, song)

	return enrichedSong, nil
}

func (s *PluginService) enrichSong(ctx context.Context, provider string, song model.Song) model.EnrichedSong {
	downloaded := false
	// repository.GetSongByUserIDByISRC(userId, song.Isrc)
	// if err == nil {
	// 	downloaded = true
	// }

	return model.EnrichedSong{
		Provider:   provider,
		Downloaded: downloaded,
		Song:       song,
	}
}

func (s *PluginService) GetAlbum(ctx context.Context, instances []model.Instance, id string) (model.EnrichedAlbum, error) {
	var album model.Album
	var err error
	var provider string
	for _, instance := range instances {
		plugin, ok := s.store.GetPluginByName(instance.Plugin)
		if !ok {
			err = model.ErrPluginNotFound
			continue
		}
		album, err = plugin.Album(ctx, instance.Url, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
	}
	if err != nil {
		return model.EnrichedAlbum{}, err
	}

	enrichedAlbum := s.enrichAlbum(ctx, provider, album)

	return enrichedAlbum, nil
}

func (s *PluginService) enrichAlbum(ctx context.Context, provider string, album model.Album) model.EnrichedAlbum {
	downloaded := true
	songs := []model.EnrichedSong{}
	if len(album.Songs) == 0 {
		downloaded = false
	}
	for _, song := range album.Songs {
		enrichedSong := s.enrichSong(ctx, provider, song)
		if enrichedSong.Downloaded == false {
			downloaded = false
		}

		songs = append(songs, enrichedSong)
	}

	return model.EnrichedAlbum{
		Provider:   provider,
		Downloaded: downloaded,
		Album:      album,
		Songs:      songs,
	}
}

func (s *PluginService) GetArtist(ctx context.Context, instances []model.Instance, id string) (model.EnrichedArtist, error) {
	var artist model.Artist
	var err error
	var provider string
	for _, instance := range instances {
		plugin, ok := s.store.GetPluginByName(instance.Plugin)
		if !ok {
			err = model.ErrPluginNotFound
			continue
		}
		artist, err = plugin.Artist(ctx, instance.Url, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
	}
	if err != nil {
		return model.EnrichedArtist{}, err
	}

	enrichedArtist := s.enrichArtist(ctx, provider, instances, artist)

	return enrichedArtist, nil
}

func (s *PluginService) enrichArtist(ctx context.Context, provider string, instances []model.Instance, artist model.Artist) model.EnrichedArtist {
	followed := uuid.UUID{}
	// if follow, err := repository.GetFollowByProviderByArtistID(data.Provider, data.Id); err == nil {
	// 	followed = follow.ID
	// }

	albums := []model.EnrichedAlbum{}
	eps := []model.EnrichedAlbum{}
	singles := []model.EnrichedAlbum{}
	var wg sync.WaitGroup
	for _, album := range artist.Albums {
		wg.Go(func() {
			enrichedAlbum, err := s.GetAlbum(ctx, instances, album.Id)
			if err != nil {
				return
			}
			albums = append(albums, enrichedAlbum)
		})
	}
	for _, ep := range artist.Ep {
		wg.Go(func() {
			enrichedAlbum, err := s.GetAlbum(ctx, instances, ep.Id)
			if err != nil {
				return
			}
			eps = append(eps, enrichedAlbum)
		})
	}
	for _, single := range artist.Albums {
		wg.Go(func() {
			enrichedAlbum, err := s.GetAlbum(ctx, instances, single.Id)
			if err != nil {
				return
			}
			singles = append(singles, enrichedAlbum)
		})
	}
	wg.Wait()

	return model.EnrichedArtist{
		Provider: provider,
		Followed: followed,
		Artist:   artist,
		Albums:   albums,
		Ep:       eps,
		Singles:  singles,
	}

}

func (s *PluginService) GetPlaylist(ctx context.Context, instances []model.Instance, id string) (model.EnrichedPlaylist, error) {
	var playlist model.Playlist
	var err error
	var provider string
	for _, instance := range instances {
		plugin, ok := s.store.GetPluginByName(instance.Plugin)
		if !ok {
			err = model.ErrPluginNotFound
			continue
		}
		playlist, err = plugin.Playlist(ctx, instance.Url, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
	}
	if err != nil {
		return model.EnrichedPlaylist{}, err
	}

	enrichedPlaylist := s.enrichPlaylist(ctx, provider, playlist)

	return enrichedPlaylist, nil
}

func (s *PluginService) enrichPlaylist(ctx context.Context, provider string, playlist model.Playlist) model.EnrichedPlaylist {
	downloaded := true
	songs := []model.EnrichedSong{}
	if len(playlist.Songs) == 0 {
		downloaded = false
	}
	for _, song := range playlist.Songs {
		enrichedSong := s.enrichSong(ctx, provider, song)
		if enrichedSong.Downloaded == false {
			downloaded = false
		}

		songs = append(songs, enrichedSong)
	}

	return model.EnrichedPlaylist{
		Provider:   provider,
		Downloaded: downloaded,
		Playlist:   playlist,
		Songs:      songs,
	}
}

func (s *PluginService) Search(ctx context.Context, instances []model.Instance, q string) (model.EnrichedSearch, error) {
	var result model.Search
	var err error
	var provider string
	for _, instance := range instances {
		plugin, ok := s.store.GetPluginByName(instance.Plugin)
		if !ok {
			err = model.ErrPluginNotFound
			continue
		}
		result, err = plugin.Search(ctx, instance.Url, q, q, q, q)
		if err == nil {
			provider = plugin.Provider()
			break
		}
	}
	if err != nil {
		return model.EnrichedSearch{}, err
	}

	songs := []model.EnrichedSong{}
	for _, song := range result.Songs {
		enrichedSong := s.enrichSong(ctx, provider, song)

		songs = append(songs, enrichedSong)
	}

	albums := []model.EnrichedAlbum{}
	for _, album := range result.Albums {
		enrichedAlbum := s.enrichAlbum(ctx, provider, album)

		albums = append(albums, enrichedAlbum)
	}

	artists := []model.EnrichedArtist{}
	for _, artist := range result.Artists {
		enrichedArtist := s.enrichArtist(ctx, provider, instances, artist)

		artists = append(artists, enrichedArtist)
	}

	playlists := []model.EnrichedPlaylist{}
	for _, playlist := range result.Playlists {
		enrichedPlaylist := s.enrichPlaylist(ctx, provider, playlist)

		playlists = append(playlists, enrichedPlaylist)
	}

	return model.EnrichedSearch{
		Songs:     songs,
		Albums:    albums,
		Artists:   artists,
		Playlists: playlists,
	}, nil
}
