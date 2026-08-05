package app

import (
	"context"
	"errors"
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
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/sqlite"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func Run(cfg *config.Config) {
	repo, err := sqlite.New(cfg.DB.Path)
	if err != nil {
		slog.Error("failed to create a new postgres repository", slog.String("err", err.Error()))
		os.Exit(1)
	}

	if err := repo.Migrate(cfg.DB.Path); err != nil {
		slog.Error("failed to migrate the database", slog.String("err", err.Error()))
		os.Exit(1)
	}

	pluginStore := service.NewPluginStoreService()
	pluginHifi := hifi.NewHifi()
	pluginStore.Register(&pluginHifi)

	pluginCache := service.NewPluginCacheService(cfg.Plugin)

	authS := service.NewAuthService(cfg.Session, &repo, &repo)
	userS := service.NewUserService(cfg.Download, &repo)
	pluginS := service.NewPluginService(&pluginCache, &pluginStore, &repo, &repo, &repo)
	instanceS := service.NewInstanceService(&pluginS, &repo)
	metadataS := service.NewMetadataService(&pluginS)
	downloadS := service.NewDownloadService(cfg.Download, cfg.Plugin, &pluginS, &metadataS, &repo, &repo)
	followS := service.NewFollowService(&pluginS, &repo)
	fetchNewReleasesS := service.NewFetchNewReleasesService(cfg.Plugin, &pluginS, &followS, downloadS)

	userH := handler.NewUserHandler(&userS)
	authH := handler.NewAuthHandler(cfg.HTTP, cfg.Session, &authS)
	instanceH := handler.NewInstanceHandler(&instanceS)
	followH := handler.NewFollowHandler(&followS)
	pluginH := handler.NewPluginHandler(&pluginS)
	downloadH := handler.NewDownloadHandler(downloadS)

	authMW := middleware.AuthMiddleware(cfg.Session, &authS)
	adminMW := middleware.AdminMiddleware()
	targetUserMW := middleware.TargetUserMiddleware(&repo)
	guestMW := middleware.GuestMiddleware()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := service.InitAdmin(ctx, cfg.Admin, &userS, &repo); err != nil && !errors.Is(err, model.ErrAdminAlreadyExist) {
		slog.Warn("admin initialisation failed", slog.String("err", err.Error()))
	}

	c := cron.New()
	if err := job.DownloadNewReleases(c, ctx, &fetchNewReleasesS); err != nil {
		slog.Error("start the job: FetchFollows", slog.String("err", err.Error()))
		os.Exit(1)
	}
	if err := job.CleanExpiredSession(c, ctx, &authS); err != nil {
		slog.Error("start the job: CleanExpiredSession", slog.String("err", err.Error()))
		os.Exit(1)
	}
	if err := job.CleanExpiredPluginCache(c, ctx, &pluginCache); err != nil {
		slog.Error("start the job: CleanExpiredPluginCache", slog.String("err", err.Error()))
		os.Exit(1)
	}
	c.Start()

	app := gin.New()
	router.New(app,
		authMW,
		adminMW,
		targetUserMW,
		guestMW,
		&userH,
		&authH,
		&instanceH,
		&pluginH,
		&downloadH,
		&followH,
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
