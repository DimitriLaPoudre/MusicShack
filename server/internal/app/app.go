package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/handler"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/middleware"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/router"
	"github.com/DimitriLaPoudre/MusicShack/internal/job"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/postgres"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/DimitriLaPoudre/MusicShack/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func Run(cfg *config.Config) {
	repo, err := postgres.New(cfg.DB.DSN())
	if err != nil {
		slog.Error("failed to create a new postgres repository", slog.String("err", err.Error()))
		os.Exit(1)
	}

	if err := repo.Migrate(cfg.DB.DSN()); err != nil {
		slog.Error("failed to migrate the database", slog.String("err", err.Error()))
		os.Exit(1)
	}

	pluginStore := service.NewPluginStoreService()
	pluginHifi := hifi.NewHifi()
	pluginStore.Register(&pluginHifi)

	authS := service.NewAuthService(cfg.Session, &repo, &repo)
	userS := service.NewUserService(cfg.Library, &repo)
	pluginS := service.NewPluginService(&pluginStore)

	authU := usecase.NewAuthUseCase(cfg.Session, &repo, &repo)
	userU := usecase.NewUserUseCase(cfg.Library, &userS, &repo)
	instanceU := usecase.NewInstanceUseCase(&pluginS, &repo)
	pluginU := usecase.NewPluginUseCase(&pluginS, &pluginStore, &repo)

	meH := handler.NewMeHandler(&userU)
	userH := handler.NewUserHandler(&userU)
	authH := handler.NewAuthHandler(cfg.HTTP, cfg.Session, &authU)
	instanceH := handler.NewInstanceHandler(&instanceU)
	pluginH := handler.NewPluginHandler(&pluginU)

	authMW := middleware.AuthMiddleware(cfg.Session, &authS)
	adminMW := middleware.AdminMiddleware(&authS)
	userMW := middleware.UserMiddleware(&authS)
	guestMW := middleware.GuestMiddleware(&authS)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	service.InitAdmin(ctx, cfg.Admin, &userS, &repo)

	c := cron.New()
	// if err := job.FetchFollows(c, ctx, l); err != nil {
	// 	l.Fatal().Err(err).Msg("failed to start job: FetchFollows")
	// }
	if err := job.CleanExpiredSession(c, ctx, &repo); err != nil {
		slog.Error("start the job: CleanExpiredSession", slog.String("err", err.Error()))
		os.Exit(1)
	}
	c.Start()

	app := gin.New()
	router.New(app, cfg,
		authMW,
		adminMW,
		userMW,
		guestMW,
		&meH,
		&userH,
		&authH,
		&instanceH,
		&pluginH,
	)
	httpServ := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: app,
	}
	go func() {
		slog.Info("server starting on " + httpServ.Addr)
		if err := httpServ.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", slog.String("err", err.Error()))
		}
	}()

	<-ctx.Done()
	slog.Info("CTRL-C successfully handle")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := httpServ.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.String("err", err.Error()))
	}
	c.Stop().Done()
}
