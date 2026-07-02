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

func (u *PluginUseCase) GetSong(c context.Context, user *model.User, provider string, id string) (model.Song, error) {
	instances, err := u.instance.ListInstancesByFilter(c, &model.InstanceFilter{UserID: &user.ID, Provider: &provider})
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
			if song.Id != "" {
				break
			}
		}
	}
	if song.Id == "" {
		return model.Song{}, err
	}

	// _, err = repository.GetSongByUserIDByISRC(userId, data.Isrc)
	// if err == nil {
	// 	data.Downloaded = true
	// }

	return song, nil
}

// func (u *PluginUseCase) GetPlaylist(c *gin.Context) {
// 	userId, err := utils.GetFromContext[uint](c, "userId")
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	provider := c.Param("provider")
// 	id := c.Param("id")
//
// 	data, err := plugins.GetPlaylist(c.Request.Context(), userId, provider, id)
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	data.Downloaded = true
// 	for i, song := range data.Songs {
// 		_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
// 		if err == nil {
// 			data.Songs[i].Downloaded = true
// 		} else {
// 			data.Downloaded = false
// 		}
// 	}
//
// 	c.JSON(http.StatusOK, data)
// }
//
// func (u *PluginUseCase) GetAlbum(c *gin.Context) {
// 	userId, err := utils.GetFromContext[uint](c, "userId")
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	provider := c.Param("provider")
// 	id := c.Param("id")
//
// 	data, err := plugins.GetAlbum(c.Request.Context(), userId, provider, id)
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	data.Downloaded = true
// 	for i, song := range data.Songs {
// 		_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
// 		if err == nil {
// 			data.Songs[i].Downloaded = true
// 		} else {
// 			data.Downloaded = false
// 		}
// 	}
//
// 	c.JSON(http.StatusOK, data)
// }
//
// func (u *PluginUseCase) GetArtist(c *gin.Context) {
// 	userId, err := utils.GetFromContext[uint](c, "userId")
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	provider := c.Param("provider")
// 	id := c.Param("id")
//
// 	data, err := plugins.GetArtist(c.Request.Context(), userId, provider, id)
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	follow, err := repository.GetFollowByProviderByArtistID(data.Provider, data.Id)
// 	if err == nil {
// 		data.Followed = follow.ID
// 	}
//
// 	var wg sync.WaitGroup
// 	for i, album := range data.Albums {
// 		wg.Add(1)
// 		go func(i int) {
// 			defer wg.Done()
// 			albumData, err := plugins.GetAlbum(c.Request.Context(), userId, provider, album.Id)
// 			if err != nil {
// 				return
// 			}
// 			downloaded := true
// 			for _, song := range albumData.Songs {
// 				_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
// 				if err != nil {
// 					downloaded = false
// 					break
// 				}
// 			}
// 			data.Albums[i].Downloaded = downloaded
// 		}(i)
// 	}
// 	wg.Wait()
//
// 	c.JSON(http.StatusOK, data)
// }
//
// func (u *PluginUseCase) Search(c *gin.Context) {
// 	userId, err := utils.GetFromContext[uint](c, "userId")
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	search := c.Query("q")
// 	finding := make(map[string]models.SearchData)
//
// 	for provider, pluginsList := range plugins.GetAllPluginsByProvider() {
// 		var tmp models.SearchData
// 		var err error
// 		for _, plugin := range pluginsList {
// 			if urlItem, err := plugin.Url(c.Request.Context(), userId, search); err == nil {
// 				c.JSON(http.StatusOK, gin.H{"url": urlItem})
// 				return
// 			}
// 			tmp, err = plugin.Search(c.Request.Context(), userId, search, search, search)
// 			if err == nil {
// 				break
// 			}
// 		}
// 		if err == nil {
// 			var wg sync.WaitGroup
// 			for i, song := range tmp.Songs {
// 				_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
// 				if err == nil {
// 					tmp.Songs[i].Downloaded = true
// 				}
// 			}
// 			for i, album := range tmp.Albums {
// 				wg.Add(1)
// 				go func(i int) {
// 					defer wg.Done()
// 					albumData, err := plugins.GetAlbum(c.Request.Context(), userId, provider, album.Id)
// 					if err != nil {
// 						return
// 					}
// 					downloaded := true
// 					for _, song := range albumData.Songs {
// 						_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
// 						if err != nil {
// 							downloaded = false
// 							break
// 						}
// 					}
// 					tmp.Albums[i].Downloaded = downloaded
// 				}(i)
// 			}
// 			for i, artist := range tmp.Artists {
// 				follow, err := repository.GetFollowByProviderByArtistID(provider, artist.Id)
// 				if err == nil {
// 					tmp.Artists[i].Followed = follow.ID
// 				}
// 			}
// 			for i, playlist := range tmp.Playlists {
// 				wg.Add(1)
// 				go func(i int) {
// 					defer wg.Done()
// 					playlist, err := plugins.GetPlaylist(c.Request.Context(), userId, provider, playlist.ID)
// 					if err != nil {
// 						return
// 					}
// 					downloaded := true
// 					for _, song := range playlist.Songs {
// 						_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
// 						if err != nil {
// 							downloaded = false
// 							break
// 						}
// 					}
// 					tmp.Playlists[i].Downloaded = downloaded
// 				}(i)
// 			}
// 			wg.Wait()
// 			finding[provider] = tmp
// 		}
// 	}
// 	c.JSON(http.StatusOK, gin.H{"result": finding})
// }
