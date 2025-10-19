package routes

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/handlers"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	r := gin.Default()

	r.GET("/api", handlers.Info)
	r.GET("/api/me", middlewares.Logged(), handlers.Me)

	r.POST("/api/signup", middlewares.LoggedOut(), handlers.Signup)
	r.POST("/api/login", middlewares.LoggedOut(), handlers.Login)

	// users := r.Group("/users")
	// {
	// 	users.POST("/", CreateUser)
	// 	users.GET("/", ListUsers)
	// 	users.GET("/:id", GetUser)
	// 	users.PUT("/:id", UpdateUser)
	// 	users.DELETE("/:id", DeleteUser)
	// }

	return r
}
