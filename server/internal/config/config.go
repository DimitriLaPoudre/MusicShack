package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

var PORT string
var HTTPS bool
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

	https, err := strconv.ParseBool(os.Getenv("HTTPS"))
	if err != nil {
		log.Println("HTTPS is invalid: ", err.Error())
		https = false
	}
	HTTPS = https

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is missing - defaulting to 8080")
		port = "8080"
	}
	PORT = port

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
