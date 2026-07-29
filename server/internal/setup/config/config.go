package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		Admin    AdminConfig   `envPrefix:"ADMIN_"`
		DB       DBConfig      `envPrefix:"DB_"`
		Session  SessionConfig `envPrefix:"SESSION_"`
		HTTP     HTTPConfig
		Log      LogConfig      `envPrefix:"LOG_"`
		Download DownloadConfig `envPrefix:"DOWNLOAD_"`
		Plugin   PluginConfig   `envPrefix:"PLUGIN_"`
	}

	AdminConfig struct {
		DefaultUsername string `env:"DEFAULT_USERNAME,unset" envDefault:"admin"`
		DefaultPassword string `env:"DEFAULT_PASSWORD,unset" envDefault:"changemenow"`
	}

	DBConfig struct {
		Path string `env:"PATH,unset,required"`
	}

	SessionConfig struct {
		CookieName  string        `env:"COOKIE_NAME" envDefault:"user_token"`
		Exp         time.Duration `env:"EXP" envDefault:"168h"`
		RememberExp time.Duration `env:"REMEMBER_EXP" envDefault:"720h"`
	}

	HTTPConfig struct {
		Port  int  `env:"PORT" envDefault:"8080"`
		HTTPS bool `env:"HTTPS" envDefault:"false"`
	}

	LogConfig struct {
		Pretty bool   `env:"PRETTY" envDefault:"false"`
		Level  string `env:"LEVEL" envDefault:"info"`
	}

	DownloadConfig struct {
		Path       string `env:"PATH,unset,required"`
		Concurrent int    `env:"CONCURRENT" envDefault:"3"`
	}

	PluginConfig struct {
		Pagination struct {
			Limit  int `env:"LIMIT" envDefault:"25"`
			Offset int `env:"OFFSET" envDefault:"0"`
		} `envPrefix:"PAGINATION_"`
		Cache struct {
			Expiration time.Duration `env:"EXPIRATION" envDefault:"30m"`
		} `envPrefix:"CACHE_"`
	}
)

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	cfg.Download.Path = filepath.Clean(cfg.Download.Path)
	libraryPath := cfg.Download.Path
	info, err := os.Stat(libraryPath)
	if err != nil {
		slog.Error(fmt.Sprintf("library path %s is invalid", libraryPath), slog.String("err", err.Error()))
		os.Exit(1)
	}
	if !info.IsDir() {
		slog.Error(fmt.Sprintf("library path %s is not a directory", libraryPath))
		os.Exit(1)
	}
	if err := IsWritableDirectory(libraryPath); err != nil {
		slog.Error(fmt.Sprintf("library path %s is not writable", libraryPath), slog.String("err", err.Error()))
		os.Exit(1)
	}

	return cfg, nil
}

func IsWritableDirectory(dir string) error {
	testFile := filepath.Join(dir, ".write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("testFile creation: %w", err)
	}
	_ = f.Close()
	if err := os.Remove(testFile); err != nil {
		return fmt.Errorf("testFile deletion: %w", err)
	}
	return nil
}
