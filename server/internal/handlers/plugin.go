package handlers

import (
	"fmt"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	api := c.Param("api")
	id := c.Param("id")

	p, ok := plugins.Get(api)
	if !ok {
		fmt.Println("invalid api name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	data, err := p.Song(c.Request.Context(), userId, id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetAlbum(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	api := c.Param("api")
	id := c.Param("id")

	p, ok := plugins.Get(api)
	if !ok {
		fmt.Println("invalid api name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	data, err := p.Album(c.Request.Context(), userId, id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetArtist(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	api := c.Param("api")
	id := c.Param("id")

	p, ok := plugins.Get(api)
	if !ok {
		fmt.Println("invalid api name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	data, err := p.Artist(c.Request.Context(), userId, id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func Search(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	search := c.Query("q")
	finding := make(map[string]models.SearchData)

	for key, value := range plugins.GetRegistry() {
		tmp, err := value.Search(c.Request.Context(), userId, search, search, search)
		if err == nil &&
			(len(tmp.Songs) != 0 || len(tmp.Albums) != 0 || len(tmp.Artists) != 0) {
			finding[key] = tmp
		}
	}
	c.JSON(http.StatusOK, finding)
}
