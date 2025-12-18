package routes

import (
	"fmt"
	"net/http"
	"os"
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
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		fmt.Println("not exist")
	}

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
			me.PUT("", handlers.UpdateMe)
			me.DELETE("", handlers.DeleteMe)
		}

		api.POST("/signup", middlewares.LoggedOut(), handlers.Signup)
		api.POST("/login", middlewares.LoggedOut(), handlers.Login)
		api.POST("/logout", middlewares.Logged(), handlers.Logout)

		apiInstance := api.Group("/instances")
		{
			apiInstance.Use(middlewares.Logged())
			apiInstance.POST("", handlers.AddInstance)
			apiInstance.GET("", handlers.ListInstances)
			apiInstance.GET("/:id", handlers.GetInstance)
			apiInstance.DELETE("/:id", handlers.RemoveInstance)
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

		downloads := api.Group("/users/downloads")
		{
			downloads.Use(middlewares.Logged())
			downloads.POST("/song/:api/:id", handlers.AddDownloadSong)
			downloads.POST("/album/:api/:id", handlers.AddDownloadAlbum)
			downloads.POST("/artist/:api/:id", handlers.AddDownloadArtist)
			downloads.GET("", handlers.ListDownload)
			downloads.DELETE("/:id", handlers.DeleteDownload)
			downloads.POST("/retry/:id", handlers.RetryDownload)
			downloads.POST("/cancel/:id", handlers.CancelDownload)
		}

		follows := api.Group("/follows")
		{
			follows.Use(middlewares.Logged())
			follows.POST("", handlers.AddFollow)
			follows.GET("", handlers.ListFollows)
			follows.DELETE("/:id", handlers.DeleteFollow)
		}

		api.GET("/song/:api/:id", middlewares.Logged(), handlers.GetSong)
		api.GET("/album/:api/:id", middlewares.Logged(), handlers.GetAlbum)
		api.GET("/artist/:api/:id", middlewares.Logged(), handlers.GetArtist)
		api.GET("/search", middlewares.Logged(), handlers.Search)
	}

	return r
}
