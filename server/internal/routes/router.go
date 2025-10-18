package routes

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	r := gin.Default()

	r.GET("/api", handlers.Info)

	r.POST("/api/signup", handlers.Signup)
	r.POST("/api/login", handlers.Login)

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
