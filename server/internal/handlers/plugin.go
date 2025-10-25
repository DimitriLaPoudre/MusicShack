package handlers

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"

	"github.com/gin-gonic/gin"
)

func GetSong(c *gin.Context) {
	api := c.Param("api")
	id := c.Param("id")

	p, ok := plugins.Get(api)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	data, err := p.Song(id)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}
	data, err := p.Album(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
