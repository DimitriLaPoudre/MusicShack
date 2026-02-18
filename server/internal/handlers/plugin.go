package handlers

import (
	"log"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	data, err := plugins.GetSong(c.Request.Context(), userId, provider, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = repository.GetSongByUserIDByISRC(userId, data.Isrc)
	if err == nil {
		data.Downloaded = true
	}

	c.JSON(http.StatusOK, data)
}

func GetAlbum(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	data, err := plugins.GetAlbum(c.Request.Context(), userId, provider, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data.Downloaded = true
	for i, song := range data.Songs {
		_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
		if err == nil {
			data.Songs[i].Downloaded = true
		} else {
			data.Downloaded = false
		}
	}

	c.JSON(http.StatusOK, data)
}

func GetArtist(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	data, err := plugins.GetArtist(c.Request.Context(), userId, provider, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	follow, err := repository.GetFollowByProviderByArtistID(data.Provider, data.Id)
	if err == nil {
		data.Followed = follow.ID
	}

	c.JSON(http.StatusOK, data)
}

func Search(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := c.Query("q")
	finding := make(map[string]models.SearchData)

	for provider, plugins := range plugins.GetAllPluginsByProvider() {
		var tmp models.SearchData
		var err error
		for _, plugin := range plugins {
			if urlItem, err := plugin.Url(c.Request.Context(), userId, search); err == nil {
				c.JSON(http.StatusOK, gin.H{"url": urlItem})
				return
			}
			tmp, err = plugin.Search(c.Request.Context(), userId, search, search, search)
			if err == nil {
				break
			}
		}
		if err == nil {
			for i, song := range tmp.Songs {
				_, err = repository.GetSongByUserIDByISRC(userId, song.Isrc)
				if err == nil {
					tmp.Songs[i].Downloaded = true
				}
			}
			for i, artist := range tmp.Artists {
				follow, err := repository.GetFollowByProviderByArtistID(provider, artist.Id)
				if err == nil {
					tmp.Artists[i].Followed = follow.ID
				}
			}
			finding[provider] = tmp
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": finding})
}
