package app

import (
	"context"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/handler"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/router"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/postgres"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

func Run(l *zerolog.Logger, cfg *config.Config) {
	repo, err := postgres.New(l, cfg.DB.DSN(), cfg.DB.Migration)
	if err != nil {
		l.Fatal().Msg(err.Error())
	}

	adminS := usecase.NewAdminUseCase(l, &cfg.Admin, &repo, &repo)

	adminH := handler.NewAdminHandler(l, &adminS)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	c := cron.New()
	// if err := job.FetchFollows(c, ctx, l); err != nil {
	// 	l.Fatal().Err(err).Msg("failed to start job: FetchFollows")
	// }
	// if err := job.CleanExpiredSession(c, ctx, l, &repo); err != nil {
	// 	l.Fatal().Err(err).Msg("failed to start job: FetchFollows")
	// }
	c.Start()

	app := gin.New()
	router.New(app, l, cfg,
		&adminH,
	)
	httpServ := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: app,
	}
	go func() {
		l.Info().Msg("server starting on " + httpServ.Addr)
		if err := httpServ.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error().Err(err).Msg("failed to start server")
		}
	}()

	<-ctx.Done()
	l.Info().Msg("CTRL-C successfully handle")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := httpServ.Shutdown(ctx); err != nil {
		l.Error().Err(err).Msg("server forced to shutdown")
	}
	c.Stop().Done()
}
