package service

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type PluginService struct {
	store *PluginStoreService
}

func NewPluginService(store *PluginStoreService) PluginService {
	return PluginService{
		store: store,
	}
}

func (s *PluginService) InstancesToMapPluginInstances(instances []model.Instance) map[model.Plugin][]model.Instance {
	pluginInstances := map[model.Plugin][]model.Instance{}
	pluginsNotFound := map[string]struct{}{}
	for _, i := range instances {
		plugin, ok := s.store.GetPluginByName(i.Plugin)
		if !ok {
			pluginsNotFound[i.Plugin] = struct{}{}
		}
		if instances, ok := pluginInstances[plugin]; !ok {
			pluginInstances[plugin] = []model.Instance{i}
		} else {
			instances = append(instances, i)
			pluginInstances[plugin] = instances
		}
	}

	pluginsNotFoundList := []string{}
	for pluginName := range pluginsNotFound {
		pluginsNotFoundList = append(pluginsNotFoundList, pluginName)
	}
	if len(pluginsNotFoundList) > 0 {
		slog.Warn("not found plugins", slog.Any("plugins", pluginsNotFoundList))
	}

	return pluginInstances
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

var isrcRegexp = regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{3}[0-9]{2}[0-9]{5}$`)

func (s *PluginService) IsISRC(isrc string) bool {
	return isrcRegexp.MatchString(strings.ToUpper(isrc))
}

func (s *PluginService) GetSong(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string) (model.EnrichedSong, error) {
	var song model.Song
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		song, err = plugin.Song(ctx, instances, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedSong{}, fmt.Errorf("song %w: %v", model.ErrPluginDataNotFound, errMap)
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

func (s *PluginService) GetAlbum(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string) (model.EnrichedAlbum, error) {
	var album model.Album
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		album, err = plugin.Album(ctx, instances, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedAlbum{}, fmt.Errorf("album %w: %v", model.ErrPluginDataNotFound, errMap)
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

func (s *PluginService) GetArtist(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string) (model.EnrichedArtist, error) {
	var artist model.Artist
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		artist, err = plugin.Artist(ctx, instances, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedArtist{}, fmt.Errorf("artist %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedArtist := s.enrichArtist(ctx, provider, pluginInstances, artist)

	return enrichedArtist, nil
}

func (s *PluginService) enrichArtist(ctx context.Context, provider string, pluginInstances map[model.Plugin][]model.Instance, artist model.Artist) model.EnrichedArtist {
	followed := uuid.UUID{}
	// if follow, err := repository.GetFollowByProviderByArtistID(data.Provider, data.Id); err == nil {
	// 	followed = follow.ID
	// }

	albums := []model.EnrichedAlbum{}
	eps := []model.EnrichedAlbum{}
	singles := []model.EnrichedAlbum{}
	var wg sync.WaitGroup
	for _, album := range artist.Albums {
		wg.Add(1)
		go func(album model.Album) {
			defer wg.Done()
			enrichedAlbum := model.EnrichedAlbum{}
			if album.NumberTracks == 0 || len(album.Songs) == 0 || album.Songs[0].Isrc != "" {
				enrichedAlbum = s.enrichAlbum(ctx, provider, album)
			} else {
				if tmp, err := s.GetAlbum(ctx, pluginInstances, album.Id); err != nil {
					return
				} else {
					enrichedAlbum = tmp
				}
			}
			albums = append(albums, enrichedAlbum)
		}(album)
	}
	for _, ep := range artist.Ep {
		wg.Add(1)
		go func(ep model.Album) {
			defer wg.Done()
			enrichedAlbum := model.EnrichedAlbum{}
			if ep.NumberTracks == 0 || len(ep.Songs) == 0 || ep.Songs[0].Isrc != "" {
				enrichedAlbum = s.enrichAlbum(ctx, provider, ep)
			} else {
				if tmp, err := s.GetAlbum(ctx, pluginInstances, ep.Id); err != nil {
					return
				} else {
					enrichedAlbum = tmp
				}
			}
			eps = append(eps, enrichedAlbum)
		}(ep)
	}
	for _, single := range artist.Albums {
		wg.Add(1)
		go func(single model.Album) {
			defer wg.Done()
			enrichedAlbum := model.EnrichedAlbum{}
			if single.NumberTracks == 0 || len(single.Songs) == 0 || single.Songs[0].Isrc != "" {
				enrichedAlbum = s.enrichAlbum(ctx, provider, single)
			} else {
				if tmp, err := s.GetAlbum(ctx, pluginInstances, single.Id); err != nil {
					return
				} else {
					enrichedAlbum = tmp
				}
			}
			singles = append(singles, enrichedAlbum)
		}(single)
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

func (s *PluginService) GetPlaylist(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string) (model.EnrichedPlaylist, error) {
	var playlist model.Playlist
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		playlist, err = plugin.Playlist(ctx, instances, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedPlaylist{}, fmt.Errorf("playlist %w: %v", model.ErrPluginDataNotFound, errMap)
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

func (s *PluginService) Search(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, q string) (model.EnrichedSearch, error) {
	var result model.Search
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		result, err = plugin.Search(ctx, instances, q, q, q, q)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedSearch{}, fmt.Errorf("search %w: %v", model.ErrPluginDataNotFound, errMap)
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
		enrichedArtist := s.enrichArtist(ctx, provider, pluginInstances, artist)

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

func (s *PluginService) Url(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, q string) (model.TypedItem, error) {
	var urlItem model.UrlItem
	errMap := map[string]string{}
	for plugin := range pluginInstances {
		var err error
		urlItem, err = plugin.Url(ctx, q)
		if err == nil {
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.TypedItem{}, fmt.Errorf("url parsed %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	var data any
	var err error
	switch urlItem.Type {
	case model.TypeSong:
		data, err = s.GetSong(ctx, pluginInstances, urlItem.Id)
	case model.TypeAlbum:
		data, err = s.GetAlbum(ctx, pluginInstances, urlItem.Id)
	case model.TypeArtist:
		data, err = s.GetArtist(ctx, pluginInstances, urlItem.Id)
	case model.TypePlaylist:
		data, err = s.GetPlaylist(ctx, pluginInstances, urlItem.Id)
	}
	if err != nil {
		return model.TypedItem{}, fmt.Errorf("from url %s: %w", q, err)
	}

	return model.TypedItem{
		Type: urlItem.Type,
		Data: data,
	}, nil
}

func (s *PluginService) GetSongByISRC(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, isrc string) (model.EnrichedSong, error) {
	var song model.Song
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		song, err = plugin.SongByISRC(ctx, instances, isrc)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedSong{}, fmt.Errorf("song by ISRC %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedSong := s.enrichSong(ctx, provider, song)

	return enrichedSong, nil
}
