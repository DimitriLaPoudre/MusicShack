package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

func New(level string, pretty bool) *slog.Logger {
	var output io.Writer = os.Stdout

	handlerOpts := &slog.HandlerOptions{Level: parseLevel(level)}

	var handler slog.Handler
	if pretty {
		handler = slog.NewTextHandler(output, handlerOpts)
	} else {
		handler = slog.NewJSONHandler(output, handlerOpts)
	}

	handler = &contextHandler{
		Handler: handler,
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	MeKey        contextKey = "me"
)

type contextHandler struct {
	slog.Handler
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		r.AddAttrs(slog.String("request_id", requestID))

		if me, ok := ctx.Value(MeKey).(model.User); ok {
			r.AddAttrs(slog.Group("user", slog.String("id", me.ID.String()), slog.String("username", me.Username), slog.String("role", string(me.Role))))
		}
	}

	return h.Handler.Handle(ctx, r)
}
