package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
)

func ListSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var limit uint
	if result := c.Query("limit"); result == "" {
		limit = 20
	} else {
		if result, err := strconv.ParseUint(result, 10, 0); err != nil {
			limit = 20
		} else {
			limit = uint(result)
		}
	}

	var offset uint
	if result := c.Query("offset"); result == "" {
		offset = 0
	} else {
		if result, err := strconv.ParseUint(result, 10, 0); err != nil {
			offset = 0
		} else {
			offset = uint(result)
		}
	}

	songs, err := repository.ListSongByUserID(userId, limit, offset)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(songs)

	list := make([]models.ResponseSong, 0)
	for _, song := range songs {
		if song, err := services.GetLibrarySong(song); err != nil {
			log.Println("ListSong:", err)
			continue
		} else {
			list = append(list, song)
		}
	}

	c.JSON(http.StatusOK, list)
}

func DeleteSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id uint
	if result, err := strconv.ParseUint(c.Param("id"), 10, 0); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		id = uint(result)
	}

	if err := repository.DeleteSongByUserID(userId, id); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func SyncLibrary(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.SyncUserLibrary(userId); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
