package handlers

import (
	"fmt"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetSong(c *gin.Context) {
	api := c.Param("api")
	id := c.Param("id")

	p, ok := plugins.Get(api)
	if !ok {
		fmt.Println("invalid api name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	data, err := p.Song(c.Request.Context(), id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetAlbum(c *gin.Context) {
	api := c.Param("api")
	id := c.Param("id")

	p, ok := plugins.Get(api)
	if !ok {
		fmt.Println("invalid api name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	data, err := p.Album(c.Request.Context(), id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetArtist(c *gin.Context) {
	api := c.Param("api")
	id := c.Param("id")

	p, ok := plugins.Get(api)
	if !ok {
		fmt.Println("invalid api name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	data, err := p.Artist(c.Request.Context(), id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func Search(c *gin.Context) {
	search := c.Query("q")
	finding := make(map[string]models.SearchData)

	for key, value := range plugins.GetRegistry() {
		tmp, err := value.Search(c.Request.Context(), search, search, search)
		if err == nil {
			finding[key] = tmp
		}
	}
	c.JSON(http.StatusOK, finding)
}

func AddDownloadSong(c *gin.Context) {
	api := c.Param("api")
	id := c.Param("id")
	qualityAudio := c.Query("qualityAudio")

	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, ok := plugins.Get(api)
	if !ok {
		fmt.Println("invalid api name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	services.DownloadManager.Add(userId, p, id, qualityAudio)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
