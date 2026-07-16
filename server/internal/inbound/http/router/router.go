package router

import (
	"net/http"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/handler"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/middleware"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/gin-gonic/gin"
)

func New(
	app *gin.Engine,
	cfg *config.Config,

	authMW gin.HandlerFunc,
	adminMW gin.HandlerFunc,
	userMW gin.HandlerFunc,
	guestMW gin.HandlerFunc,

	meH *handler.MeHandler,
	userH *handler.UserHandler,
	authH *handler.AuthHandler,
	instanceH *handler.InstanceHandler,
	pluginH *handler.PluginHandler,
) {
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	gin.Recovery()

	app.GET("/healthz", healthz)

	api := app.Group("/api")
	{
		meGroup := api.Group("/me")
		{
			meGroup.Use(middleware.RateLimiter(time.Minute, 100))
			meGroup.Use(authMW)
			meGroup.GET("", meH.Get)
			meGroup.PUT("", meH.Update)

			instancesGroup := meGroup.Group("/instances")
			{
				instancesGroup.GET("", instanceH.ListForMe)
				instancesGroup.POST("", instanceH.CreateForMe)
				instancesGroup.DELETE(":id", instanceH.DeleteForMe)
			}
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
			authGroup.POST("/login", middleware.RateLimiter(time.Minute, 10), guestMW, authH.Login)
			authGroup.DELETE("/logout", middleware.RateLimiter(time.Minute, 10), authMW, authH.Logout)
		}

		pluginGroup := api.Group("/plugin")
		{
			pluginGroup.GET("/song/:provider/:id", authMW, pluginH.GetSong)
			pluginGroup.GET("/album/:provider/:id", authMW, pluginH.GetAlbum)
			pluginGroup.GET("/artist/:provider/:id", authMW, pluginH.GetArtist)
			pluginGroup.GET("/playlist/:provider/:id", authMW, pluginH.GetPlaylist)
			pluginGroup.GET("/search", authMW, pluginH.GetSearch)
		}
	}
}

func healthz(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
