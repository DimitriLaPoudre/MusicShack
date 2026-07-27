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
	targetUserMW gin.HandlerFunc,
	guestMW gin.HandlerFunc,

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
		meGroup := api.Group("/me",
			middleware.RateLimiter(time.Minute, 100),
			authMW,
			targetUserMW,
		)
		{
			meGroup.GET("", userH.GetByID)
			meGroup.PUT("", userH.Update)
			meGroup.DELETE("", userH.Delete)
		}

		instancesGroup := meGroup.Group("/instances",
			middleware.RateLimiter(time.Minute, 100),
			authMW,
		)
		{
			instancesGroup.GET("", instanceH.List)
			instancesGroup.POST("", instanceH.Create)
			instancesGroup.DELETE(":id", instanceH.Delete)
		}

		usersGroup := api.Group("/users",
			middleware.RateLimiter(time.Minute, 100),
			authMW,
			adminMW,
		)
		{
			usersGroup.POST("", userH.Create)
			usersGroup.GET("", userH.List)
			usersGroup.GET("/:user_id", userH.GetByID)
			usersGroup.PUT("/:user_id", userH.Update)
			usersGroup.DELETE("/:user_id", userH.Delete)
		}

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", middleware.RateLimiter(time.Minute, 5), guestMW, authH.Login)
			authGroup.DELETE("/logout", middleware.RateLimiter(time.Minute, 10), authMW, authH.Logout)
		}

		pluginGroup := api.Group("",
			middleware.RateLimiter(time.Minute, 100),
			authMW,
		)
		{
			pluginGroup.GET("/song/:provider/:id", pluginH.GetSongInfo)
			pluginGroup.GET("/album/:provider/:id", pluginH.GetAlbumInfo)
			pluginGroup.GET("/album/:provider/:id/songs", pluginH.GetAlbumSongs)
			pluginGroup.GET("/artist/:provider/:id", pluginH.GetArtistInfo)
			pluginGroup.GET("/artist/:provider/:id/albums", pluginH.GetArtistAlbums)
			pluginGroup.GET("/playlist/:provider/:id", pluginH.GetPlaylistInfo)
			pluginGroup.GET("/playlist/:provider/:id/songs", pluginH.GetPlaylistSongs)
			pluginGroup.GET("/search", pluginH.GetSearch)
		}

		downloadGroup := api.Group("/download", middleware.RateLimiter(time.Minute, 100), authMW)
		{

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

		followGroup := api.Group("/follows", middleware.RateLimiter(time.Minute, 100), authMW)
		{
			followGroup.POST("/artist", followH.Add)
			followGroup.GET("", followH.List)
			followGroup.DELETE("/:id", followH.Delete)
		}

		// libraryGroup := api.Group("/library", middleware.RateLimiter(time.Minute, 100), authMW)
		{
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
