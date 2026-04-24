package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GinPrettyError(c *gin.Context, code int, err error) {
	log.Println(err)
	c.JSON(code, gin.H{"error": err.Error()})
}
