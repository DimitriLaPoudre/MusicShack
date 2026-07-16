package middleware

import (
	"log/slog"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()
		latency := time.Since(start)
		status := c.Writer.Status()

		if query != "" {
			path += "?" + query
		}

		requestID, err := utils.GetFromContext[string](c, "request_id")
		if err != nil {
			requestID = "unknown"
		}

		attrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("ip", clientIP),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
		}

		var level slog.Level
		switch {
		case status >= 500:
			level = slog.LevelError
		case status >= 400:
			level = slog.LevelWarn
		default:
			level = slog.LevelInfo
		}
		slog.LogAttrs(c.Request.Context(), level, "HTTP request", attrs...)
	}
}
