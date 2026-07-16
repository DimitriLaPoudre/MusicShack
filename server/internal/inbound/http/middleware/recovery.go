package middleware

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/gin-gonic/gin"
)

func isBrokenPipe(err any) bool {
	if ne, ok := err.(*net.OpError); ok {
		var se *os.SyscallError
		if errors.As(ne, &se) {
			errStr := strings.ToLower(se.Error())

			if strings.Contains(errStr, "broken pipe") {
				return true
			}

			if strings.Contains(errStr, "connection reset by peer") {
				return true
			}

			if strings.Contains(errStr, "connection aborted") {
				return true
			}
		}
	}
	return false
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, err := utils.GetFromContext[string](c, "request_id")
				if err != nil {
					requestID = "unknown"
				}
				attrs := []slog.Attr{
					slog.String("request_id", requestID),
					slog.String("ip", c.ClientIP()),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.Any("error", err),
				}

				if isBrokenPipe(err) {
					slog.LogAttrs(
						c.Request.Context(),
						slog.LevelWarn,
						"client connection broken",
						attrs...,
					)

					c.Abort()
					return
				}

				attrs = append(attrs,
					slog.String("stack", string(debug.Stack())),
				)

				slog.LogAttrs(
					c.Request.Context(),
					slog.LevelError,
					"panic recovered",
					attrs...,
				)

				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					response.Error{Message: "Internal server error"},
				)
			}
		}()

		c.Next()
	}
}
