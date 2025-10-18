package main

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/routes"
	"os"
)

func main() {
	database.InitDB()
	r := routes.SetupRouters()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
