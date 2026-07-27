package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type PluginService struct {
	store    *PluginStoreService
	instance model.InstanceRepository
	user     model.UserRepository
	follow   model.FollowRepository
}

func NewPluginService(store *PluginStoreService, instance model.InstanceRepository, user model.UserRepository, follow model.FollowRepository) PluginService {
	return PluginService{
		store:    store,
		instance: instance,
		user:     user,
		follow:   follow,
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

func (s *PluginService) GetStatus(ctx context.Context, pluginName string, url string) error {
	p, ok := s.store.GetPluginByName(pluginName)
	if !ok {
		return model.ErrPluginNotFound
	}

	if err := p.Status(ctx, url); err != nil {
		return err
	}

	return nil
}

// -- SongInfo -- //

func (s *PluginService) GetSongInfo(ctx context.Context, userID uuid.UUID, provider string, id string) (model.EnrichedSongInfo, error) {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
	if err != nil {
		return model.EnrichedSongInfo{}, fmt.Errorf("list instances of user %s for provider %s: %w", userID.String(), provider, err)
	}

	pluginInstances := s.InstancesToMapPluginInstances(instances)

	songInfo, err := s.GetSongInfoFromPluginInstances(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedSongInfo{}, fmt.Errorf("get song %s info from each plugin: %w", id, err)
	}

	return songInfo, nil
}

func (s *PluginService) GetSongInfoFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string) (model.EnrichedSongInfo, error) {
	var songInfo model.SongInfo
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		songInfo, err = plugin.SongInfo(ctx, instances, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedSongInfo{}, fmt.Errorf("song info %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedSongInfo := s.enrichSongInfo(ctx, provider, songInfo)

	return enrichedSongInfo, nil
}

func (s *PluginService) enrichSongInfo(ctx context.Context, provider string, songInfo model.SongInfo) model.EnrichedSongInfo {
	return model.EnrichedSongInfo{
		Provider: provider,
		SongInfo: songInfo,
	}
}

// -- AlbumInfo -- //

func (s *PluginService) GetAlbumInfo(ctx context.Context, userID uuid.UUID, provider string, id string) (model.EnrichedAlbumInfo, error) {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
	if err != nil {
		return model.EnrichedAlbumInfo{}, fmt.Errorf("list instances of user %s for provider %s: %w", userID.String(), provider, err)
	}

	pluginInstances := s.InstancesToMapPluginInstances(instances)

	albumInfo, err := s.GetAlbumInfoFromPluginInstances(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedAlbumInfo{}, fmt.Errorf("get album %s info from each plugin: %w", id, err)
	}

	return albumInfo, nil
}

func (s *PluginService) GetAlbumInfoFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string) (model.EnrichedAlbumInfo, error) {
	var albumInfo model.AlbumInfo
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		albumInfo, err = plugin.AlbumInfo(ctx, instances, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedAlbumInfo{}, fmt.Errorf("album info %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedAlbumInfo := s.enrichAlbumInfo(ctx, provider, albumInfo)

	return enrichedAlbumInfo, nil
}

func (s *PluginService) enrichAlbumInfo(ctx context.Context, provider string, albumInfo model.AlbumInfo) model.EnrichedAlbumInfo {
	return model.EnrichedAlbumInfo{
		Provider:  provider,
		AlbumInfo: albumInfo,
	}
}

// -- AlbumSongs -- //

func (s *PluginService) GetAlbumSongs(ctx context.Context, userID uuid.UUID, provider string, id string, limit int, offset int) (model.EnrichedAlbumSongs, error) {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
	if err != nil {
		return model.EnrichedAlbumSongs{}, fmt.Errorf("list instances of user %s for provider %s: %w", userID.String(), provider, err)
	}

	pluginInstances := s.InstancesToMapPluginInstances(instances)

	albumSongs, err := s.GetAlbumSongsFromPluginInstances(ctx, pluginInstances, id, limit, offset)
	if err != nil {
		return model.EnrichedAlbumSongs{}, fmt.Errorf("get album %s songs from each plugin: %w", id, err)
	}

	return albumSongs, nil
}

func (s *PluginService) GetAlbumSongsFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string, limit int, offset int) (model.EnrichedAlbumSongs, error) {
	var albumSongs model.AlbumSongs
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		albumSongs, err = plugin.AlbumSongs(ctx, instances, id, limit, offset)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedAlbumSongs{}, fmt.Errorf("album songs %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedAlbumSongs := s.enrichAlbumSongs(ctx, provider, albumSongs)

	return enrichedAlbumSongs, nil
}

func (s *PluginService) enrichAlbumSongs(ctx context.Context, provider string, albumSongs model.AlbumSongs) model.EnrichedAlbumSongs {
	songs := []model.EnrichedSongInfo{}
	for _, song := range albumSongs.Songs {
		songs = append(songs, s.enrichSongInfo(ctx, provider, song))
	}

	return model.EnrichedAlbumSongs{
		Provider:   provider,
		AlbumSongs: albumSongs,
		Songs:      songs,
	}
}

// -- ArtistInfo -- //

func (s *PluginService) GetArtistInfo(ctx context.Context, userID uuid.UUID, provider string, id string) (model.EnrichedArtistInfo, error) {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
	if err != nil {
		return model.EnrichedArtistInfo{}, fmt.Errorf("list instances of user %s for provider %s: %w", userID.String(), provider, err)
	}

	pluginInstances := s.InstancesToMapPluginInstances(instances)

	artistInfo, err := s.GetArtistInfoFromPluginInstances(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedArtistInfo{}, fmt.Errorf("get artist %s info from each plugin: %w", id, err)
	}

	return artistInfo, nil
}

func (s *PluginService) GetArtistInfoFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string) (model.EnrichedArtistInfo, error) {
	var artistInfo model.ArtistInfo
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		artistInfo, err = plugin.ArtistInfo(ctx, instances, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedArtistInfo{}, fmt.Errorf("artist info %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedArtistInfo := s.enrichArtistInfo(ctx, provider, artistInfo)

	return enrichedArtistInfo, nil
}

func (s *PluginService) enrichArtistInfo(ctx context.Context, provider string, artistInfo model.ArtistInfo) model.EnrichedArtistInfo {
	followed := uuid.UUID{}
	// if follow, err := s.follow.GetFollowByFilter(ctx, model.FollowFilter{Provider: &provider, ID: &artistInfo.ID}); err == nil {
	// 	followed = follow.ID
	// }

	return model.EnrichedArtistInfo{
		Provider:   provider,
		Followed:   followed,
		ArtistInfo: artistInfo,
	}

}

// -- ArtistAlbums -- //

func (s *PluginService) GetArtistAlbums(ctx context.Context, userID uuid.UUID, provider string, id string, limit int, offset int) (model.EnrichedArtistAlbums, error) {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
	if err != nil {
		return model.EnrichedArtistAlbums{}, fmt.Errorf("list instances of user %s for provider %s: %w", userID.String(), provider, err)
	}

	pluginInstances := s.InstancesToMapPluginInstances(instances)

	artistAlbums, err := s.GetArtistAlbumsFromPluginInstances(ctx, pluginInstances, id, limit, offset)
	if err != nil {
		return model.EnrichedArtistAlbums{}, fmt.Errorf("get artist %s albums from each plugin: %w", id, err)
	}

	return artistAlbums, nil
}

func (s *PluginService) GetArtistAlbumsFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string, limit int, offset int) (model.EnrichedArtistAlbums, error) {
	var artist model.ArtistAlbums
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		artist, err = plugin.ArtistAlbums(ctx, instances, id, limit, offset)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedArtistAlbums{}, fmt.Errorf("artist albums %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedArtistAlbums := s.enrichArtistAlbums(ctx, provider, artist)

	return enrichedArtistAlbums, nil
}

func (s *PluginService) enrichArtistAlbums(ctx context.Context, provider string, artistAlbums model.ArtistAlbums) model.EnrichedArtistAlbums {
	albums := []model.EnrichedAlbumInfo{}
	eps := []model.EnrichedAlbumInfo{}
	singles := []model.EnrichedAlbumInfo{}
	for _, album := range artistAlbums.Albums {
		albums = append(albums, s.enrichAlbumInfo(ctx, provider, album))
	}
	for _, ep := range artistAlbums.Ep {
		eps = append(eps, s.enrichAlbumInfo(ctx, provider, ep))
	}
	for _, single := range artistAlbums.Singles {
		singles = append(singles, s.enrichAlbumInfo(ctx, provider, single))
	}
	// var wg sync.WaitGroup
	// for _, album := range artistAlbums.Albums {
	// 	wg.Add(1)
	// 	go func(album model.Album) {
	// 		defer wg.Done()
	// 		enrichedAlbum := model.EnrichedAlbumInfo{}
	// 		// if len(album.Songs) == 0 || album.Songs[0].Isrc != "" {
	// 		// 	enrichedAlbum = s.enrichAlbum(ctx, provider, album)
	// 		// } else {
	// 		if tmp, err := s.GetAlbumFromPluginInstances(ctx, pluginInstances, album.ID); err != nil {
	// 			slog.Error("get album", slog.String("err", err.Error()))
	// 			return
	// 		} else {
	// 			enrichedAlbum = tmp
	// 		}
	// 		// }
	// 		albums = append(albums, enrichedAlbum)
	// 	}(album)
	// }
	// for _, ep := range artistAlbums.Ep {
	// 	wg.Add(1)
	// 	go func(ep model.Album) {
	// 		defer wg.Done()
	// 		enrichedAlbum := model.EnrichedAlbumInfo{}
	// 		// if len(ep.Songs) == 0 || ep.Songs[0].Isrc != "" {
	// 		// 	enrichedAlbum = s.enrichAlbum(ctx, provider, ep)
	// 		// } else {
	// 		if tmp, err := s.GetAlbumFromPluginInstances(ctx, pluginInstances, ep.ID); err != nil {
	// 			slog.Error("get ep", slog.String("err", err.Error()))
	// 			return
	// 		} else {
	// 			enrichedAlbum = tmp
	// 		}
	// 		// }
	// 		eps = append(eps, enrichedAlbum)
	// 	}(ep)
	// }
	// for _, single := range artistAlbums.Singles {
	// 	wg.Add(1)
	// 	go func(single model.Album) {
	// 		defer wg.Done()
	// 		enrichedAlbum := model.EnrichedAlbumInfo{}
	// 		// if len(single.Songs) == 0 || single.Songs[0].Isrc != "" {
	// 		// 	enrichedAlbum = s.enrichAlbum(ctx, provider, single)
	// 		// } else {
	// 		if tmp, err := s.GetAlbumFromPluginInstances(ctx, pluginInstances, single.ID); err != nil {
	// 			slog.Error("get single", slog.String("err", err.Error()))
	// 			return
	// 		} else {
	// 			enrichedAlbum = tmp
	// 		}
	// 		// }
	// 		singles = append(singles, enrichedAlbum)
	// 	}(single)
	// }
	// wg.Wait()

	return model.EnrichedArtistAlbums{
		Provider:     provider,
		ArtistAlbums: artistAlbums,
		Albums:       albums,
		Ep:           eps,
		Singles:      singles,
	}

}

// -- PlaylistInfo -- //

func (s *PluginService) GetPlaylistInfo(ctx context.Context, userID uuid.UUID, provider string, id string) (model.EnrichedPlaylistInfo, error) {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
	if err != nil {
		return model.EnrichedPlaylistInfo{}, fmt.Errorf("list instances of user %s for provider %s: %w", userID.String(), provider, err)
	}

	pluginInstances := s.InstancesToMapPluginInstances(instances)

	playlistInfo, err := s.GetPlaylistInfoFromPluginInstances(ctx, pluginInstances, id)
	if err != nil {
		return model.EnrichedPlaylistInfo{}, fmt.Errorf("get playlist %s info from each: %w", id, err)
	}

	return playlistInfo, nil
}

func (s *PluginService) GetPlaylistInfoFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string) (model.EnrichedPlaylistInfo, error) {
	var playlistInfo model.PlaylistInfo
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		playlistInfo, err = plugin.PlaylistInfo(ctx, instances, id)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedPlaylistInfo{}, fmt.Errorf("playlist info %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedPlaylistInfo := s.enrichPlaylistInfo(ctx, provider, playlistInfo)

	return enrichedPlaylistInfo, nil
}

func (s *PluginService) enrichPlaylistInfo(ctx context.Context, provider string, playlistInfo model.PlaylistInfo) model.EnrichedPlaylistInfo {
	return model.EnrichedPlaylistInfo{
		Provider:     provider,
		PlaylistInfo: playlistInfo,
	}
}

// -- PlaylistSongs -- //

func (s *PluginService) GetPlaylistSongs(ctx context.Context, userID uuid.UUID, provider string, id string, limit int, offset int) (model.EnrichedPlaylistSongs, error) {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
	if err != nil {
		return model.EnrichedPlaylistSongs{}, fmt.Errorf("list instances of user %s for provider %s: %w", userID.String(), provider, err)
	}

	pluginInstances := s.InstancesToMapPluginInstances(instances)

	playlistSongs, err := s.GetPlaylistSongsFromPluginInstances(ctx, pluginInstances, id, limit, offset)
	if err != nil {
		return model.EnrichedPlaylistSongs{}, fmt.Errorf("get playlist %s songs from each: %w", id, err)
	}

	return playlistSongs, nil
}

func (s *PluginService) GetPlaylistSongsFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string, limit int, offset int) (model.EnrichedPlaylistSongs, error) {
	var playlistSongs model.PlaylistSongs
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		playlistSongs, err = plugin.PlaylistSongs(ctx, instances, id, limit, offset)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedPlaylistSongs{}, fmt.Errorf("playlist songs %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedPlaylistSongs := s.enrichPlaylistSongs(ctx, provider, playlistSongs)

	return enrichedPlaylistSongs, nil
}

func (s *PluginService) enrichPlaylistSongs(ctx context.Context, provider string, playlistSongs model.PlaylistSongs) model.EnrichedPlaylistSongs {
	songs := []model.EnrichedSongInfo{}
	for _, song := range playlistSongs.Songs {
		songs = append(songs, s.enrichSongInfo(ctx, provider, song))
	}

	return model.EnrichedPlaylistSongs{
		Provider:      provider,
		PlaylistSongs: playlistSongs,
		Songs:         songs,
	}
}

// -- Search -- //

func (s *PluginService) Search(ctx context.Context, userID uuid.UUID, q string, limit int, offset int) (model.SearchResult, error) {
	if q == "" {
		return model.SearchResult{}, model.ErrPluginSearchEmptyQuery
	}

	providerResult := map[string]model.EnrichedSearch{}

	for provider := range s.store.ListPluginsByProvider() {
		instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
		if err != nil {
			continue
		}

		pluginInstances := s.InstancesToMapPluginInstances(instances)

		if item, err := s.UrlFromPluginInstances(ctx, pluginInstances, q); err == nil {
			return model.SearchResult{ItemFound: item}, nil
		}

		if s.IsISRC(q) {
			if data, err := s.GetSongByISRCFromPluginInstances(ctx, pluginInstances, q); err == nil {
				return model.SearchResult{ItemFound: model.TypedItem{
					Type: model.TypeSong,
					Data: data,
				}}, nil
			}
		}

		result, err := s.SearchFromPluginInstances(ctx, pluginInstances, q, limit, offset)
		if err != nil {
			continue
		}

		providerResult[provider] = result
	}

	return model.SearchResult{ProviderResult: providerResult}, nil
}

func (s *PluginService) SearchFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, q string, limit int, offset int) (model.EnrichedSearch, error) {
	var result model.Search
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		result, err = plugin.Search(ctx, instances, q, q, q, q, limit, offset)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedSearch{}, fmt.Errorf("search %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	songs := []model.EnrichedSongInfo{}
	for _, song := range result.Songs {
		songs = append(songs, s.enrichSongInfo(ctx, provider, song))
	}

	albums := []model.EnrichedAlbumInfo{}
	for _, album := range result.Albums {
		albums = append(albums, s.enrichAlbumInfo(ctx, provider, album))
	}

	artists := []model.EnrichedArtistInfo{}
	for _, artist := range result.Artists {
		artists = append(artists, s.enrichArtistInfo(ctx, provider, pluginInstances, artist))
	}

	playlists := []model.EnrichedPlaylistInfo{}
	for _, playlist := range result.Playlists {
		playlists = append(playlists, s.enrichPlaylistInfo(ctx, provider, playlist))
	}

	return model.EnrichedSearch{
		Songs:     songs,
		Albums:    albums,
		Artists:   artists,
		Playlists: playlists,
	}, nil
}

func (s *PluginService) UrlFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, q string) (model.TypedItem, error) {
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
		data, err = s.GetSongInfoFromPluginInstances(ctx, pluginInstances, urlItem.ID)
	case model.TypeAlbum:
		data, err = s.GetAlbumInfoFromPluginInstances(ctx, pluginInstances, urlItem.ID)
	case model.TypeArtist:
		data, err = s.GetArtistInfoFromPluginInstances(ctx, pluginInstances, urlItem.ID)
	case model.TypePlaylist:
		data, err = s.GetPlaylistInfoFromPluginInstances(ctx, pluginInstances, urlItem.ID)
	}
	if err != nil {
		return model.TypedItem{}, fmt.Errorf("from url %s: %w", q, err)
	}

	return model.TypedItem{
		Type: urlItem.Type,
		Data: data,
	}, nil
}

func (s *PluginService) GetSongByISRCFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, isrc string) (model.EnrichedSongInfo, error) {
	var song model.SongInfo
	var provider string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		song, err = plugin.SongInfoByISRC(ctx, instances, isrc)
		if err == nil {
			provider = plugin.Provider()
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return model.EnrichedSongInfo{}, fmt.Errorf("song by ISRC %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	enrichedSongInfo := s.enrichSongInfo(ctx, provider, song)

	return enrichedSongInfo, nil
}

func (s *PluginService) Download(ctx context.Context, userID uuid.UUID, provider string, id string, hiRes bool) (io.ReadCloser, string, error) {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID, Provider: &provider})
	if err != nil {
		return nil, "", fmt.Errorf("list instances of user %s for provider %s: %w", userID.String(), provider, err)
	}

	pluginInstances := s.InstancesToMapPluginInstances(instances)

	reader, extension, err := s.DownloadFromPluginInstances(ctx, pluginInstances, id, hiRes)
	if err != nil {
		return nil, "", fmt.Errorf("download %s from each: %w", id, err)
	}

	return reader, extension, nil
}

func (s *PluginService) DownloadFromPluginInstances(ctx context.Context, pluginInstances map[model.Plugin][]model.Instance, id string, hiRes bool) (io.ReadCloser, string, error) {
	var reader io.ReadCloser
	var extension string
	errMap := map[string]string{}
	for plugin, instances := range pluginInstances {
		var err error
		reader, extension, err = plugin.Download(ctx, instances, id, hiRes)
		if err == nil {
			break
		}
		errMap[plugin.Name()] = err.Error()
	}
	if len(errMap) == len(pluginInstances) {
		return nil, "", fmt.Errorf("download %w: %v", model.ErrPluginDataNotFound, errMap)
	}

	return reader, extension, nil
}
