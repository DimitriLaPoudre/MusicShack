package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"reflect"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/sync"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
)

type PluginCacheService struct {
	cfgPlugin config.PluginConfig
	cache     sync.Cache[CacheKey]
}

func NewPluginCacheService(cfgPlugin config.PluginConfig) PluginCacheService {
	return PluginCacheService{
		cfgPlugin: cfgPlugin,
		cache:     sync.NewCache[CacheKey](cfgPlugin.Cache.Expiration),
	}
}

type CacheKey [32]byte

func (s *PluginCacheService) HashKey(fn any, args ...any) CacheKey {
	v := reflect.ValueOf(fn)

	if v.Kind() != reflect.Func {
		panic("not a function")
	}

	ptr := v.Pointer()
	b, _ := json.Marshal(map[string]any{"ptr": ptr, "args": args})
	return sha256.Sum256(b)
}

func (s *PluginCacheService) SongInfo(plugin model.Plugin, ctx context.Context, instances []model.Instance, id string) (model.SongInfo, error) {
	key := s.HashKey(plugin.SongInfo, instances, id)
	if cached, ok := s.cache.Load(key); ok {
		if songInfo := cached.(model.SongInfo); ok {
			return songInfo, nil
		}
	}

	songInfo, err := plugin.SongInfo(ctx, instances, id)
	if err == nil {
		s.cache.Store(key, songInfo)
	}

	return songInfo, nil
}

func (s *PluginCacheService) SongInfoByISRC(plugin model.Plugin, ctx context.Context, instances []model.Instance, isrc string) (model.SongInfo, error) {
	key := s.HashKey(plugin.SongInfoByISRC, instances, isrc)
	if cached, ok := s.cache.Load(key); ok {
		if songInfo := cached.(model.SongInfo); ok {
			return songInfo, nil
		}
	}

	songInfo, err := plugin.SongInfoByISRC(ctx, instances, isrc)
	if err == nil {
		s.cache.Store(key, songInfo)
	}

	return songInfo, nil
}

func (s *PluginCacheService) AlbumInfo(plugin model.Plugin, ctx context.Context, instances []model.Instance, id string) (model.AlbumInfo, error) {
	key := s.HashKey(plugin.AlbumInfo, instances, id)
	if cached, ok := s.cache.Load(key); ok {
		if albumInfo := cached.(model.AlbumInfo); ok {
			return albumInfo, nil
		}
	}

	albumInfo, err := plugin.AlbumInfo(ctx, instances, id)
	if err == nil {
		s.cache.Store(key, albumInfo)
	}

	return albumInfo, nil
}

func (s *PluginCacheService) AlbumSongs(plugin model.Plugin, ctx context.Context, instances []model.Instance, id string, limit int, offset int) (model.PaginatedSongs, error) {
	key := s.HashKey(plugin.AlbumSongs, instances, id, limit, offset)
	if cached, ok := s.cache.Load(key); ok {
		if albumSongs := cached.(model.PaginatedSongs); ok {
			return albumSongs, nil
		}
	}

	albumSongs, err := plugin.AlbumSongs(ctx, instances, id, limit, offset)
	if err == nil {
		s.cache.Store(key, albumSongs)
	}

	return albumSongs, nil
}

func (s *PluginCacheService) ArtistInfo(plugin model.Plugin, ctx context.Context, instances []model.Instance, id string) (model.ArtistInfo, error) {
	key := s.HashKey(plugin.ArtistInfo, instances, id)
	if cached, ok := s.cache.Load(key); ok {
		if artistInfo := cached.(model.ArtistInfo); ok {
			return artistInfo, nil
		}
	}

	artistInfo, err := plugin.ArtistInfo(ctx, instances, id)
	if err == nil {
		s.cache.Store(key, artistInfo)
	}

	return artistInfo, nil
}

func (s *PluginCacheService) ArtistAlbums(plugin model.Plugin, ctx context.Context, instances []model.Instance, id string, limit int, offset int) (model.ArtistPaginatedAlbums, error) {
	key := s.HashKey(plugin.ArtistAlbums, instances, id, limit, offset)
	if cached, ok := s.cache.Load(key); ok {
		if artistAlbums := cached.(model.ArtistPaginatedAlbums); ok {
			return artistAlbums, nil
		}
	}

	artistAlbums, err := plugin.ArtistAlbums(ctx, instances, id, limit, offset)
	if err == nil {
		s.cache.Store(key, artistAlbums)
	}

	return artistAlbums, nil
}

func (s *PluginCacheService) PlaylistInfo(plugin model.Plugin, ctx context.Context, instances []model.Instance, id string) (model.PlaylistInfo, error) {
	key := s.HashKey(plugin.PlaylistInfo, instances, id)
	if cached, ok := s.cache.Load(key); ok {
		if playlistInfo := cached.(model.PlaylistInfo); ok {
			return playlistInfo, nil
		}
	}

	playlistInfo, err := plugin.PlaylistInfo(ctx, instances, id)
	if err == nil {
		s.cache.Store(key, playlistInfo)
	}

	return playlistInfo, nil
}

func (s *PluginCacheService) PlaylistSongs(plugin model.Plugin, ctx context.Context, instances []model.Instance, id string, limit int, offset int) (model.PaginatedSongs, error) {
	key := s.HashKey(plugin.PlaylistSongs, instances, id, limit, offset)
	if cached, ok := s.cache.Load(key); ok {
		if playlistSongs := cached.(model.PaginatedSongs); ok {
			return playlistSongs, nil
		}
	}

	playlistSongs, err := plugin.PlaylistSongs(ctx, instances, id, limit, offset)
	if err == nil {
		s.cache.Store(key, playlistSongs)
	}

	return playlistSongs, nil
}

func (s *PluginCacheService) Search(plugin model.Plugin, ctx context.Context, instances []model.Instance, song string, album string, artist string, playlist string, limit int, offset int) (model.Search, error) {
	key := s.HashKey(plugin.Search, instances, song, album, artist, playlist, limit, offset)
	if cached, ok := s.cache.Load(key); ok {
		if search := cached.(model.Search); ok {
			return search, nil
		}
	}

	search, err := plugin.Search(ctx, instances, song, album, artist, playlist, limit, offset)
	if err == nil {
		s.cache.Store(key, search)
	}

	return search, nil
}

func (s *PluginCacheService) URL(plugin model.Plugin, ctx context.Context, url string) (model.URLItem, error) {
	key := s.HashKey(plugin.URL, url)
	if cached, ok := s.cache.Load(key); ok {
		if urlItem := cached.(model.URLItem); ok {
			return urlItem, nil
		}
	}

	urlItem, err := plugin.URL(ctx, url)
	if err == nil {
		s.cache.Store(key, urlItem)
	}

	return urlItem, nil
}
