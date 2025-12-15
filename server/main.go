package main

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	_ "github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/routes"
)

func main() {
	r := routes.SetupRouters()
	r.Run(":" + config.URL.Port())
}
