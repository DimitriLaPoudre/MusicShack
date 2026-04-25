package middleware

import (
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func AdminMiddleware(l *zerolog.Logger, admin *usecase.AdminUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn, err := c.Cookie("admin_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err := admin.Authenticate(c.Request.Context(), tkn); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
