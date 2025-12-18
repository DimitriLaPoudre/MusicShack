package config

import (
	"log"
	"net/url"
	"os"
)

var PORT string
var URL url.URL
var JWT_SECRET []byte
var DOWNLOAD_FOLDER string

func init() {
	urlRaw := os.Getenv("URL")
	if urlRaw == "" {
		log.Fatal("URL is missing")
	}
	parsed, err := url.Parse(urlRaw)
	if err != nil {
		log.Fatal("URL invalid: ", err)
	}
	URL = *parsed

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is missing")
	}
	PORT = port

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is missing")
	}
	JWT_SECRET = []byte(secret)

	folder := os.Getenv("DOWNLOAD")
	if folder == "" {
		log.Fatal("DOWNLOAD is missing")
	}
	info, err := os.Stat(folder)
	if err != nil {
		log.Fatal("DOWNLOAD_FOLDER: ", err)
	}
	if !info.IsDir() {
		log.Fatal("DOWNLOAD is not a directory")
	}
	DOWNLOAD_FOLDER = folder
}
