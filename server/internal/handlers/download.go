package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
)

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
	services.DownloadManager.AddSong(userId, p, id, qualityAudio)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AddDownloadAlbum(c *gin.Context) {
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
	services.DownloadManager.AddAlbum(userId, p, id, qualityAudio)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AddDownloadArtist(c *gin.Context) {
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
	services.DownloadManager.AddArtist(userId, p, id, qualityAudio)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ListDownload(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks := services.DownloadManager.List(userId)
	c.JSON(http.StatusOK, tasks)
}

func DeleteDownload(c *gin.Context) {
	taskIdRaw := c.Param("id")
	taskIdBadType, err := strconv.ParseUint(taskIdRaw, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId := uint(taskIdBadType)
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = services.DownloadManager.Remove(userId, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func RetryDownload(c *gin.Context) {
	taskIdRaw := c.Param("id")
	taskIdBadType, err := strconv.ParseUint(taskIdRaw, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId := uint(taskIdBadType)
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = services.DownloadManager.Retry(userId, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

func CancelDownload(c *gin.Context) {
	taskIdRaw := c.Param("id")
	taskIdBadType, err := strconv.ParseUint(taskIdRaw, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId := uint(taskIdBadType)
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = services.DownloadManager.Cancel(userId, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}
