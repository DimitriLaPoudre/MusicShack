package middlewares

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/gin-gonic/gin"
)

func Logged() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		tokenUnsigned, err := services.GetTokenContent(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("userId", tokenUnsigned.Id)
		c.Next()
	}
}

func LoggedOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.Next()
			return
		}

		_, err = services.GetTokenContent(token)
		if err != nil {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "user already logged in"})
		c.Abort()
	}
}
