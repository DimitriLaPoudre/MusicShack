package config

import (
	"log"
	"net/url"
	"os"
)

var URL url.URL
var JWT_SECRET []byte
var DOWNLOAD_FOLDER string

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
		log.Fatal("JWT_SECRET is missing")
	}
	JWT_SECRET = []byte(secret)

	folder := os.Getenv("DOWNLOAD_FOLDER")
	if folder == "" {
		log.Fatal("DOWNLOAD_FOLDER is missing")
	}
	info, err := os.Stat(folder)
	if err != nil {
		log.Fatal("DOWNLOAD_FOLDER: ", err)
	}
	if !info.IsDir() {
		log.Fatal("DOWNLOAD_FOLDER is not a directory")
	}
	DOWNLOAD_FOLDER = folder
}
