package usecase

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/service"
	"github.com/rs/zerolog"
)

type PluginUseCase struct {
	l        *zerolog.Logger
	store    *service.PluginStoreService
	instance model.InstanceRepository
}

func NewPluginUseCase(l *zerolog.Logger, store *service.PluginStoreService, instance model.InstanceRepository) PluginUseCase {
	return PluginUseCase{
		l:        l,
		store:    store,
		instance: instance,
	}
}

func (u *PluginUseCase) GetSong(c context.Context, user model.User, provider string, id string) (model.Song, error) {
	instances, err := u.instance.ListInstancesByFilter(c, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.Song{}, err
	}

	var song model.Song
	for _, instance := range instances {
		if plugin, ok := u.store.GetPluginByName(instance.Plugin); !ok {
			err = model.ErrPluginNotFound
			continue
		} else {
			song, err = plugin.Song(c, instance.Url, id)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return model.Song{}, err
	}

	// _, err = repository.GetSongByUserIDByISRC(userId, data.Isrc)
	// if err == nil {
	// 	data.Downloaded = true
	// }

	return song, nil
}

func (u *PluginUseCase) GetAlbum(c context.Context, user model.User, provider string, id string) (model.Album, error) {
	instances, err := u.instance.ListInstancesByFilter(c, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.Album{}, err
	}

	var album model.Album
	for _, instance := range instances {
		if plugin, ok := u.store.GetPluginByName(instance.Plugin); !ok {
			err = model.ErrPluginNotFound
			continue
		} else {
			album, err = plugin.Album(c, instance.Url, id)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return model.Album{}, err
	}

	// album.Downloaded = true
	// for i, song := range album.Songs {
	// 	_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
	// 	if err == nil {
	// 		data.Songs[i].Downloaded = true
	// 	} else {
	// 		data.Downloaded = false
	// 	}
	// }

	return album, nil
}

func (u *PluginUseCase) GetArtist(c context.Context, user model.User, provider string, id string) (model.Artist, error) {
	instances, err := u.instance.ListInstancesByFilter(c, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.Artist{}, err
	}

	var artist model.Artist
	for _, instance := range instances {
		if plugin, ok := u.store.GetPluginByName(instance.Plugin); !ok {
			err = model.ErrPluginNotFound
			continue
		} else {
			artist, err = plugin.Artist(c, instance.Url, id)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return model.Artist{}, err
	}

	// follow, err := repository.GetFollowByProviderByArtistID(data.Provider, data.Id)
	// if err == nil {
	// 	data.Followed = follow.ID
	// }

	// var wg sync.WaitGroup
	// for i, album := range data.Albums {
	// 	wg.Add(1)
	// 	go func(i int) {
	// 		defer wg.Done()
	// 		albumData, err := plugins.GetAlbum(c.Request.Context(), userId, provider, album.Id)
	// 		if err != nil {
	// 			return
	// 		}
	// 		downloaded := true
	// 		for _, song := range albumData.Songs {
	// 			_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
	// 			if err != nil {
	// 				downloaded = false
	// 				break
	// 			}
	// 		}
	// 		data.Albums[i].Downloaded = downloaded
	// 	}(i)
	// }
	// wg.Wait()

	return artist, nil
}

func (u *PluginUseCase) GetPlaylist(c context.Context, user model.User, provider string, id string) (model.Playlist, error) {
	instances, err := u.instance.ListInstancesByFilter(c, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return model.Playlist{}, err
	}

	var playlist model.Playlist
	for _, instance := range instances {
		if plugin, ok := u.store.GetPluginByName(instance.Plugin); !ok {
			err = model.ErrPluginNotFound
			continue
		} else {
			playlist, err = plugin.Playlist(c, instance.Url, id)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return model.Playlist{}, err
	}

	// data.Downloaded = true
	// for i, song := range data.Songs {
	// 	_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
	// 	if err == nil {
	// 		data.Songs[i].Downloaded = true
	// 	} else {
	// 		data.Downloaded = false
	// 	}
	// }

	return playlist, nil
}

func (u *PluginUseCase) GetSearch(c context.Context, user model.User, q string) (map[string]model.Search, error) {
	providerResult := make(map[string]model.Search)

	for provider, _ := range u.store.ListPluginsByProvider() {
		instances, err := u.instance.ListInstancesByFilter(c, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
		if err != nil {
			continue
		}

		var search model.Search
		for _, instance := range instances {
			if plugin, ok := u.store.GetPluginByName(instance.Plugin); !ok {
				err = model.ErrPluginNotFound
				continue
			} else {
				search, err = plugin.Search(c, instance.Url, q, q, q)
				if err == nil {
					break
				}
			}
		}
		if err == nil {
			// var wg sync.WaitGroup
			// for i, song := range tmp.Songs {
			// 	_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
			// 	if err == nil {
			// 		tmp.Songs[i].Downloaded = true
			// 	}
			// }
			// for i, album := range tmp.Albums {
			// 	wg.Add(1)
			// 	go func(i int) {
			// 		defer wg.Done()
			// 		albumData, err := plugins.GetAlbum(c.Request.Context(), userId, provider, album.Id)
			// 		if err != nil {
			// 			return
			// 		}
			// 		downloaded := true
			// 		for _, song := range albumData.Songs {
			// 			_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
			// 			if err != nil {
			// 				downloaded = false
			// 				break
			// 			}
			// 		}
			// 		tmp.Albums[i].Downloaded = downloaded
			// 	}(i)
			// }
			// for i, artist := range tmp.Artists {
			// 	follow, err := repository.GetFollowByProviderByArtistID(provider, artist.Id)
			// 	if err == nil {
			// 		tmp.Artists[i].Followed = follow.ID
			// 	}
			// }
			// for i, playlist := range tmp.Playlists {
			// 	wg.Add(1)
			// 	go func(i int) {
			// 		defer wg.Done()
			// 		playlist, err := plugins.GetPlaylist(c.Request.Context(), userId, provider, playlist.ID)
			// 		if err != nil {
			// 			return
			// 		}
			// 		downloaded := true
			// 		for _, song := range playlist.Songs {
			// 			_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
			// 			if err != nil {
			// 				downloaded = false
			// 				break
			// 			}
			// 		}
			// 		tmp.Playlists[i].Downloaded = downloaded
			// 	}(i)
			// }
			// wg.Wait()
			providerResult[provider] = search
		}
	}

	return providerResult, nil
}
