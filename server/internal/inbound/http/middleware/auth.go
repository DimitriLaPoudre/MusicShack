package middleware

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg config.SessionConfig, auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn, err := c.Cookie(cfg.CookieName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		me, err := auth.Authenticate(c.Request.Context(), tkn)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("me", me)

		c.Next()
	}
}

func AdminMiddleware(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		me, err := utils.GetFromContext[model.User](c, "me")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !auth.IsAdmin(c, me) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func UserMiddleware(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		me, err := utils.GetFromContext[model.User](c, "me")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !auth.IsUser(c, me) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func GuestMiddleware(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := utils.GetFromContext[model.User](c, "me"); err != nil {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}
