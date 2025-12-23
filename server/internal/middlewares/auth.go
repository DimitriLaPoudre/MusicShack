package middlewares

import (
	"fmt"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/gin-gonic/gin"
)

func Logged() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("user_session")
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		session, err := repository.FindUserSessionByToken(token)
		if err != nil {
			fmt.Println(err)
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

		_, err = repository.FindUserSessionByToken(token)
		repository.DeleteExpiredUserSession()
		if err != nil {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "user already logged in"})
		c.Abort()
	}
}
