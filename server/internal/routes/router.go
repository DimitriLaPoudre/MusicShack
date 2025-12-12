package routes

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/handlers"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/api", handlers.Info)
	r.GET("/api/me", middlewares.Logged(), handlers.Me)

	r.POST("/api/signup", middlewares.LoggedOut(), handlers.Signup)
	r.POST("/api/login", middlewares.LoggedOut(), handlers.Login)
	r.POST("/api/logout", middlewares.Logged(), handlers.Logout)

	apiInstance := r.Group("/api/instances")
	{
		apiInstance.Use(middlewares.Logged())
		apiInstance.POST("/", handlers.AddInstance)
		apiInstance.GET("/", handlers.ListInstances)
		apiInstance.GET("/:id", handlers.GetInstance)
		apiInstance.DELETE("/:id", handlers.RemoveInstance)
	}

	users := r.Group("/api/users")
	{
		users.Use(middlewares.Logged())
		users.POST("/", handlers.CreateUser)
		users.GET("/", handlers.ListUsers)
		users.GET("/:id", handlers.GetUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)

		downloads := users.Group("/downloads")
		{
			downloads.POST("/song/:api/:id", handlers.AddDownloadSong)
			downloads.POST("/album/:api/:id", handlers.AddDownloadAlbum)
			downloads.POST("/artist/:api/:id", handlers.AddDownloadArtist)
			downloads.GET("/", handlers.ListDownload)
			downloads.DELETE("/:id", handlers.DeleteDownload)
			downloads.POST("/retry/:id", handlers.RetryDownload)
			downloads.POST("/cancel/:id", handlers.CancelDownload)
		}
	}

	r.GET("/api/song/:api/:id", middlewares.Logged(), handlers.GetSong)
	r.GET("/api/album/:api/:id", middlewares.Logged(), handlers.GetAlbum)
	r.GET("/api/artist/:api/:id", middlewares.Logged(), handlers.GetArtist)
	r.GET("/api/search", middlewares.Logged(), handlers.Search)

	return r
}
