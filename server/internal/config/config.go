package config

import (
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var PORT string
var URL url.URL
var JWT_SECRET []byte
var DOWNLOAD_FOLDER string

func checkDownloadDirectory(dir string) error {
	testFile := filepath.Join(dir, ".write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return err
	}
	f.Close()
	if err := os.Remove(testFile); err != nil {
		return err
	}
	return nil
}

func init() {
	godotenv.Load("../.env")

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
		log.Println("PORT is missing - defaulting to 8080")
		port = "8080"
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
		log.Fatal("DOWNLOAD: ", err)
	}
	if !info.IsDir() {
		log.Fatal("DOWNLOAD is not a directory")
	}
	if err := checkDownloadDirectory(folder); err != nil {
		log.Fatal("DOWNLOAD can't be written in: ", err)
	}

	DOWNLOAD_FOLDER = folder
}
