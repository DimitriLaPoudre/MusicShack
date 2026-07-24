package router

import (
	"net/http"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/handler"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/middleware"
	"github.com/gin-gonic/gin"
)

func New(
	app *gin.Engine,

	authMW gin.HandlerFunc,
	adminMW gin.HandlerFunc,
	userMW gin.HandlerFunc,
	guestMW gin.HandlerFunc,

	meH *handler.MeHandler,
	userH *handler.UserHandler,
	authH *handler.AuthHandler,
	instanceH *handler.InstanceHandler,
	pluginH *handler.PluginHandler,
	downloadH *handler.DownloadHandler,
	followH *handler.FollowHandler,
) {
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())

	app.GET("/healthz", healthz)

	api := app.Group("/api")
	{
		meGroup := api.Group("/me")
		{
			meGroup.Use(middleware.RateLimiter(time.Minute, 100))
			meGroup.Use(authMW)

			meGroup.GET("", meH.Get)
			meGroup.PUT("", meH.Update)

		}

		instancesGroup := meGroup.Group("/instances")
		{
			instancesGroup.Use(middleware.RateLimiter(time.Minute, 100))
			instancesGroup.Use(authMW)

			instancesGroup.GET("", instanceH.ListForMe)
			instancesGroup.POST("", instanceH.CreateForMe)
			instancesGroup.DELETE(":id", instanceH.DeleteForMe)
		}

		usersGroup := api.Group("/users")
		{
			usersGroup.Use(middleware.RateLimiter(time.Minute, 100))
			usersGroup.Use(authMW)
			usersGroup.Use(adminMW)

			usersGroup.POST("", userH.Create)
			usersGroup.GET("", userH.List)
			usersGroup.GET("/:id", userH.GetByID)
			usersGroup.PUT("/:id", userH.Update)
			usersGroup.DELETE("/:id", userH.Delete)
		}

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", middleware.RateLimiter(time.Minute, 5), guestMW, authH.Login)
			authGroup.DELETE("/logout", middleware.RateLimiter(time.Minute, 10), authMW, authH.Logout)
		}

		pluginGroup := api.Group("/plugin")
		{
			pluginGroup.Use(middleware.RateLimiter(time.Minute, 100))
			pluginGroup.Use(authMW)

			pluginGroup.GET("/song/:provider/:id", pluginH.GetSong)
			pluginGroup.GET("/album/:provider/:id", pluginH.GetAlbum)
			pluginGroup.GET("/artist/:provider/:id", pluginH.GetArtist)
			pluginGroup.GET("/playlist/:provider/:id", pluginH.GetPlaylist)
			pluginGroup.GET("/search", pluginH.GetSearch)
		}

		downloadGroup := api.Group("/download")
		{
			downloadGroup.Use(middleware.RateLimiter(time.Minute, 100))
			downloadGroup.Use(authMW)

			downloadGroup.POST("/song", downloadH.DownloadSong)
			downloadGroup.POST("/album", downloadH.DownloadAlbum)
			downloadGroup.POST("/artist", downloadH.DownloadArtist)
			downloadGroup.POST("/playlist", downloadH.DownloadPlaylist)

			downloadGroup.PUT("/:id/retry", downloadH.Retry)
			downloadGroup.PUT("/retry", downloadH.RetryAll)
			downloadGroup.PUT("/:id/cancel", downloadH.Cancel)
			downloadGroup.PUT("/remove", downloadH.RemoveDone)
			downloadGroup.PUT("/:id/remove", downloadH.Remove)
			downloadGroup.GET("", downloadH.List)
		}

		followGroup := api.Group("/follows")
		{
			followGroup.Use(middleware.RateLimiter(time.Minute, 100))
			followGroup.Use(authMW)

			followGroup.POST("/artist", followH.Add)
			followGroup.GET("", followH.List)
			followGroup.DELETE("/:id", followH.Delete)
		}

		libraryGroup := api.Group("/library")
		{
			libraryGroup.Use(middleware.RateLimiter(time.Minute, 100))
			libraryGroup.Use(authMW)

			// libraryGroup.GET("", libraryH.List)
			// libraryGroup.POST("", libraryH.Upload)
			// libraryGroup.PUT("/:id", libraryH.Edit)
			// libraryGroup.GET("/:id/img", libraryH.GetCover)
			// libraryGroup.DELETE("/:id", libraryH.Delete)
			// libraryGroup.PUT("", libraryH.Sync)
		}
	}
}

func healthz(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
