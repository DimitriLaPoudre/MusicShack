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

	apiInstance := r.Group("/api/instances/")
	{
		apiInstance.POST("/", middlewares.Logged(), handlers.AddInstance)
		apiInstance.GET("/", middlewares.Logged(), handlers.ListInstances)
		apiInstance.GET("/:id", middlewares.Logged(), handlers.GetInstance)
		apiInstance.POST("/:id", middlewares.Logged(), handlers.RemoveInstance)
	}

	users := r.Group("/api/users")
	{
		users.POST("/", handlers.CreateUser)
		users.GET("/", handlers.ListUsers)
		users.GET("/:id", handlers.GetUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}

	return r
}
