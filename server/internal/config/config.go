package config

import (
	"log"
	"net/url"
	"os"
)

var URL url.URL
var JWT_SECRET []byte

func init() {
	parsed, err := url.Parse(os.Getenv("URL"))
	if err != nil {
		log.Fatal("URL invalid: ", err)
	}
	if parsed.Port() == "" {
		log.Fatal("URL port invalid")
	}
	URL = *parsed

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is empty")
	}
	JWT_SECRET = []byte(secret)
}
