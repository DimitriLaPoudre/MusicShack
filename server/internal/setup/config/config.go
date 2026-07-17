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
		Admin   AdminConfig   `envPrefix:"ADMIN_"`
		Library LibraryConfig `envPrefix:"LIBRARY_"`
		DB      DBConfig      `envPrefix:"DB_"`
		Session SessionConfig `envPrefix:"SESSION_"`
		HTTP    HTTPConfig
		Log     LogConfig `envPrefix:"LOG_"`
	}

	AdminConfig struct {
		DefaultUsername string `env:"DEFAULT_USERNAME,unset" envDefault:"admin"`
		DefaultPassword string `env:"DEFAULT_PASSWORD,unset" envDefault:"changemenow"`
	}

	LibraryConfig struct {
		Path string `env:"PATH,required"`
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
)

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	info, err := os.Stat(cfg.Library.Path)
	if err != nil {
		slog.Error(fmt.Sprintf("library path invalid: %v", err))
		os.Exit(1)
	}
	if !info.IsDir() {
		slog.Error("library path is not a directory")
		os.Exit(1)
	}
	if err := IsWritableDirectory(cfg.Library.Path); err != nil {
		slog.Error(fmt.Sprintf("library path is not writable: %v", err))
		os.Exit(1)
	}

	return cfg, nil
}

func IsWritableDirectory(dir string) error {
	testFile := filepath.Join(dir, ".write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return err
	}
	_ = f.Close()
	if err := os.Remove(testFile); err != nil {
		return err
	}
	return nil
}
