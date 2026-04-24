package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/routes"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils/autofetch"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cron := autofetch.AutoFetch(ctx)
	r := routes.SetupRouters()
	defer r.Close()
	if err := r.RunWithContext(ctx); err != nil && err != context.Canceled {
		panic(err)
	}
	<-cron.Stop().Done()
	log.Println("CTRL-C successfully handle")
}
