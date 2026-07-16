package main

import (
	"log/slog"
	"os"

	"github.com/DimitriLaPoudre/MusicShack/internal/app"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	l := logger.New(cfg.Log.Level, cfg.Log.Pretty)
	slog.SetDefault(l)

	app.Run(cfg)
}
