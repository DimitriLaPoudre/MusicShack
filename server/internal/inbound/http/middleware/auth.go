package middleware

import (
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/utils"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/service"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func AuthMiddleware(l *zerolog.Logger, cfg config.SessionConfig, auth *service.AuthService) gin.HandlerFunc {
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

func AdminMiddleware(l *zerolog.Logger, auth *service.AuthService) gin.HandlerFunc {
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

func UserMiddleware(l *zerolog.Logger, auth *service.AuthService) gin.HandlerFunc {
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

func GuestMiddleware(l *zerolog.Logger, auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := utils.GetFromContext[model.User](c, "me"); err != nil {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}
