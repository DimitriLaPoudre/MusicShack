package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddFollow(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.RequestFollow
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := req.Id
	api, ok := plugins.Get(req.Api)
	if !ok {
		fmt.Println("invalid api name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api name"})
		return
	}

	if _, err := api.Artist(c.Request.Context(), id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.AddFollow(userId, api.Name(), id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ListFollows(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followsRaw, err := repository.ListFollowsByUserID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	follows := make([]models.FollowItem, len(followsRaw))
	var wg sync.WaitGroup

	for index, value := range followsRaw {
		wg.Add(1)
		go func(i int, follow models.Follow) {
			defer wg.Done()

			p, ok := plugins.Get(follow.Api)
			if !ok {
				return
			}

			artist, err := p.Artist(c.Request.Context(), follow.ArtistId)
			if err != nil {
				return
			}

			follows[i] = models.FollowItem{Id: follow.ID, Artist: artist}
		}(index, value)
	}

	wg.Wait()

	c.JSON(http.StatusOK, follows)
}

func DeleteFollow(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, strconv.IntSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(idUint64)

	if err := repository.DeleteFollowByUserID(userId, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "follow not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
