package config

import (
	"fmt"
	"log"
	"net/url"
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
		Host      string `env:"HOST" envDefault:"localhost"`
		Port      int    `env:"PORT" envDefault:"5432"`
		Name      string `env:"NAME,unset,required"`
		User      string `env:"USER,unset,required"`
		Password  string `env:"PASS,unset,required"`
		Params    string `env:"PARAMS" envDefault:"sslmode=disable"`
		Migration string `env:"MIGRATION"`
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

func (c *DBConfig) DSN() string {
	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:     "/" + c.Name,
		RawQuery: c.Params,
	}
	return u.String()
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	info, err := os.Stat(cfg.Library.Path)
	if err != nil {
		log.Fatal("LIBRARY_PATH: ", err)
	}
	if !info.IsDir() {
		log.Fatal("LIBRARY_PATH is not a directory")
	}
	if err := checkLibraryDirectory(cfg.Library.Path); err != nil {
		log.Fatal("LIBRARY_PATH can't be written in: ", err)
	}

	return cfg, nil
}

func checkLibraryDirectory(dir string) error {
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
