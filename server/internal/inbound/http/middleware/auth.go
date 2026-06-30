package middleware

import (
	"net/http"

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
