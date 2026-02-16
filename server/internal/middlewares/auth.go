package middlewares

import (
	"errors"
	"log"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/gin-gonic/gin"
)

func Logged() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("user_session")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		session, err := repository.GetUserSessionByToken(token)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		repository.DeleteExpiredUserSession()

		c.Set("userId", session.UserId)
		c.Next()
	}
}

func LoggedOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("user_session")
		if err != nil {
			c.Next()
			return
		}

		_, err = repository.GetUserSessionByToken(token)
		if err != nil {
			c.Next()
			return
		}
		repository.DeleteExpiredUserSession()

		err = errors.New("user already logged in")
		log.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		c.Abort()
	}
}
