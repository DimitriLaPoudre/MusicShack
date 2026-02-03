package handlers

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
)

func AddFollow(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.RequestFollow
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plugins, ok := plugins.GetProvider(req.Provider)
	if !ok {
		log.Println("invalid provider name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider name"})
		return
	}

	var artist models.ArtistData
	for _, plugin := range plugins {
		artist, err = plugin.Artist(c.Request.Context(), userId, req.Id)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	follow, err := repository.AddFollow(userId, req.Provider, req.Id, artist.Name, artist.PictureUrl)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := models.FollowItem{
		Id:               follow.ID,
		Provider:         follow.Provider,
		ArtistId:         follow.ArtistId,
		ArtistName:       follow.ArtistName,
		ArtistPictureUrl: follow.ArtistPictureUrl,
	}

	c.JSON(http.StatusOK, data)
}

func ListFollows(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followsRaw, err := repository.ListFollowsByUserID(userId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	follows := make([]models.FollowItem, len(followsRaw))
	var wg sync.WaitGroup

	for index, value := range followsRaw {
		wg.Add(1)
		go func(i int, follow models.Follow) {
			defer wg.Done()

			follows[i] = models.FollowItem{Id: follow.ID, Provider: follow.Provider, ArtistId: follow.ArtistId, ArtistName: follow.ArtistName, ArtistPictureUrl: follow.ArtistPictureUrl}
		}(index, value)
	}

	wg.Wait()

	c.JSON(http.StatusOK, follows)
}

func DeleteFollow(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(idUint64)

	if err := repository.DeleteFollowByUserID(userId, id); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
