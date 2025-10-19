package main

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/routes"
)

func main() {
	database.InitDB()
	r := routes.SetupRouters()
	r.Run(":" + config.URL.Port())
}
