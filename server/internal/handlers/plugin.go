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

	for provider, plugins := range plugins.GetAllProvider() {
		var tmp models.SearchData
		var err error
		for _, plugin := range plugins {
			tmp, err = plugin.Search(c.Request.Context(), userId, search, search, search)
			if err == nil {
				break
			}
		}
		if err == nil {
			for i, artist := range tmp.Artists {
				follow, err := repository.GetFollowByProviderByArtistID(provider, artist.Id)
				if err == nil {
					tmp.Artists[i].Followed = follow.ID
				}
			}
			finding[provider] = tmp
		}
	}
	c.JSON(http.StatusOK, finding)
}
