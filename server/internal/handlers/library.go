package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListSong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func DeleteSong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func SyncLibrary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
