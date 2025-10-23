package main

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	_ "github.com/DimitriLaPoudre/MusicShack/server/internal/plugin"
	_ "github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/routes"
)

func main() {
	database.InitDB()
	r := routes.SetupRouters()
	r.Run(":" + config.URL.Port())
}
