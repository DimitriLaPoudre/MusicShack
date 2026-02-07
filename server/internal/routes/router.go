package routes

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/handlers"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
)

func SetupRouters() *graceful.Graceful {
	r, err := graceful.Default(graceful.WithAddr(":" + config.PORT))
	if err != nil {
		panic(err)
	}
	r.Use(cors.Default())

	buildDir := "../client_web/build"

	r.Static("/_app", filepath.Join(buildDir, "_app"))
	r.Static("/assets", filepath.Join(buildDir, "assets"))

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}
		c.File(filepath.Join(buildDir, "index.html"))
	})

	api := r.Group("/api")
	{

		api.GET("", handlers.Info)
		me := api.Group("/me")
		{
			me.Use(middlewares.Logged())
			me.GET("", handlers.Me)
			me.PUT("", handlers.MeUpdate)
		}

		api.POST("/login", middlewares.RateLimiter("5-M"), middlewares.LoggedOut(), handlers.Login)
		api.POST("/logout", middlewares.Logged(), handlers.Logout)

		api.GET("/song/:provider/:id", middlewares.Logged(), handlers.GetSong)
		api.GET("/album/:provider/:id", middlewares.Logged(), handlers.GetAlbum)
		api.GET("/artist/:provider/:id", middlewares.Logged(), handlers.GetArtist)
		api.GET("/search", middlewares.Logged(), handlers.Search)

		admin := api.Group("/admin")
		{
			admin.GET("", middlewares.Admin(), handlers.Admin)
			admin.POST("/login", middlewares.RateLimiter("5-M"), middlewares.Admout(), handlers.AdminLogin)
			admin.PUT("/password", middlewares.Admin(), handlers.AdminPassword)
			admin.POST("/logout", middlewares.Admin(), handlers.AdminLogout)
		}

		users := api.Group("/users")
		{
			users.Use(middlewares.Admin())
			users.POST("", handlers.CreateUser)
			users.GET("", handlers.ListUsers)
			users.GET("/:id", handlers.GetUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}

		apiInstance := api.Group("/instances")
		{
			apiInstance.Use(middlewares.Logged())
			apiInstance.POST("", handlers.AddInstance)
			apiInstance.GET("", handlers.ListInstances)
			apiInstance.DELETE("/:id", handlers.RemoveInstance)
		}

		downloads := api.Group("/downloads")
		{
			downloads.Use(middlewares.Logged())
			downloads.POST("", handlers.AddDownload)
			downloads.GET("", handlers.ListDownload)
			downloads.DELETE("/:id", handlers.DeleteDownload)
			downloads.POST("/:id/retry", handlers.RetryDownload)
			downloads.POST("/:id/cancel", handlers.CancelDownload)
			downloads.POST("/retry", handlers.RetryDownloads)
			downloads.POST("/done", handlers.DoneDownloads)
		}

		follows := api.Group("/follows")
		{
			follows.Use(middlewares.Logged())
			follows.POST("", handlers.AddFollow)
			follows.GET("", handlers.ListFollows)
			follows.DELETE("/:id", handlers.DeleteFollow)
		}

		library := api.Group("/library")
		{
			library.Use(middlewares.Logged())
			library.GET("", handlers.ListSong)
			library.GET("/:id/img", handlers.GetSongCover)
			library.DELETE("/:id", handlers.DeleteSong)
			library.PUT("", handlers.SyncLibrary)
		}
	}

	return r
}
