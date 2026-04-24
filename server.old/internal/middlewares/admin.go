package middlewares

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/gin-gonic/gin"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("admin_session")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		admin, err := repository.GetAdmin()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if admin.Token != token || admin.ExpiresAt.Before(time.Now()) {
			err := errors.New("invalid token")
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func Admout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("admin_session")
		if err != nil {
			c.Next()
			return
		}

		admin, err := repository.GetAdmin()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if admin.Token != token || admin.ExpiresAt.Before(time.Now()) {
			c.Next()
			return
		}
		err = errors.New("admin can't access this ressources")
		log.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		c.Abort()
	}
}
