package middleware

import (
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/utils"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func UserMiddleware(l *zerolog.Logger, auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		me, err := utils.GetFromContext[*model.User](c, "me")
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
