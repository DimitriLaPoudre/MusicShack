package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	PORT         string
	HTTPS        bool
	LIBRARY_PATH string
)

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

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println(".env not found")
	}

	https := os.Getenv("HTTPS")
	if https == "" {
		log.Println("HTTPS is missing - defaulting to false")
		HTTPS = false
	} else {
		if value, err := strconv.ParseBool(https); err != nil {
			log.Println("HTTPS is invalid: ", err.Error(), " - defaulting to false")
			HTTPS = false
		} else {
			HTTPS = value
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is missing - defaulting to 8080")
		port = "8080"
	}
	PORT = port

	folder := os.Getenv("LIBRARY_PATH")
	if folder == "" {
		log.Fatal("LIBRARY_PATH is missing")
	}
	info, err := os.Stat(folder)
	if err != nil {
		log.Fatal("LIBRARY_PATH: ", err)
	}
	if !info.IsDir() {
		log.Fatal("LIBRARY_PATH is not a directory")
	}
	if err := checkLibraryDirectory(folder); err != nil {
		log.Fatal("LIBRARY_PATH can't be written in: ", err)
	}
	LIBRARY_PATH = folder
}
